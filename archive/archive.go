package archive

import (
	"archive/zip"
	"bytes"
	"github.com/cloudogu/cesapp-lib/core"
	"io"
	"os"
	"strings"
	"time"
)

var log = core.GetLogger()

type FileReaderFunc = func(name string) ([]byte, error)
type FileStatFunc = func(name string) (os.FileInfo, error)
type WriteToZipFunc = func(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error
type SaveFileFunc = func(name string, data []byte, perm os.FileMode) error
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
	save       SaveFileFunc
}

type File struct {
	NameOutside string
	NameInside  string
}

func NewManager() *Manager {
	buffer := bytes.NewBuffer(nil)
	writer := zip.NewWriter(buffer)
	return &Manager{
		buffer:     buffer,
		writer:     writer,
		readFile:   os.ReadFile,
		stat:       os.Stat,
		writeToZip: WriteContentToZip,
		close:      writer.Close,
		save:       os.WriteFile,
	}
}

func WriteContentToZip(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error {
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
	return m.AddContentAsFileWithModifiedDate(content, fileNameInArchive, time.Now())
}

func (m Manager) AddContentAsFileWithModifiedDate(content string, fileNameInArchive string, modified time.Time) error {
	header := zip.FileHeader{
		Name:     fileNameInArchive,
		Modified: modified,
		Method:   zip.Deflate,
	}

	return m.writeToZip(m.writer, &header, []byte(content))
}

func (m Manager) AddFileToArchive(file File) error {
	log.Debugf("Adding file '%s' as '%s' to archive.", file.NameOutside, file.NameInside)
	content, err := m.readFile(file.NameOutside)
	if err != nil {
		return err
	}

	fileInfo, err := m.stat(file.NameOutside)
	if err != nil {
		return err
	}

	// Error can be ignored because the function actually never returns an error
	header, _ := zip.FileInfoHeader(fileInfo)

	header.Method = zip.Deflate
	header.Name = file.NameInside

	return m.writeToZip(m.writer, header, content)
}

func (m Manager) AddFilesToArchive(files []File, closeAfterFinish bool) error {
	for _, file := range files {
		err := m.AddFileToArchive(file)
		if err != nil {
			return err
		}
	}

	if closeAfterFinish {
		return m.Close()
	}

	return nil
}

func (m Manager) SaveArchiveAsFile(archivePath string) error {
	if !strings.HasSuffix(archivePath, ".zip") {
		log.Warning("File ending .zip was not provided.")
	}
	if strings.HasSuffix(archivePath, "/") {
		log.Warning("Incorrect file path was provided. Adding 'archive.zip' to file path.")
		archivePath += "archive.zip"
	}

	content := m.GetContent()

	return m.save(archivePath, content, 0755)
}

func (m Manager) Close() error {
	return m.close()
}
