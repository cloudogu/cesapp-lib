package mocks

import (
	"github.com/stretchr/testify/mock"
)

const (
	AnyLifetime = -1
	Anything    = mock.Anything
)

type RegistryMocks struct {
	Registry      *Registry
	SubRegistries map[string]*ConfigurationContext
	DoguRegistry  *DoguRegistry
	WatchContext  *WatchConfigurationContext
	Node          *RegistryNode
}

// CreateMockRegistry creates a mock registry containing default values for any sub-registry.
// The doguRegs parameter contains a list with all dogu registries that should be exist. Pass nil if you don't need any dogu registry.
// Note: AssertExpectations on the registry may fail because of the predefined expectations.
func CreateMockRegistry(doguRegs []string) RegistryMocks {
	registry := &Registry{}
	globalConfig := &ConfigurationContext{}
	blueprintConfig := &ConfigurationContext{}
	doguReg := &DoguRegistry{}
	rootConfig := &WatchConfigurationContext{}
	registryNode := &RegistryNode{}
	registry.On("GlobalConfig").Return(globalConfig)
	registry.On("DoguRegistry").Return(doguReg)
	registry.On("BlueprintRegistry").Return(blueprintConfig)
	registry.On("RootConfig").Return(rootConfig)
	registry.On("GetNode").Return(registryNode)

	registries := map[string]*ConfigurationContext{}

	for _, doguReg := range doguRegs {
		registries[doguReg] = addDoguRegistry(registry, doguReg)
	}

	registries["_global"] = globalConfig
	registries["blueprints"] = blueprintConfig

	return RegistryMocks{
		Registry:      registry,
		SubRegistries: registries,
		DoguRegistry:  doguReg,
		WatchContext:  rootConfig,
		Node:          registryNode,
	}
}

// OnGet provides a helper function to mock the "Get" method of a configuration context
func OnGet(config *ConfigurationContext, key string, returnValue string, returnError error) {
	config.On("Get", key).Return(returnValue, returnError)
}

// OnSet provides a helper function to mock the "Set" method of a configuration context
func OnSet(config *ConfigurationContext, keyToSet string, valueToSet string, returnError error) {
	config.On("Set", keyToSet, valueToSet).Return(returnError)
}

// OnDelete provides a helper function to mock the "Delete" method of a configuration context
func OnDelete(config *ConfigurationContext, key string, returnError error) {
	config.On("Delete", key).Return(returnError)
}

// OnDeleteRecursive provides a helper function to mock the "Delete" method of a configuration context
func OnDeleteRecursive(config *ConfigurationContext, key string, returnError error) {
	config.On("DeleteRecursive", key).Return(returnError)
}

// OnExists provides a helper function to mock the "Exists" method of a configuration context
func OnExists(config *ConfigurationContext, key string, returnExists bool, returnError error) {
	config.On("Exists", key).Return(returnExists, returnError)
}

// OnGetOrFalse provides a helper function to mock the "GetOrFalse" method of a configuration context
func OnGetOrFalse(config *ConfigurationContext, key string, returnValue string, returnExists bool, returnError error) {
	config.On("GetOrFalse", key).Return(returnExists, returnValue, returnError)
}

// OnRefresh provides a helper function to mock the "Refresh" method of a configuration context
func OnRefresh(config *ConfigurationContext, key string, ttl int, returnError error) {
	if ttl == AnyLifetime {
		config.On("Refresh", key, mock.AnythingOfType("int")).Return(returnError)
	} else {
		config.On("Refresh", key, ttl).Return(returnError)
	}
}

// OnSetWithLifetime provides a helper function to mock the "SetWithLifetime" method of a configuration context
func OnSetWithLifetime(config *ConfigurationContext, key string, value string, ttl int, returnError error) {
	if ttl == AnyLifetime {
		config.On("SetWithLifetime", key, value, mock.AnythingOfType("int")).Return(returnError)
	} else {
		config.On("SetWithLifetime", key, value, ttl).Return(returnError)
	}
}

// OnRemoveAll provides a helper function to mock the "RemoveAll" method of a configuration context
func OnRemoveAll(config *ConfigurationContext, returnError error) {
	config.On("RemoveAll").Return(returnError)
}

func addDoguRegistry(registry *Registry, regName string) *ConfigurationContext {
	doguReg := &ConfigurationContext{}

	registry.On("DoguConfig", mock.MatchedBy(func(name string) bool {
		return name == regName
	})).Return(doguReg)

	return doguReg
}
