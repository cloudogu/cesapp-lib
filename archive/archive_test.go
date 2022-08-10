package archive

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

type mockFileCreator struct{}
type mockFailedFileCreator struct{}
type mockFailedFileOpener struct{}
type mockFailedFileCopier struct{}
type mockZipWriter struct{}

func (mfc *mockFileCreator) create(filename string) (*os.File, error) {
	return ioutil.TempFile("", filename)
}

func (mfc *mockFailedFileCreator) create(filename string) (*os.File, error) {
	return nil, errors.New("failed to create file")
}

func (mfc *mockFailedFileOpener) open(filename string) (*os.File, error) {
	return nil, errors.New("failed to open file")
}

func (mzw *mockZipWriter) Close() error {
	return errors.New("failed to close file")
}

func (mzw *mockZipWriter) Create(name string) (io.Writer, error) {
	return nil, errors.New("failed to create file with writer")
}

func (mzw *mockFailedFileCopier) copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return 0, errors.New("failed to copy file in zip archive")
}

func TestSupportArchiveHandler_CreateZipArchiveFile(t *testing.T) {
	//test success
	handler := DefaultHandler{fileCreator: &mockFileCreator{}}
	file, err := handler.CreateZipArchiveFile("test.zip")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	// test failure
	handler = DefaultHandler{fileCreator: &mockFailedFileCreator{}}
	_, err = handler.CreateZipArchiveFile("test.zip")
	require.Error(t, err)
}

func TestSupportArchiveHandler_New(t *testing.T) {
	supportArchiveHandler := NewHandler()
	assert.NotNil(t, supportArchiveHandler)
}

func TestSupportArchiveHandler_Close(t *testing.T) {
	handler := DefaultHandler{
		writer: &mockZipWriter{},
	}

	err := handler.Close()
	assert.Error(t, err)

	zipFile, _ := ioutil.TempFile("", "*.zip")
	handler.InitialiseZipWriter(zipFile)
	err = handler.Close()
	assert.NoError(t, err)
}

func TestSupportArchiveHandler_InitialiseZipWriter(t *testing.T) {
	handler := DefaultHandler{fileCreator: &mockFileCreator{}}
	file, _ := handler.CreateZipArchiveFile("test.zip")
	handler.InitialiseZipWriter(file)

	assert.NotNil(t, handler.writer)
}

func TestSupportArchiveHandler_AppendFileToArchive_Success(t *testing.T) {
	handler := DefaultHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &defaultFileHandler{},
		fileCopier:  &defaultFileHandler{},
	}
	zipFile, _ := ioutil.TempFile("", "*.zip")
	handler.InitialiseZipWriter(zipFile)

	tmpFile, err := ioutil.TempFile("", "test.txt")
	if err != nil {
		fmt.Println("Failed to create temp zipFile: ", tmpFile.Name())
		t.Fail()
	}
	fmt.Println("Created temp zipFile: ", tmpFile.Name())
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString("test data"); err != nil {
		fmt.Print("Unable to write to temporary file")
		t.Fail()
	}
	handler.AppendFileToArchive(tmpFile.Name(), "/test.txt")
	handler.Close()

	assert.FileExists(t, zipFile.Name())
	fi, err := zipFile.Stat()
	if err != nil {
		fmt.Print("Unable to get info about temporary zip file")
		t.Fail()
	}
	assert.True(t, fi.Size() > 25) // empty zip archive is around ~20
}

func TestSupportArchiveHandler_AppendFileToArchive_Failure(t *testing.T) {
	// first point of failure: Cannot open file
	handler := DefaultHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &mockFailedFileOpener{},
	}

	tmpFile, _ := ioutil.TempFile("", "test.txt")
	err := handler.AppendFileToArchive(tmpFile.Name(), "/test.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")

	// second point of failure: Cannot open file
	handler = DefaultHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &defaultFileHandler{},
		writer:      &mockZipWriter{},
	}
	err = handler.AppendFileToArchive(tmpFile.Name(), "/test.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create file with writer")

	// third point of failure: Cannot copy file into zip archive
	handler = DefaultHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &defaultFileHandler{},
		fileCopier:  &mockFailedFileCopier{},
	}
	zipFile, _ := ioutil.TempFile("", "*.zip")
	handler.InitialiseZipWriter(zipFile)
	err = handler.AppendFileToArchive(tmpFile.Name(), "/test.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to copy file in zip archive")

}

