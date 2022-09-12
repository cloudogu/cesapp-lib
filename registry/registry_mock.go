package registry

// MockRegistry is a in memory implementation of the registry interface. This registry is only for testing purposes
// and should never be used in production environments.
//
// Deprecated: MockRegistry exists for historical compatibility
// and should not be used. Use mocks.Registry instead.
type MockRegistry struct {
	global       *mockConfigurationContext
	hosts        map[string]*mockConfigurationContext
	dogus        map[string]*mockConfigurationContext
	doguRegistry *mockDoguRegistry
}

// GlobalConfig returns a mock implementation of the global configuration context
func (mr *MockRegistry) GlobalConfig() ConfigurationContext {
	if mr.global == nil {
		mr.global = createMockConfigurationContext()
	}
	return mr.global
}

func (mr *MockRegistry) HostConfig(hostService string) ConfigurationContext {
	if mr.hosts == nil {
		mr.hosts = make(map[string]*mockConfigurationContext)
	}
	hostCtx := mr.dogus[hostService]
	if hostCtx == nil {
		hostCtx = createMockConfigurationContext()
		mr.hosts[hostService] = hostCtx
	}
	return hostCtx
}

// DoguConfig returns a mock implementation of a dogu configuration context
func (mr *MockRegistry) DoguConfig(dogu string) ConfigurationContext {
	if mr.dogus == nil {
		mr.dogus = make(map[string]*mockConfigurationContext)
	}
	doguCtx := mr.dogus[dogu]
	if doguCtx == nil {
		doguCtx = createMockConfigurationContext()
		mr.dogus[dogu] = doguCtx
	}
	return doguCtx
}

// State is currently not implemented and returns always nil
func (mr *MockRegistry) State(dogu string) State {
	return nil
}

// DoguRegistry returns a mock implementation of the dogu registry
func (mr *MockRegistry) DoguRegistry() DoguRegistry {
	if mr.doguRegistry == nil {
		mr.doguRegistry = createMockDoguRegistry()
	}
	return mr.doguRegistry
}

// GlobalConfig returns a mock implementation of the global configuration context
func (mr *MockRegistry) BlueprintRegistry() ConfigurationContext {
	if mr.global == nil {
		mr.global = createMockConfigurationContext()
	}
	return mr.global
}

// RootConfig was added to deprecated mock for legacy code support and has no functionality
func (mr *MockRegistry) RootConfig() WatchConfigurationContext {
	return nil
}

// GetNode was added to deprecated mock for legacy code support and has no functionality
func (mr *MockRegistry) GetNode() (Node, error) {
	return Node{}, nil
}
