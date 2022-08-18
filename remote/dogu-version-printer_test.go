package remote

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	remocks "github.com/cloudogu/cesapp-lib/remote/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimitVersions_inttest(t *testing.T) {
	latest := core.Version{
		Raw: "0.3",
	}
	older := core.Version{
		Raw: "0.2",
	}
	evenOlder := core.Version{
		Raw: "0.1",
	}
	versions := []core.Version{latest, older, evenOlder}

	t.Run("no change when limit is higher then length", func(t *testing.T) {
		limited := limitVersions(versions, 5)
		assert.Len(t, limited, 3)
		assert.Equal(t, latest, limited[0])
		assert.Equal(t, older, limited[1])
		assert.Equal(t, evenOlder, limited[2])
		assert.Equal(t, versions, limited)
	})
	t.Run("no change when limit is equal then length", func(t *testing.T) {
		limited := limitVersions(versions, 3)
		assert.Len(t, limited, 3)
		assert.Equal(t, latest, limited[0])
		assert.Equal(t, older, limited[1])
		assert.Equal(t, evenOlder, limited[2])
		assert.Equal(t, versions, limited)
	})
	t.Run("no change when limit is lower then length", func(t *testing.T) {
		limited := limitVersions(versions, 2)
		assert.Len(t, limited, 2)
		assert.Equal(t, latest, limited[0])
		assert.Equal(t, older, limited[1])
	})
	t.Run("no change when limit is zero", func(t *testing.T) {
		limited := limitVersions(versions, 0)
		assert.Len(t, limited, 0)
	})
	t.Run("no change when limit is negative", func(t *testing.T) {
		limited := limitVersions(versions, -1)
		assert.Len(t, limited, 0)
	})
}

func TestDoguVersionPrinter_PrintForAllDogus_inttest(t *testing.T) {
	realStdOut := os.Stdout
	//core.GetLogger().InitForUnitTests(logrus.InfoLevel)

	t.Run("should fail for a limit less than 0", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		// when
		err := sut.PrintForAllDogus(-1)

		// then
		require.Error(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Equal(t, "", actualOutput)
		mockRemote.AssertExpectations(t)
	})
	t.Run("should print all versions for a limit equal to 0", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		dogu1 := &core.Dogu{
			Name: "namespace/dogu1",
		}
		mockRemote := &remocks.Registry{}
		mockRemote.On("GetAll").Return([]*core.Dogu{dogu1}, nil)
		mockRemote.On("GetVersionsOf", dogu1.Name).Return([]core.Version{{Raw: "v1"}, {Raw: "v2"}}, nil)

		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintForAllDogus(0)

		// then
		require.NoError(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Contains(t, actualOutput, "level=info msg=\"  namespace/dogu1:\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v1\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v2\\n\"")
		mockRemote.AssertExpectations(t)
	})

	t.Run("should print only 1 of 2 versions for a limit equal to 1", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		dogu1 := &core.Dogu{
			Name: "namespace/dogu1",
		}
		mockRemote := &remocks.Registry{}
		mockRemote.On("GetAll").Return([]*core.Dogu{dogu1}, nil)
		mockRemote.On("GetVersionsOf", dogu1.Name).Return([]core.Version{{Raw: "v1"}, {Raw: "v2"}}, nil)

		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintForAllDogus(1)

		// then
		require.NoError(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Contains(t, actualOutput, "level=info msg=\"  namespace/dogu1:\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v1\\n\"")
		mockRemote.AssertExpectations(t)
	})
}

func TestDoguVersionPrinter_PrintForSingleDogu_inttest(t *testing.T) {
	realStdOut := os.Stdout
	//logging.InitForUnitTests(logrus.InfoLevel)

	dogu1 := &core.Dogu{
		Name: "namespace/dogu1",
	}

	t.Run("should fail for a limit less than 0", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		// when
		err := sut.PrintForSingleDogu(dogu1, -1)

		// then
		require.Error(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Equal(t, "", actualOutput)
		mockRemote.AssertExpectations(t)
	})
	t.Run("should print all dogu versions for a limit equal to 0", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		mockRemote.On("GetVersionsOf", dogu1.Name).Return([]core.Version{{Raw: "v1"}, {Raw: "v2"}}, nil)
		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintForSingleDogu(dogu1, 0)

		// then
		require.NoError(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Contains(t, actualOutput, "level=info msg=\"  namespace/dogu1:\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v1\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v2\\n\"")
		mockRemote.AssertExpectations(t)
	})
	t.Run("should print 1 of 2 dogu version for a limit equal to 1", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		mockRemote.On("GetVersionsOf", dogu1.Name).Return([]core.Version{{Raw: "v1"}, {Raw: "v2"}}, nil)
		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()

		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintForSingleDogu(dogu1, 1)

		// then
		require.NoError(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Contains(t, actualOutput, "level=info msg=\"  namespace/dogu1:\\n\"")
		assert.Contains(t, actualOutput, "level=info msg=\"     - v1\\n\"")
		mockRemote.AssertExpectations(t)
	})
}

