package log

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

const (
	supportArchiveFileName = "support-archive.zip"
)

// TODO: remove me; Anleitung fileszippen: https://golang.cafe/blog/golang-zip-file-example.html
func WriteSupportArchive(logfilePaths []string) error {
	zippedArchiveFile, err := os.Create(supportArchiveFileName)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zippedArchiveFile)

	for i, file := range logfilePaths {
		fmt.Printf("opening file: %v, %s\n", i, file)
		doguLogFile, err := SelectLogFile(file)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = doguLogFile.file.Close()
		}()

		createdFileInZip, err := zipWriter.Create(file)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(createdFileInZip, doguLogFile.file); err != nil {
			panic(err)
		}

	}
	defer func() {
		err = zippedArchiveFile.Close()
	}()

	return err
}
