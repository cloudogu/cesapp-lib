package archive

import (
	"io"
	"os"
)

type ReverseReader struct {
	file *os.File
}

// seekEnd seeks to the final byte of the file
func (r *ReverseReader) seekEnd() error {
	_, err := r.file.Seek(0, io.SeekEnd)
	return err
}

// Read the file backwards
func (r *ReverseReader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	// This no-op gives us the current offset value
	offset, err := r.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	var m int
	for i := 0; i < len(b); i++ {
		if offset == 0 {
			return m, io.EOF
		}
		// Seek in case someone else is relying on seek too
		offset, err = r.file.Seek(-1, io.SeekCurrent)
		if err != nil {
			return m, err // Should never happen
		}

		// Just read one byte at a time
		n, err := r.file.ReadAt(b[i:i+1], offset)
		if err != nil {
			return m + n, err // Should never happen
		}
		m += n
	}
	return m, nil
}
