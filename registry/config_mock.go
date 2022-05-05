package registry

import (
	"github.com/coreos/etcd/client"
	"github.com/pkg/errors"
)

func createMockConfigurationContext() *mockConfigurationContext {
	return &mockConfigurationContext{
		values:    make(map[string]string),
		lifetimes: make(map[string]int),
	}
}

// Deprecated: mockConfigurationContext exists for historical compatibility
// and should not be used. Use mocks.Registry instead.
type mockConfigurationContext struct {
	values    map[string]string
	lifetimes map[string]int
}

func (mcc *mockConfigurationContext) Set(key, value string) error {
	mcc.values[key] = value
	return nil
}

func (mcc *mockConfigurationContext) SetWithLifetime(key, value string, timeToLiveInSeconds int) error {
	_ = mcc.Set(key, value)
	mcc.lifetimes[key] = timeToLiveInSeconds
	return nil
}

func (mcc *mockConfigurationContext) Get(key string) (string, error) {
	if _, exists := mcc.values[key]; exists {
		return mcc.values[key], nil
	}
	return "", client.Error{Code: client.ErrorCodeKeyNotFound}
}

func (mcc *mockConfigurationContext) Refresh(key string, timeToLiveInSeconds int) error {
	// Nothing to be done here because mock etcd does not support timeout anyway
	return nil
}

func (mcc *mockConfigurationContext) GetAll() (map[string]string, error) {
	return mcc.values, nil
}

func (mcc *mockConfigurationContext) Delete(key string) error {
	if _, exists := mcc.values[key]; exists {
		delete(mcc.values, key)
		return nil
	}
	return errors.Errorf("key %s does not exist", key)
}

func (mcc *mockConfigurationContext) DeleteRecursive(key string) error {
	return nil
}

func (mcc *mockConfigurationContext) Exists(key string) (bool, error) {
	_, exists := mcc.values[key]
	return exists, nil
}

func (mcc *mockConfigurationContext) RemoveAll() error {
	mcc.values = make(map[string]string)
	return nil
}

func (mcc *mockConfigurationContext) GetOrFalse(key string) (exists bool, value string, err error) {
	if _, exists := mcc.values[key]; exists {
		return true, mcc.values[key], nil
	}
	return false, "", nil
}
