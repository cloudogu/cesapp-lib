package archive

import (
	"archive/zip"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"io"
	"os"
	"strings"
)

var log = core.GetLogger()

type fileHandler interface {
	Create(filename string) (*os.File, error)
	Open(filename string) (*os.File, error)
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	GetFileInfoHeader(filePath string) (*zip.FileHeader, error)
}

type zipWriter interface {
	CreateHeader(fh *zip.FileHeader) (io.Writer, error)
	Close() error
}

type osFileHandler struct{}

func (d *osFileHandler) Open(filename string) (*os.File, error) {
	return os.Open(filename)
}

func (d *osFileHandler) Create(filename string) (*os.File, error) {
	return os.Create(filename)
}

func (d *osFileHandler) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (d *osFileHandler) GetFileInfoHeader(filePath string) (*zip.FileHeader, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return nil, err
	}

	header.Method = zip.Deflate

	return header, nil
}

// Handler can Create a zip archive and add files to it.
type Handler struct {
	writer      zipWriter
	fileHandler fileHandler
}

func InitIn(path string) (*Handler, error) {
	handler := &osFileHandler{}
	file, err := handler.Create(path)
	if err != nil {
		return nil, err
	}

	return &Handler{
		writer:      zip.NewWriter(file),
		fileHandler: handler,
	}, nil
}

// AppendFileToArchive takes a path to file that is read and appended to an archive.
// make sure to call the Close method when you're done with appending files to the archive.
func (ar *Handler) AppendFileToArchive(fileToZipPath string, filepathInZip string) error {
	log.Debugf("Append file %s to Archive as %s", fileToZipPath, filepathInZip)
	file, err := ar.fileHandler.Open(fileToZipPath)
	if err != nil {
		return fmt.Errorf("failed to read base file for appending to archive: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	handler, err := ar.fileHandler.GetFileInfoHeader(fileToZipPath)
	if err != nil {
		return err
	}

	fileInZip, err := ar.writer.CreateHeader(handler)
	if err != nil {
		return fmt.Errorf("failed to Create file in archive: %w", err)
	}

	if _, err := ar.fileHandler.Copy(fileInZip, file); err != nil {
		return fmt.Errorf("failed to Copy file into archive: %w", err)
	}
	return nil
}

func (ar *Handler) AppendFilesIntoArchive(filePaths []string, closeAfterFinish bool) error {
	log.Debugf("Append files %s to archive, (close: %v)", strings.Join(filePaths, ","), closeAfterFinish)
	if closeAfterFinish {
		defer func() {
			_ = ar.Close()
		}()
	}

	for _, filePath := range filePaths {
		err := ar.AppendFileToArchive(filePath, filePath)
		if err != nil {
			return fmt.Errorf("failed to write logfiles into archive: %w", err)
		}
	}

	return nil
}

func (ar *Handler) Close() error {
	log.Debugf("Close archive")
	err := ar.writer.Close()
	if err != nil {
		return fmt.Errorf("could not close archive file: %w", err)
	}
	return nil
}
