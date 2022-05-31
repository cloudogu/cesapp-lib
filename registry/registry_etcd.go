package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"os"

	"github.com/pkg/errors"
)

type etcdRegistry struct {
	client *resilentEtcdClient
}

func newEtcdRegistry(configuration core.Registry) (*etcdRegistry, error) {
	client, err := createEtcdClient(configuration)
	if err != nil {
		return nil, err
	}
	return &etcdRegistry{client}, nil
}

func createEtcdClient(configuration core.Registry) (*resilentEtcdClient, error) {
	if configuration.Type != "etcd" {
		return nil, errors.New("currently only etcd registries are supported")
	}

	var endpoints []string
	endpoint := os.Getenv("REGISTRY_ENDPOINT")
	if len(endpoint) > 0 {
		endpoints = append(endpoints, endpoint)
	} else {
		endpoints = configuration.Endpoints
	}

	return newResilentEtcdClient(endpoints)
}

// GlobalConfig returns a ConfigurationContext for the global context.
func (er *etcdRegistry) GlobalConfig() ConfigurationContext {
	return &etcdConfigurationContext{"/config/" + DirectoryGlobal, er.client}
}

// HostConfig returns a ConfigurationContext for the host context.
func (er *etcdRegistry) HostConfig(hostService string) ConfigurationContext {
	return &etcdConfigurationContext{"/config/" + directoryHost + "/" + hostService, er.client}
}

// DoguConfig returns a ConfigurationContext for the given dogu.
func (er *etcdRegistry) DoguConfig(dogu string) ConfigurationContext {
	return &etcdConfigurationContext{"/config/" + dogu, er.client}
}

// State returns the state object for the given dogu.
func (er *etcdRegistry) State(dogu string) State {
	return &etcdState{"/state/" + dogu, er.client}
}

// DoguRegistry returns an object which is able to manage dogus.
func (er *etcdRegistry) DoguRegistry() DoguRegistry {
	return newCombinedEtcdDoguRegistry(er.client, "/dogu", "/dogu_v2")
}

// BlueprintRegistry returns a registry ConfigurationContext for manipulating blueprint registry entries.
func (er *etcdRegistry) BlueprintRegistry() ConfigurationContext {
	return &etcdConfigurationContext{"/blueprint", er.client}
}

// RootRegistry returns a ConfigurationContext for the root context
func (er *etcdRegistry) RootRegistry() WatchConfigurationContext {
	return &etcdWatchConfigurationContext{"/", er.client}
}
