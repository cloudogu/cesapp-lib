package remote

import (
	"github.com/cloudogu/cesapp-lib/core"
)

// Registry is able to manage the remote dogu registry.
type Registry interface {
	// Create the dogu on the remote server.
	Create(dogu *core.Dogu) error
	// Get returns the detail about a dogu from the remote server.
	Get(name string) (*core.Dogu, error)
	// GetVersion returns a version specific detail about the dogu. Name is mandatory. Version is optional; if no version
	// is given then the newest version will be returned.
	GetVersion(name, version string) (*core.Dogu, error)
	// GetAll returns all dogus from the remote server.
	GetAll() ([]*core.Dogu, error)
	// GetVersionsOf return all versions of a dogu.
	GetVersionsOf(name string) ([]core.Version, error)
	// SetUseCache disables or enables the caching for the remote registry.
	SetUseCache(useCache bool)
	// Delete removes a specific dogu descriptor from the dogu registry.
	Delete(dogu *core.Dogu) error
}

// New creates a new remote registry
func New(remoteConfig *core.Remote, credentials *core.Credentials) (Registry, error) {
	return newHTTPRemote(remoteConfig, credentials)
}
