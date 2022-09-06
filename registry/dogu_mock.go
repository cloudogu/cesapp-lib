package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client/v2"
)

func createMockDoguRegistry() *mockDoguRegistry {
	return &mockDoguRegistry{
		enabled: []string{},
		dogus:   make(map[string]*core.Dogu),
		// faultyDogus:   make(map[string]*core.Dogu),
	}
}

type mockDoguRegistry struct {
	enabled []string
	dogus   map[string]*core.Dogu
	// faultyDogus   map[string]*core.Dogu
}

func (reg *mockDoguRegistry) Enable(dogu *core.Dogu) error {
	reg.enabled = append(reg.enabled, dogu.GetSimpleName())
	return nil
}

func (reg *mockDoguRegistry) Register(dogu *core.Dogu) error {
	reg.dogus[dogu.GetSimpleName()] = dogu
	return nil
}

func (reg *mockDoguRegistry) RegisterToFail(dogu *core.Dogu) error {
	reg.dogus[dogu.GetSimpleName()] = dogu
	return nil
}

func (reg *mockDoguRegistry) Unregister(name string) error {
	enabled := []string{}
	for _, dogu := range reg.enabled {
		if dogu != name {
			enabled = append(enabled, dogu)
		}
	}
	reg.enabled = enabled
	delete(reg.dogus, name)
	return nil
}

func (reg *mockDoguRegistry) Get(name string) (*core.Dogu, error) {
	if dogu, ok := reg.dogus[name]; ok {
		if dogu.Category == "nil" {
			return nil, nil
		}
		return dogu, nil
	}
	return nil, errors.Errorf("could not find dogu %s", name)
}

func (reg *mockDoguRegistry) GetAll() ([]*core.Dogu, error) {
	dogus := []*core.Dogu{}
	for _, dogu := range reg.dogus {
		dogus = append(dogus, dogu)
	}
	return dogus, nil
}

func (reg *mockDoguRegistry) IsEnabled(name string) (bool, error) {
	for _, dogu := range reg.enabled {
		if dogu == name {
			return true, nil
		}
	}
	return false, nil
}

func (reg *mockDoguRegistry) Watch(_ string, _ bool, _ chan *client.Response) {

}
