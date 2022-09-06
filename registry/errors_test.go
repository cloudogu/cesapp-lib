package registry

import (
	"errors"
	"strconv"
	"testing"

	"github.com/coreos/etcd/client"

	errors2 "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestIsKeyNotFoundError(t *testing.T) {
	t.Run("should return false if error is nil", func(t *testing.T) {
		actual := IsKeyNotFoundError(nil)
		require.False(t, actual)
	})
	t.Run("should return false if error is any other error", func(t *testing.T) {
		err := errors.New("oh noez")

		actual := IsKeyNotFoundError(err)

		require.False(t, actual)
	})
	t.Run("should return true if error is an original client error with code 100", func(t *testing.T) {
		originalErr := &client.Error{Code: client.ErrorCodeKeyNotFound}

		actual := IsKeyNotFoundError(originalErr)

		require.True(t, actual)
	})
	t.Run("should return true if original client error is wrapped several times by other errors", func(t *testing.T) {
		originalErr := &client.Error{Code: client.ErrorCodeKeyNotFound}
		wrappedErr := errors2.Wrap(originalErr, "oh noe")

		for i := 0; i < 3; i++ {
			wrappedErr = errors2.Wrap(wrappedErr, "oh noe:"+strconv.Itoa(i))
		}

		actual := IsKeyNotFoundError(wrappedErr)

		require.True(t, actual)
	})
}