//func TestSupportArchiveHandler_WriteLogFileIntoArchive_Success(t *testing.T) {
//	handler := DefaultHandler{
//		fileCreator: &mockFileCreator{},
//		fileOpener:  &defaultFileHandler{},
//		fileCopier:  &defaultFileHandler{},
//	}
//
//	zipFile, _ := ioutil.TempFile("", "*.zip")
//	assert.NotNil(t, zipFile)
//
//	handler.InitialiseZipWriter(zipFile)
//
//	tmpLogFile, _ := ioutil.TempFile("", "*.archive")
//	assert.NotNil(t, tmpLogFile)
//	defer tmpLogFile.Close()
//
//	_, err := tmpLogFile.WriteString("test entry in logfile data")
//	assert.NoError(t, err)
//
//	handler.WriteFilesIntoArchive(tmpLogFile.Name())
//	handler.Close()
//
//	assert.FileExists(t, zipFile.Name())
//	fi, err := zipFile.Stat()
//	if err != nil {
//		fmt.Print("Unable to get info about temporary zip file")
//		t.Fail()
//	}
//	assert.True(t, 25 < fi.Size()) // empty zip archive is around ~20
//
//}
//
//func TestSupportArchiveHandler_WriteLogFileIntoArchive_Fail(t *testing.T) {
//	// first point of failure: Cannot open file
//	handler := DefaultHandler{
//		fileCreator: &mockFileCreator{},
//		fileOpener:  &mockFailedFileOpener{},
//	}
//
//	err := handler.WriteLogFileIntoArchive("not_a_valid_file.archive")
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "no such file or directory")
//
//	// second point of failure: Cannot open file
//	handler = DefaultHandler{
//		fileCreator: &mockFileCreator{},
//		fileOpener:  &defaultFileHandler{},
//		writer:      &mockZipWriter{},
//	}
//	tmpFile, _ := ioutil.TempFile("", "test.archive")
//	err = handler.WriteLogFileIntoArchive(tmpFile.Name())
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "failed to create file with writer")
//
//	// third point of failure: Cannot copy file into zip archive
//	handler = DefaultHandler{
//		fileCreator: &mockFileCreator{},
//		fileOpener:  &defaultFileHandler{},
//		fileCopier:  &mockFailedFileCopier{},
//	}
//	zipFile, _ := ioutil.TempFile("", "*.zip")
//	handler.InitialiseZipWriter(zipFile)
//	err = handler.WriteLogFileIntoArchive(tmpFile.Name())
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "failed to copy file in zip archive")
//
//}
//

func TestSupportArchiveHandler_WriteFilesIntoArchive(t *testing.T) {
	handler := DefaultHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &defaultFileHandler{},
		fileCopier:  &defaultFileHandler{},
	}

	zipFile, _ := ioutil.TempFile("", "*.zip")
	assert.NotNil(t, zipFile)

	handler.InitialiseZipWriter(zipFile)

	tmpLogFile1, err := ioutil.TempFile("", "*.archive")
	if err != nil {
		fmt.Println("Failed to create temp zipFile: ", tmpLogFile1.Name())
		t.Fail()
	}
	fmt.Println("Created temp zipFile: ", tmpLogFile1.Name())
	defer tmpLogFile1.Close()
	if _, err := tmpLogFile1.WriteString("test data"); err != nil {
		fmt.Print("Unable to write to temporary file")
		t.Fail()
	}

	tmpLogFile2, err := ioutil.TempFile("", "*.archive")
	if err != nil {
		fmt.Println("Failed to create temp zipFile: ", tmpLogFile2.Name())
		t.Fail()
	}
	fmt.Println("Created temp zipFile: ", tmpLogFile2.Name())
	defer tmpLogFile2.Close()
	if _, err := tmpLogFile2.WriteString("test data"); err != nil {
		fmt.Print("Unable to write to temporary file")
		t.Fail()
	}

	logFiles := []string{tmpLogFile1.Name(), tmpLogFile2.Name()}

	handler.WriteFilesIntoArchive(logFiles, true)

	assert.FileExists(t, zipFile.Name())
	fi, err := zipFile.Stat()
	if err != nil {
		fmt.Print("Unable to get info about temporary zip file")
		t.Fail()
	}
	assert.True(t, 300 < fi.Size()) // empty zip archive is around ~20, with two logs around ~300
}
