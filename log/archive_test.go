package log

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

type mockFileCreator struct{}
type mockFailedFileCreator struct{}

func (mfc *mockFileCreator) create(filename string) (*os.File, error) {
	return ioutil.TempFile("", filename)
}

func (mfc *mockFailedFileCreator) create(filename string) (*os.File, error) {
	return nil, errors.New("failed to create file")
}

func TestSupportArchiveHandler_CreateZipArchiveFile(t *testing.T) {
	//test success
	handler := SupportArchiveHandler{fileCreator: &mockFileCreator{}}
	file, err := handler.CreateZipArchiveFile("test.zip")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	// test failure
	handler = SupportArchiveHandler{fileCreator: &mockFailedFileCreator{}}
	_, err = handler.CreateZipArchiveFile("test.zip")
	require.Error(t, err)
}

func TestSupportArchiveHandler_New(t *testing.T) {
	supportArchiveHandler := New()
	assert.NotNil(t, supportArchiveHandler)
}

func TestSupportArchiveHandler_InitialiseZipWriter(t *testing.T) {
	handler := SupportArchiveHandler{fileCreator: &mockFileCreator{}}
	file, _ := handler.CreateZipArchiveFile("test.zip")
	handler.InitialiseZipWriter(file)

	assert.NotNil(t, handler.writer)
}

func TestSupportArchiveHandler_AppendFileToArchive(t *testing.T) {
	handler := SupportArchiveHandler{
		fileCreator: &mockFileCreator{},
		fileOpener:  &defaultFileHandler{},
	}
	file, _ := handler.CreateZipArchiveFile("test.zip")
	handler.InitialiseZipWriter(file)

	tmpFile, err := ioutil.TempFile("", "test.file")
	if err != nil {
		fmt.Println("Failed to create temp file: ", tmpFile.Name())
		t.Fail()
	}
	fmt.Println("Created temp file: ", tmpFile.Name())
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString("test data"); err != nil {
		fmt.Print("Unable to write to temporary file")
		t.Fail()
	}
	handler.AppendFileToArchive(tmpFile.Name(), "/test.file")
	//TODO prüfen ob zip geschrieben wurde
}
