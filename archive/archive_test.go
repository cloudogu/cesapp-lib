package archive

import (
	"errors"
	"github.com/cloudogu/cesapp-lib/archive/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClose(t *testing.T) {
	t.Run("fail on close", func(t *testing.T) {
		mockWriter := &mocks.ZipWriter{}
		mockWriter.On("Close").Return(errors.New("testerror"))

		handler := &Handler{
			writer: mockWriter,
		}

		err := handler.Close()
		assert.Error(t, err)
	})
}
