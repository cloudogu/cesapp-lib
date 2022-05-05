package remote

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

// NewMockRegistry return a mock implementation of the remote registry and should be used for testing purposes only
func NewMockRegistry() Registry {
	return &mockRegistry{
		values: make(map[string][]*core.Dogu),
	}
}

type mockRegistry struct {
	values map[string][]*core.Dogu
}

// Create the dogu on the remote server
func (reg *mockRegistry) Create(dogu *core.Dogu) error {
	dogus := reg.values[dogu.Name]
	dogus = append(dogus, dogu)
	reg.values[dogu.Name] = dogus
	return nil
}

// Get returns the detail about a dogu from the remote server
func (reg *mockRegistry) Get(name string) (*core.Dogu, error) {
	dogus := reg.values[name]
	if len(dogus) >= 1 {
		return dogus[len(dogus)-1], nil
	}
	return nil, errors.Errorf("dogu %s not found", name)
}

// GetVersion returns a version specific detail about the dogu
func (reg *mockRegistry) GetVersion(name, version string) (*core.Dogu, error) {
	dogus := reg.values[name]
	for _, dogu := range dogus {
		if dogu.Name == name && dogu.Version == version {
			return dogu, nil
		}
	}
	return nil, errors.Errorf("could not find dogu %s in version %s", name, version)
}

// GetAll returns all dogus from the remote server
func (reg *mockRegistry) GetAll() ([]*core.Dogu, error) {
	dogus := []*core.Dogu{}
	for _, doguVersions := range reg.values {
		if len(doguVersions) > 0 {
			dogus = append(dogus, doguVersions[len(doguVersions)-1])
		}
	}
	return dogus, nil
}

// GetVersionOf returns all versions from the remote server
func (reg *mockRegistry) GetVersionsOf(name string) ([]core.Version, error) {
	versions := make([]core.Version, 0)
	dogus := reg.values[name]
	for _, dogu := range dogus {
		version, err := dogu.GetVersion()
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// SetUseCache disables or enables the caching for the remote mock registry
func (reg *mockRegistry) SetUseCache(useCache bool) {
	// mockRegistry doesn't use a cache
}

// Delete removes a specific dogu descriptor from the dogu registry.
func (reg *mockRegistry) Delete(dogu *core.Dogu) error {
	delete(reg.values, dogu.Name)

	return nil
}
