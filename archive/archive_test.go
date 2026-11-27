package archive

import (
	"archive/zip"
	bytesPkg "bytes"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/archive/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"
)

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

	// compare the archive entries instead of the raw bytes, as the compressed output may differ between go versions
	reader, err := zip.NewReader(bytesPkg.NewReader(bytes), int64(len(bytes)))
	require.NoError(t, err)
	expected := []struct{ name, content string }{
		{"myfile", "test"},
		{"myfile1", "test1"},
		{"myfile2", "test2"},
	}
	require.Len(t, reader.File, len(expected))
	for i, file := range reader.File {
		assert.Equal(t, expected[i].name, file.Name)
		assert.True(t, date.Truncate(time.Second).Equal(file.Modified), "unexpected modified date %s", file.Modified)

		rc, err := file.Open()
		require.NoError(t, err)
		content, err := io.ReadAll(rc)
		require.NoError(t, err)
		require.NoError(t, rc.Close())
		assert.Equal(t, expected[i].content, string(content))
	}
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
