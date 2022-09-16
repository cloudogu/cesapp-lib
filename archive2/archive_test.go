package archive2

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/archive2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
	"time"
)

// testArchiveByteArray contains the bytes of a test archive
// The test archive contains:
// * File "myfile" containing text "test"
// * File "myfile1" containing text "test1"
// * File "myfile2" containing text "test2"
var testArchiveByteArray = []byte{80, 75, 3, 4, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 9, 0, 109, 121, 102, 105, 108, 101, 85, 84, 5, 0, 1, 205, 167, 207, 97, 42, 73, 45, 46, 1, 4, 0, 0, 255, 255, 80, 75, 7, 8, 12, 126, 127, 216, 10, 0, 0, 0, 4, 0, 0, 0, 80, 75, 3, 4, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 9, 0, 109, 121, 102, 105, 108, 101, 49, 85, 84, 5, 0, 1, 205, 167, 207, 97, 42, 73, 45, 46, 49, 4, 4, 0, 0, 255, 255, 80, 75, 7, 8, 226, 220, 178, 138, 11, 0, 0, 0, 5, 0, 0, 0, 80, 75, 3, 4, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 9, 0, 109, 121, 102, 105, 108, 101, 50, 85, 84, 5, 0, 1, 205, 167, 207, 97, 42, 73, 45, 46, 49, 2, 4, 0, 0, 255, 255, 80, 75, 7, 8, 88, 141, 187, 19, 11, 0, 0, 0, 5, 0, 0, 0, 80, 75, 1, 2, 20, 0, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 12, 126, 127, 216, 10, 0, 0, 0, 4, 0, 0, 0, 6, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 102, 105, 108, 101, 85, 84, 5, 0, 1, 205, 167, 207, 97, 80, 75, 1, 2, 20, 0, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 226, 220, 178, 138, 11, 0, 0, 0, 5, 0, 0, 0, 7, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 71, 0, 0, 0, 109, 121, 102, 105, 108, 101, 49, 85, 84, 5, 0, 1, 205, 167, 207, 97, 80, 75, 1, 2, 20, 0, 20, 0, 8, 0, 8, 0, 32, 8, 33, 84, 88, 141, 187, 19, 11, 0, 0, 0, 5, 0, 0, 0, 7, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 144, 0, 0, 0, 109, 121, 102, 105, 108, 101, 50, 85, 84, 5, 0, 1, 205, 167, 207, 97, 80, 75, 5, 6, 0, 0, 0, 0, 3, 0, 3, 0, 185, 0, 0, 0, 217, 0, 0, 0, 0, 0}

func Test_AddContentAsFileWithModifiedDate(t *testing.T) {
	manager := NewManager()

	bytes := manager.GetContent()
	assert.True(t, len(bytes) == 0)

	date := time.Date(2022, 1, 1, 1, 1, 1, 1, &time.Location{})

	err := manager.AddContentAsFileWithModifiedDate("test", "myfile", date)
	require.NoError(t, err)

	err = manager.AddContentAsFileWithModifiedDate("test1", "myfile1", date)
	require.NoError(t, err)

	err = manager.AddContentAsFileWithModifiedDate("test2", "myfile2", date)
	require.NoError(t, err)

	err = manager.Close()
	require.NoError(t, err)

	bytes = manager.GetContent()
	require.True(t, len(bytes) > 0)
	require.Equal(t, testArchiveByteArray, bytes)
}

func TestNewManager(t *testing.T) {
	t.Run("can initialize, no field is nil", func(t *testing.T) {
		manager := NewManager()
		require.NotNil(t, manager.close)
		require.NotNil(t, manager.writer)
		require.NotNil(t, manager.save)
		require.NotNil(t, manager.readFile)
		require.NotNil(t, manager.stat)
		require.NotNil(t, manager.writeToZip)
		require.NotNil(t, manager.buffer)
	})
}

func TestGetContent(t *testing.T) {
	t.Run("content is filled", func(t *testing.T) {
		manager := NewManager()

		bytes := manager.GetContent()
		assert.True(t, len(bytes) == 0)

		err := manager.AddContentAsFile("test", "myfile")

		err = manager.Close()
		require.NoError(t, err)

		bytes = manager.GetContent()
		require.True(t, len(bytes) > 0)
	})
}

