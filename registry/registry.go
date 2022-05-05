package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

const (
	// DirectoryGlobal contains the registry key to the global registry.
	DirectoryGlobal = "_global"
	// DirectoryHost contains the registry key to the host registry
	DirectoryHost = "_host"
	// KeyDoguPublicKey contains the key to a dogu's public key.
	KeyDoguPublicKey = "public.pem"
	// HostServiceCesappd is the service name of cesappd
	HostServiceCesappd = "cesappd"
	// HostServiceBackupWatcher is the service name of backup-watcher
	HostServiceBackupWatcher = "backup-watcher"
)

var (
	HostServices = map[string]struct{}{HostServiceCesappd: {}, HostServiceBackupWatcher: {}}
)

// Registry represents the main registry of a cloudogu ecosystem. The registry
// manage dogus, their configuration and their states.
type Registry interface {
	// GlobalConfig returns a ConfigurationContext for the global context
	GlobalConfig() ConfigurationContext
	// HostConfig returns a ConfigurationContext for the host context
	HostConfig(hostService string) ConfigurationContext
	// DoguConfig returns a ConfigurationContext for the given dogu
	DoguConfig(dogu string) ConfigurationContext
	// State returns the state object for the given dogu
	State(dogu string) State
	// DoguRegistry returns an object which is able to manage dogus
	DoguRegistry() DoguRegistry
	// BlueprintRegistry to maintain a blueprint history
	BlueprintRegistry() ConfigurationContext
}

// New creates a new registry
func New(configuration core.Registry) (Registry, error) {
	if configuration.Type != "etcd" {
		return nil, errors.Errorf("currently only etcd registry is supported, %s was provided", configuration.Type)
	}
	return newEtcdRegistry(configuration)
}