func TestDoguVersionPrinter_PrintInDefaultFormat_inttest(t *testing.T) {
	realStdOut := os.Stdout
	//logging.InitForUnitTests(logrus.InfoLevel)

	dogu1 := &core.Dogu{
		Name:        "namespace/dogu1",
		Description: "dogu1desc",
		Version:     "dogu1version",
		DisplayName: "dogu1displayname",
	}
	dogu2 := &core.Dogu{
		Name:        "namespace/dogu2",
		Description: "dogu2desc",
		Version:     "dogu2version",
		DisplayName: "dogu2displayname",
	}

	t.Run("should print table in default format", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		mockRemote.On("GetAll").Return([]*core.Dogu{dogu1, dogu2}, nil)
		sut := &DoguVersionPrinter{Remote: mockRemote}

		fakeReader, fakeWriter := routeStdoutToReplacement()
		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintDoguListInDefaultFormat()

		// then
		assert.Nil(t, err)
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.Contains(t, actualOutput, "NAME")
		assert.Contains(t, actualOutput, "VERSION")
		assert.Contains(t, actualOutput, "DESCRIPTION")
		assert.Contains(t, actualOutput, "DISPLAYNAME")
		assert.Contains(t, actualOutput, dogu1.DisplayName)
		assert.Contains(t, actualOutput, dogu1.Description)
		assert.Contains(t, actualOutput, dogu1.Version)
		assert.Contains(t, actualOutput, dogu1.Name)
		assert.Contains(t, actualOutput, dogu2.DisplayName)
		assert.Contains(t, actualOutput, dogu2.Description)
		assert.Contains(t, actualOutput, dogu2.Version)
		assert.Contains(t, actualOutput, dogu2.Name)
		mockRemote.AssertExpectations(t)
	})

	t.Run("error on getAll", func(t *testing.T) {
		defer restoreOriginalStdout(realStdOut)

		mockRemote := &remocks.Registry{}
		mockRemote.On("GetAll").Return([]*core.Dogu{dogu1, dogu2}, assert.AnError)
		sut := &DoguVersionPrinter{Remote: mockRemote}
		fakeReader, fakeWriter := routeStdoutToReplacement()
		originalLogger := core.GetLogger
		defer func() { core.GetLogger = originalLogger }()
		myLogger := logrus.New()
		myLogger.Out = fakeWriter
		core.GetLogger = func() logrus.FieldLogger {
			return myLogger
		}

		// when
		err := sut.PrintDoguListInDefaultFormat()

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), assert.AnError.Error())
		actualOutput := captureOutput(fakeReader, fakeWriter, realStdOut)
		assert.NotContains(t, actualOutput, "NAME")
		assert.NotContains(t, actualOutput, "VERSION")
		assert.NotContains(t, actualOutput, "DESCRIPTION")
		assert.NotContains(t, actualOutput, "DISPLAYNAME")
		assert.NotContains(t, actualOutput, dogu1.DisplayName)
		assert.NotContains(t, actualOutput, dogu1.Description)
		assert.NotContains(t, actualOutput, dogu1.Version)
		assert.NotContains(t, actualOutput, dogu1.Name)
		assert.NotContains(t, actualOutput, dogu2.DisplayName)
		assert.NotContains(t, actualOutput, dogu2.Description)
		assert.NotContains(t, actualOutput, dogu2.Version)
		assert.NotContains(t, actualOutput, dogu2.Name)
		mockRemote.AssertExpectations(t)
	})
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}

func routeStdoutToReplacement() (readerPipe, writerPipe *os.File) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return r, w
}

func captureOutput(fakeReaderPipe, fakeWriterPipe, originalStdout *os.File) string {
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
