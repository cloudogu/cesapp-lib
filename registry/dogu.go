package registry

import "github.com/cloudogu/cesapp-lib/core"

// DoguRegistry manages dogus on a ecosystem
type DoguRegistry interface {
	// Enable enables the given dogu
	Enable(dogu *core.Dogu) error
	// Register registeres the dogu on the registry
	Register(dogu *core.Dogu) error
	// Unregister unregisters the dogu on the registry
	Unregister(name string) error
	// Get returns the dogu which the given name
	Get(name string) (*core.Dogu, error)
	// GetAll returns all installed dogus
	GetAll() ([]*core.Dogu, error)
	// IsEnabled returns true if the dogu is installed
	IsEnabled(name string) (bool, error)
}
