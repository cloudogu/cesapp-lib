package log

import (
	"fmt"
	"github.com/prometheus/common/log"
	"io"
	"os"
	"path"
	"path/filepath"
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

func openDoguLogFile(doguName string) (*logFile, error) {
	return openDoguLogFileFor(doguName, doguLogFilesPath)
}

func openDoguLogFileFor(doguName string, pathPrefix string) (*logFile, error) {
	fileName := fmt.Sprintf("%s.%s", doguName, logFileExtension)
	fullLogFilePath := path.Join(pathPrefix, fileName)

	cleanedFullLogFilePath, err := cleanFilePath(fullLogFilePath, pathPrefix)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(cleanedFullLogFilePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(cleanedFullLogFilePath, os.O_RDONLY, os.ModePerm)
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

func cleanFilePath(path, allowedDirectoryPrefix string) (string, error) {
	cleanedFilepath := filepath.Clean(path)
	if strings.HasPrefix(cleanedFilepath, allowedDirectoryPrefix) {
		return cleanedFilepath, nil
	}
	return "", fmt.Errorf("cannot access dogu log file outside its log directory: %s", cleanedFilepath)
}
