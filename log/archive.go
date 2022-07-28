package log

import (
	"archive/zip"
	"os"
)

const (
	supportArchiveFileName = "support-archive.zip"
)

// TODO: remove me; Anleitung fileszippen: https://golang.cafe/blog/golang-zip-file-example.html
func WriteSupportArchive(logfiles []string) {
	archive, err := os.Create(supportArchiveFileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
}
