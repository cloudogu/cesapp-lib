package registry

import (
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client/v2"
)

// IsKeyNotFoundError returns true if the given error is a or contains a registry keyNotFoundError, otherwise false.
// It returns false if the given error is nil.
func IsKeyNotFoundError(err error) bool {
	foundNaturally := isKeyNotFound(err)

	cause := errors.Cause(err)
	foundAsCause := isKeyNotFound(cause)

	return foundNaturally || foundAsCause
}

func isKeyNotFound(err error) bool {
	if cErr, ok := err.(client.Error); ok {
		return cErr.Code == client.ErrorCodeKeyNotFound
	}
	return false
}
