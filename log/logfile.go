package log

import (
	"fmt"
	"github.com/prometheus/common/log"
	"io"
	"os"
	"path"
	"strings"
)

const (
	// doguLogFilesPath contains the path to the dogu log files.
	doguLogFilesPath = "/var/log/docker/"
	// logFileExtension contains the extension name used for log files.
	logFileExtension = "log"
)

type logFile struct {
	file *os.File
	size int64
}

func SelectLogFile(logfilePath string) (*logFile, error) {
	return openDoguLogFileFor(logfilePath)
}

func openDoguLogFileFor(logfilePath string) (*logFile, error) {

	fileInfo, err := os.Stat(logfilePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logfilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &logFile{file: file, size: fileInfo.Size()}, nil
}

func (lf *logFile) close() {
	if lf.file != nil {
		_ = lf.file.Close()
	} else {
		log.Warn("cannot close uninstantiated file handle")
	}
}

// getReader returns the reader which is reading the file beginning from the end.
func (lf *logFile) getReader() io.Reader {
	reader := ReverseReader{file: lf.file}
	err := reader.seekEnd()
	if err != nil {
		log.Warnf("failed to seek to the files end: %v", err)
	}
	return &reader
}

func (lf *logFile) GetLogfilePathFromDoguName(doguName string) (string, error) {
	fileName := doguName
	if !strings.HasSuffix(doguName, ".log") {
		fileName = fmt.Sprintf("%s.%s", fileName, logFileExtension)
	}

	fullLogFilePath := path.Join(doguLogFilesPath, fileName)

	_, err := os.Stat(fullLogFilePath)
	if err != nil {
		return "", err
	}
	return fullLogFilePath, nil
}
