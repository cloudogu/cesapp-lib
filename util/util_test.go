package util

import (
	"bytes"
	"github.com/cloudogu/cesapp-lib/core"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

type Sample struct {
	Name string
}

func TestReadJSONFile(t *testing.T) {
	s := Sample{}
	err := ReadJSONFile(&s, "../resources/test/sample.json")
	assert.Nil(t, err)

	assert.Equal(t, "sample-123", s.Name, "name should be equal")
}

func TestGetContentOfFile(t *testing.T) {
	content, err := GetContentOfFile("../resources/test/sample.json")
	require.Nil(t, err)
	assert.NotNil(t, content)

	assert.Contains(t, content, "sample-123")
}

func TestReadJSONFileThatDoesNotExists(t *testing.T) {
	s := Sample{}
	err := ReadJSONFile(&s, "../resources/test/sample.jsonnn")
	assert.NotNil(t, err)
}

func TestReadJSONFileWithInvalidStructure(t *testing.T) {
	s := Sample{}
	err := ReadJSONFile(&s, "../resources/test/invalid.json")
	assert.NotNil(t, err)
}

func TestWriteJSONFile(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
	path := file.Name()
	defer os.Remove(path)

	s := Sample{"123-sample"}
	err := WriteJSONFile(s, path)
	assert.Nil(t, err)

	o := Sample{}
	err = ReadJSONFile(&o, path)
	assert.Nil(t, err)
	assert.Equal(t, "123-sample", o.Name, "name should be equal")
}

func TestExists(t *testing.T) {
	assert.False(t, Exists("../Makefilell"), "should return false")
	assert.True(t, Exists("../Makefile"), "should return true")
}

func initPackageLoggerWithDefaultLogger() func() logrus.FieldLogger {
	originalLogger := core.GetLogger()

	newLogger := logrus.New()
	newLogger.SetLevel(logrus.DebugLevel)
	core.GetLogger = func() logrus.FieldLogger { return newLogger }

	resetOriginalLogger := func() logrus.FieldLogger { return originalLogger }

	return resetOriginalLogger
}

func TestClose(t *testing.T) {
	t.Run("should close", func(t *testing.T) {
		closer := &MockCloser{}

		CloseButLogError(closer, "shouldCloseDisregardLoggingHere")

		assert.True(t, closer.wasCloseCalled)
	})
}

func TestCloseButLogError_shouldCloseWithoutLog(t *testing.T) {
	if os.Getenv("SKIP_SYSLOG_TESTS") != "" {
		t.Skip("Skipping syslog test. This test must be called in a Debian/Ubuntu environment with syslog support")
	}

	realStdout := os.Stdout
	defer restoreOriginalStdout(realStdout)
	fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

	//originalGetLogger := initPackageLoggerWithDefaultLogger()
	//defer func() {
	//	core.GetLogger = originalGetLogger
	//}()
	closer := &MockCloser{shouldErr: false}

	// when
	CloseButLogError(closer, "shouldCloseWithoutLog")

	// then
	actualOutput := captureOutputAndRestoreStdout(fakeReaderPipe, fakeWriterPipe, realStdout)
	assert.Empty(t, actualOutput)
	assert.True(t, closer.wasCloseCalled)
}

func TestCloseButLogError_shouldNotCloseAndLogInstead(t *testing.T) {
	if os.Getenv("SKIP_SYSLOG_TESTS") != "" {
		t.Skip("Skipping syslog test. This test must be called in a Debian/Ubuntu environment with syslog support")
	}

	realStdout := os.Stdout
	defer restoreOriginalStdout(realStdout)
	fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

	originalLogger := core.GetLogger
	defer func() { core.GetLogger = originalLogger }()
	myLogger := logrus.New()
	myLogger.Out = fakeWriterPipe
	core.GetLogger = func() logrus.FieldLogger {
		return myLogger
	}

	closer := &MockCloser{shouldErr: true}

	// when
	CloseButLogError(closer, "closingResultedInErrorWithLog")

	// then
	actualOutput := captureOutputAndRestoreStdout(fakeReaderPipe, fakeWriterPipe, realStdout)

	assert.NotEmpty(t, actualOutput)
	assert.Contains(t, actualOutput, "error")
	assert.Contains(t, actualOutput, "closingResultedInErrorWithLog")
	assert.Contains(t, actualOutput, "oh no, an IO error")
	assert.True(t, closer.wasCloseCalled)
}

func TestContains(t *testing.T) {
	slice := []string{"hello", "Kitty", "world"}

	assert.True(t, Contains(slice, "hello"))
	assert.True(t, Contains(slice, "Kitty"))
	assert.True(t, Contains(slice, "world"))
	assert.False(t, Contains(slice, "world1"))
}

type MockCloser struct {
	shouldErr      bool
	wasCloseCalled bool
}

func (closer *MockCloser) Close() error {
	closer.wasCloseCalled = true

	if closer.shouldErr {
		return errors.New("oh no, an IO error")
	}

	return nil
}

func TestReverse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{"empty string", args{s: ""}, ""},
		{"even", args{s: "asdf"}, "fdsa"},
		{"uneven", args{s: "asd"}, "dsa"},
		{"mixed with digits", args{s: "abcdef123456"}, "654321fedcba"},
		{"ASCII special chars", args{s: "abcde...Space> <...Tab>\t<...!§$%"}, "%$§!...<\t>baT...< >ecapS...edcba"},
		{"UTF-8 special chars", args{s: "abcdeÄÖÜ👍"}, "👍ÜÖÄedcba"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Reverse(tt.args.s); gotResult != tt.wantResult {
				t.Errorf("Reverse() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}

func routeStdoutToReplacement() (readerPipe, writerPipe *os.File) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return r, w
}

func captureOutputAndRestoreStdout(fakeReaderPipe, fakeWriterPipe, originalStdout *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, fakeReaderPipe)
		outC <- buf.String()
	}()

	// back to normal state
	fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}
