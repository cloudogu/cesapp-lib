package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type fileCreator interface {
	create(filename string) (*os.File, error)
}
type fileOpener interface {
	open(filename string) (*os.File, error)
}
type fileCopier interface {
	copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type zipWriter interface {
	Create(name string) (io.Writer, error)
	Close() error
}

type defaultFileHandler struct{}

func (d *defaultFileHandler) open(filename string) (*os.File, error) {
	return os.Open(filename)
}

func (d *defaultFileHandler) create(filename string) (*os.File, error) {
	return os.Create(filename)
}

func (d *defaultFileHandler) copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

type Handler interface {
	CreateZipArchiveFile(zipFilePath string) (io.Writer, error)
	InitialiseZipWriter(zipFile io.Writer)
	AppendFileToArchive(fileToZipPath string, filepathInZip string) error
	Close() error
	WriteFilesIntoArchive(filePaths []string, closeAfterFinish bool) error
}

// DefaultHandler
// The normal procedure should look like this.
// 		1. CreateZipArchiveFile
// 		2. InitialiseZipWriter
// 		3. AppendFileToArchive (n times)
// 		4. Close
type DefaultHandler struct {
	writer      zipWriter
	fileCreator fileCreator
	fileOpener  fileOpener
	fileCopier  fileCopier
}

func NewHandler() Handler {
	return &DefaultHandler{
		fileCreator: &defaultFileHandler{},
		fileOpener:  &defaultFileHandler{},
		fileCopier:  &defaultFileHandler{},
	}
}

func InitArchiveInPath(filePath string) (Handler, error) {
	supportArchiveHandler := NewHandler()
	supportArchive, err := supportArchiveHandler.CreateZipArchiveFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create support archive file: %w", err)
	}

	supportArchiveHandler.InitialiseZipWriter(supportArchive)

	return supportArchiveHandler, nil
}

// CreateZipArchiveFile creates the file that will be the zip archive.
// The zipFilePath expects a complete path with the correct file extension (.zip).
// If you not intend to create an io.Writer beforehand this method can be the input of InitialiseZipWriter.
func (ar *DefaultHandler) CreateZipArchiveFile(zipFilePath string) (io.Writer, error) {
	zippedArchiveFile, err := ar.fileCreator.create(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive file: %w", err)
	}
	return zippedArchiveFile, nil
}

// InitialiseZipWriter takes an existing io.Writer and initializes a zip.Writer based on it.
func (ar *DefaultHandler) InitialiseZipWriter(zipFile io.Writer) {
	zipWriter := zip.NewWriter(zipFile)
	ar.writer = zipWriter
}

// AppendFileToArchive takes a path to file that is read and appended to an archive.
// make sure to call the Close method when you're done with appending files to the archive.
func (ar *DefaultHandler) AppendFileToArchive(fileToZipPath string, filepathInZip string) error {
	file, err := ar.fileOpener.open(fileToZipPath)
	if err != nil {
		return fmt.Errorf("failed to read base file for appending to archive: %w", err)
	}
	defer file.Close()

	fileInZip, err := ar.writer.Create(filepathInZip)
	if err != nil {
		return fmt.Errorf("failed to create file in archive: %w", err)
	}

	if _, err := ar.fileCopier.copy(fileInZip, file); err != nil {
		return fmt.Errorf("failed to copy file into archive: %w", err)
	}
	return nil
}

func (ar *DefaultHandler) Close() error {
	err := ar.writer.Close()
	if err != nil {
		return fmt.Errorf("could not close archive file: %w", err)
	}
	return nil
}

func (ar *DefaultHandler) WriteFilesIntoArchive(filePaths []string, closeAfterFinish bool) error {
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
