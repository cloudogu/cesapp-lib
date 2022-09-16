package archive2_test

import (
	"github.com/cloudogu/cesapp-lib/archive2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func Test_SaveArchiveAsFile(t *testing.T) {
	t.Run("writes content", func(t *testing.T) {
		manager := archive2.NewManager()

		bufferContent := manager.GetContent()
		assert.True(t, len(bufferContent) == 0)

		err := manager.AddContentAsFile("test", "myfile")
		require.NoError(t, err)

		err = manager.AddContentAsFile("test1", "myfile1")
		require.NoError(t, err)

		err = manager.AddContentAsFile("test2", "myfile2")
		require.NoError(t, err)

		err = manager.Close()
		require.NoError(t, err)

		bufferContent = manager.GetContent()
		require.True(t, len(bufferContent) > 0)

		err = manager.SaveArchiveAsFile("./archive.zip")
		defer func() {
			_ = recover()
			_ = os.Remove("./archive.zip")
		}()
		require.NoError(t, err)

		fileContent, err := os.ReadFile("./archive.zip")
		require.NoError(t, err)

		require.Equal(t, bufferContent, fileContent)
	})

	t.Run("corrects path", func(t *testing.T) {
		manager := archive2.NewManager()

		err := os.Mkdir("./testdata", 0755)
		require.NoError(t, err)

		err = manager.SaveArchiveAsFile("./testdata/")
		defer func() {
			_ = recover()
			_ = os.Remove("./testdata/archive.zip")
			_ = os.Remove("./testdata")
		}()
		require.NoError(t, err)

		entries, err := os.ReadDir("./testdata")
		require.NoError(t, err)
		require.Len(t, entries, 1)

		require.True(t, strings.HasSuffix(entries[0].Name(), "archive.zip"))
	})

}
