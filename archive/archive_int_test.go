package archive

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitInPath_inttest(t *testing.T) {
	archive, err := InitIn("./archive.zip")
	defer func() {
		_ = os.Remove("./archive.zip")
	}()
	require.NoError(t, err)
	err = archive.Close()

	file, err := os.ReadFile("./archive.zip")
	require.NoError(t, err)
	require.NotNil(t, file)
}

func TestSupportArchiveHandler_WriteFilesIntoArchive_inttest(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "*.zip")
	require.NoError(t, err)
	assert.NotNil(t, zipFile)

	handler, err := InitIn(zipFile.Name())
	defer func() {
		_ = os.Remove(zipFile.Name())
	}()
	require.NoError(t, err)

	tmpLogFile1, err := ioutil.TempFile("", "*.archive")
	if err != nil {
		fmt.Println("Failed to Create temp zipFile: ", tmpLogFile1.Name())
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
		fmt.Println("Failed to Create temp zipFile: ", tmpLogFile2.Name())
		t.Fail()
	}
	fmt.Println("Created temp zipFile: ", tmpLogFile2.Name())
	defer tmpLogFile2.Close()
	if _, err := tmpLogFile2.WriteString("test data"); err != nil {
		fmt.Print("Unable to write to temporary file")
		t.Fail()
	}

	logFiles := []string{tmpLogFile1.Name(), tmpLogFile2.Name()}

	handler.AppendFilesIntoArchive(logFiles, true)

	assert.FileExists(t, zipFile.Name())
	fi, err := zipFile.Stat()
	if err != nil {
		fmt.Print("Unable to get info about temporary zip file")
		t.Fail()
	}
	assert.True(t, 300 < fi.Size()) // empty zip archive is around ~20, with two logs around ~300
}
