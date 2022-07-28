package log

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func teardown() {
	os.Remove(supportArchiveFileName)
}

func TestWriteSupportArchive(t *testing.T) {
	logfiles := []string{"resources/test/testlogfile.log"}
	err := WriteSupportArchive(logfiles)
	if err != nil {
		assert.NotNil(t, err)
	}
	teardown()
}
