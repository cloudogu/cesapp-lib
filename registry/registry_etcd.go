package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"go.etcd.io/etcd/client/v2"
	"os"

	"github.com/pkg/errors"
)

type etcdRegistry struct {
	client *resilentEtcdClient
}

func newEtcdRegistry(configuration core.Registry) (*etcdRegistry, error) {
	etcdClient, err := createEtcdClient(configuration)
	if err != nil {
		return nil, err
	}
	return &etcdRegistry{etcdClient}, nil
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

	return newResilientEtcdClient(endpoints, configuration.RetryPolicy)
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

// RootConfig returns a ConfigurationContext for the root context
func (er *etcdRegistry) RootConfig() WatchConfigurationContext {
	return &etcdWatchConfigurationContext{er.client}
}

// GetNode returns a ConfigurationContext for the root context
func (er *etcdRegistry) GetNode() (Node, error) {
	mainNode, err := er.client.getMainNode()
	if err != nil {
		return Node{}, err
	}
	return mapEtcdNodeToRegistryNode(mainNode), nil
}

func mapEtcdNodeToRegistryNode(node *client.Node) Node {
	result := Node{
		SubNodes: []Node{},
		IsDir:    node.Dir,
		FullKey:  node.Key,
		Value:    node.Value,
	}

	for _, child := range node.Nodes {
		result.SubNodes = append(result.SubNodes, mapEtcdNodeToRegistryNode(child))
	}

	return result
}