func TestAddFileToArchive(t *testing.T) {
	t.Run("can add file to archive", func(t *testing.T) {
		manager := NewManager()
		counter := 0
		modTime := time.Date(2022, 1, 1, 1, 1, 1, 1, &time.Location{})
		fileInfoMock := &mocks.FileInfo{}
		fileInfoMock.On("Size").Return(int64(0))
		fileInfoMock.On("Name").Return(fmt.Sprintf("filename-%v", counter))
		fileInfoMock.On("ModTime").Return(modTime)
		fileInfoMock.On("Mode").Return(fs.FileMode(0755))
		manager.readFile = func(name string) (result []byte, err error) {
			result, err = []byte(fmt.Sprintf("content-%v", counter)), nil
			counter++
			return
		}
		manager.stat = func(name string) (os.FileInfo, error) {
			if name == "myfileoutside1" || name == "myfileoutside2" {
				return fileInfoMock, nil
			}
			return nil, errors.New(name)
		}
		manager.writeToZip = func(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error {
			strContent := string(content)
			if strContent == "content-1" || strContent == "content-2" {
				return nil
			}
			return errors.New(string(content))
		}

		err := manager.AddFileToArchive(File{
			NameOutside: "myfileoutside1",
			NameInside:  "myfileinside1",
		})

		err = manager.AddFileToArchive(File{
			NameOutside: "myfileoutside2",
			NameInside:  "myfileinside2",
		})
		require.NoError(t, err)

		fileInfoMock.AssertExpectations(t)
	})
	t.Run("can add multiple file to archive", func(t *testing.T) {
		manager := NewManager()
		counter := 0
		modTime := time.Date(2022, 1, 1, 1, 1, 1, 1, &time.Location{})
		fileInfoMock := &mocks.FileInfo{}
		fileInfoMock.On("Size").Return(int64(0))
		fileInfoMock.On("Name").Return(fmt.Sprintf("filename-%v", counter))
		fileInfoMock.On("ModTime").Return(modTime)
		fileInfoMock.On("Mode").Return(fs.FileMode(0755))
		manager.readFile = func(name string) (result []byte, err error) {
			result, err = []byte(fmt.Sprintf("content-%v", counter)), nil
			counter++
			return
		}
		manager.stat = func(name string) (os.FileInfo, error) {
			if name == "myfileoutside1" || name == "myfileoutside2" {
				return fileInfoMock, nil
			}
			return nil, errors.New(name)
		}
		manager.writeToZip = func(zipWriter *zip.Writer, header *zip.FileHeader, content []byte) error {
			strContent := string(content)
			if strContent == "content-0" || strContent == "content-1" {
				return nil
			}
			return errors.New(string(content))
		}

		t.Run("with close", func(t *testing.T) {
			err := manager.AddFilesToArchive([]File{
				{
					NameOutside: "myfileoutside1",
					NameInside:  "myfileinside1",
				},
				{
					NameOutside: "myfileoutside2",
					NameInside:  "myfileinside2",
				},
			}, true)
			require.NoError(t, err)
		})

		counter = 0

		t.Run("without close", func(t *testing.T) {
			err := manager.AddFilesToArchive([]File{
				{
					NameOutside: "myfileoutside1",
					NameInside:  "myfileinside1",
				},
				{
					NameOutside: "myfileoutside2",
					NameInside:  "myfileinside2",
				},
			}, false)
			require.NoError(t, err)
		})

		fileInfoMock.AssertExpectations(t)
	})
	t.Run("fail on add multiple files to archive", func(t *testing.T) {
		manager := NewManager()
		manager.readFile = func(name string) (result []byte, err error) {
			return nil, errors.New("testerror")
		}

		err := manager.AddFilesToArchive([]File{
			{
				NameOutside: "myfileoutside1",
				NameInside:  "myfileinside1",
			},
			{
				NameOutside: "myfileoutside2",
				NameInside:  "myfileinside2",
			},
		}, false)
		require.Error(t, err)
		require.Equal(t, "testerror", err.Error())
	})
	t.Run("fail on read file", func(t *testing.T) {
		manager := NewManager()
		manager.readFile = func(name string) (result []byte, err error) {
			return nil, errors.New("testerror")
		}

		err := manager.AddFilesToArchive([]File{
			{
				NameOutside: "myfileoutside1",
				NameInside:  "myfileinside1",
			},
			{
				NameOutside: "myfileoutside2",
				NameInside:  "myfileinside2",
			},
		}, false)
		require.Error(t, err)
		require.Equal(t, "testerror", err.Error())
	})
	t.Run("fail on stat", func(t *testing.T) {
		manager := NewManager()
		manager.readFile = func(name string) (result []byte, err error) {
			return nil, nil
		}
		manager.stat = func(name string) (os.FileInfo, error) {
			return nil, errors.New("testerror")
		}

		err := manager.AddFilesToArchive([]File{
			{
				NameOutside: "myfileoutside1",
				NameInside:  "myfileinside1",
			},
			{
				NameOutside: "myfileoutside2",
				NameInside:  "myfileinside2",
			},
		}, false)
		require.Error(t, err)
		require.Equal(t, "testerror", err.Error())
	})
}
