package credentials

import (
	"github.com/cloudogu/cesapp-lib/core"
)

// DefaultStore is the id of the store which contains the credentials for backend communication
const DefaultStore = "_default"

// Store manage and stores credentials
type Store interface {
	// Add adds new credentials to the store
	Add(id string, creds *core.Credentials) error
	// Remove removes credentials from the store
	Remove(id string) error
	// Get returns the credentials with the given id
	Get(id string) *core.Credentials
}

// NewStore creates a new credential store
func NewStore(directory string) (Store, error) {
	return newSimpleStore(directory)
}
