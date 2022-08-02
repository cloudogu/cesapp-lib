package log

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
type defaultFileHandler struct{}

func (d *defaultFileHandler) open(filename string) (*os.File, error) {
	return os.Open(filename)
}

func (d *defaultFileHandler) create(filename string) (*os.File, error) {
	return os.Create(filename)
}

type ArchiveHandler interface {
	Open(path string) (io.Writer, error)
	AppendFileToArchive(path string, writer zip.Writer) error
	Close() error
}

// SupportArchiveHandler
// The normal procedure should look like this.
// 		1. CreateZipArchiveFile
// 		2. InitialiseZipWriter
// 		3. AppendFileToArchive (n times)
// 		4. Close
type SupportArchiveHandler struct {
	writer      *zip.Writer
	fileCreator fileCreator
	fileOpener  fileOpener
}

func New() *SupportArchiveHandler {
	return &SupportArchiveHandler{
		fileCreator: &defaultFileHandler{},
		fileOpener:  &defaultFileHandler{},
	}
}

// CreateZipArchiveFile creates the file that will be the zip archive.
// The zipFilePath expects a complete path with the correct file extension (.zip).
// If you not intend to create an io.Writer beforehand this method can be the input of InitialiseZipWriter.
func (ar *SupportArchiveHandler) CreateZipArchiveFile(zipFilePath string) (io.Writer, error) {
	zippedArchiveFile, err := ar.fileCreator.create(zipFilePath)
	if err != nil {
		return nil, err
	}
	return zippedArchiveFile, nil
}

// InitialiseZipWriter takes an existing io.Writer and initializes a zip.Writer based on it.
func (ar *SupportArchiveHandler) InitialiseZipWriter(zipFile io.Writer) {
	zipWriter := zip.NewWriter(zipFile)
	ar.writer = zipWriter
}

// AppendFileToArchive takes a path to file that is read and appended to an archive.
// make sure to call the Close method when you're done with appending files to the archive.
func (ar *SupportArchiveHandler) AppendFileToArchive(fileToZipPath string, filepathInZip string) error {
	file, err := ar.fileOpener.open(fileToZipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInZip, err := ar.writer.Create(filepathInZip)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fileInZip, file); err != nil {
		return err
	}
	return nil
}

func (ar *SupportArchiveHandler) Close() error {
	err := ar.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ar *SupportArchiveHandler) WriteLogFilesIntoArchive(filePaths []string, closeAfterFinish bool) {
	for _, filePath := range filePaths {
		ar.WriteLogFileIntoArchive(filePath)
	}
	if closeAfterFinish {
		defer func() {
			err := ar.Close()
			if err != nil {
				panic(err)
			}
		}()
	}
}

// WriteLogFileIntoArchive Takes the path to a single logfile and write it to an initialized and created zip-archive.
// The zipped file's dir structure matches the on the real filesystem.
func (ar *SupportArchiveHandler) WriteLogFileIntoArchive(filePath string) error {

	fmt.Printf("opening file: %s\n", filePath)
	doguLogFile, err := SelectLogFile(filePath)
	if err != nil {
		panic(err)
	}

	defer doguLogFile.file.Close()

	createdFileInZip, err := ar.writer.Create(filePath)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(createdFileInZip, doguLogFile.file); err != nil {
		panic(err)
	}

	return nil
}
