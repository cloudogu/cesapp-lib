package archive2

import (
	"archive/zip"
	"bytes"
	"github.com/cloudogu/cesapp-lib/core"
	"io"
	"os"
	"time"
)

var log = core.GetLogger()

type FileReaderFunc = func(name string) ([]byte, error)
type FileStatFunc = func(name string) (os.FileInfo, error)
type WriteToZipFunc = func(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error
type CloseFUnc = func() error

type ZipWriter interface {
	CreateHeader(fh *zip.FileHeader) (io.Writer, error)
	Close() error
}

type Manager struct {
	buffer     *bytes.Buffer
	writer     *zip.Writer
	readFile   FileReaderFunc
	stat       FileStatFunc
	writeToZip WriteToZipFunc
	close      CloseFUnc
}

type File struct {
	nameOutside string
	nameInside  string
}

func NewManager() *Manager {
	buffer := bytes.NewBuffer(nil)
	writer := zip.NewWriter(buffer)
	return &Manager{
		buffer:     buffer,
		writer:     writer,
		readFile:   os.ReadFile,
		stat:       os.Stat,
		writeToZip: WriteToZip,
		close:      writer.Close,
	}
}

func WriteToZip(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error {
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	writtenBytes, err := writer.Write(content)
	if err != nil {
		return err
	}

	log.Debugf("wrote %d byte(s) to archive", writtenBytes)

	return nil
}

func (m Manager) GetContent() []byte {
	return m.buffer.Bytes()
}

func (m Manager) AddContentAsFile(content string, fileNameInArchive string) error {
	header := zip.FileHeader{
		Name:     fileNameInArchive,
		Modified: time.Now(),
		Method:   zip.Deflate,
	}

	return m.writeToZip(m.writer, &header, []byte(content))
}

func (m Manager) AddFileToArchive(file File) error {
	content, err := m.readFile(file.nameOutside)
	if err != nil {
		return err
	}

	fileInfo, err := m.stat(file.nameOutside)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate
	header.Name = file.nameInside

	return m.writeToZip(m.writer, header, content)
}

func (m Manager) Close() error {
	return m.close()
}
