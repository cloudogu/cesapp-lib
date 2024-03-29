package registry

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"go.etcd.io/etcd/client/v2"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type etcdClient interface {
	Exists(key string) (bool, error)
	Get(key string) (string, error)
	GetRecursive(key string) (map[string]string, error)
	GetChildrenPaths(key string) ([]string, error)
	Set(key string, value string, options *client.SetOptions) (string, error)
	Delete(key string, options *client.DeleteOptions) error
	DeleteRecursive(key string) error
	Watch(ctx context.Context, key string, recursive bool, eventChannel chan *client.Response)
}

type etcdConfigurationContext struct {
	parent string
	client etcdClient
}

type etcdWatchConfigurationContext struct {
	client etcdClient
}

// Set sets a configuration value in current context
func (ecc *etcdConfigurationContext) Set(key, value string) error {
	return ecc.set(key, value, nil)
}

// SetWithLifetime sets a configuration value in current context with the given lifetime
func (ecc *etcdConfigurationContext) SetWithLifetime(key, value string, timeToLiveInSeconds int) error {
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", timeToLiveInSeconds))
	if err != nil {
		return errors.Wrapf(err, "could not create the curation '%d'", timeToLiveInSeconds)
	}
	return ecc.set(key, value, &client.SetOptions{
		TTL: duration,
	})
}

// Refresh will refresh the ttl of a key
func (ecc *etcdConfigurationContext) Refresh(key string, timeToLiveInSeconds int) error {
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", timeToLiveInSeconds))
	if err != nil {
		return errors.Wrapf(err, "could not create the curation '%d'", timeToLiveInSeconds)
	}

	err = ecc.set(key, "", &client.SetOptions{
		TTL:     duration,
		Refresh: true,
	})
	return err
}

func (ecc *etcdConfigurationContext) set(key, value string, options *client.SetOptions) error {
	path := ecc.parent + "/" + key
	core.GetLogger().Debug("try to set config key", path)

	_, err := ecc.client.Set(path, value, options)
	if err != nil {
		return errors.Wrapf(err, "could not set value %s", path)
	}

	return err
}

func (ecc *etcdConfigurationContext) Get(key string) (string, error) {
	return Get(ecc.parent, key, ecc.client)
}

// GetAll returns a map of key value pairs
func (ecc *etcdConfigurationContext) GetAll() (map[string]string, error) {
	core.GetLogger().Debugf("try to get all configuration keys and values from %s", ecc.parent)

	keyValuePairs, err := ecc.client.GetRecursive(ecc.parent)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get key value pairs recursive from %s", ecc.parent)
	}

	return keyValuePairs, nil
}

// Delete removes a configuration key and value from the current context
func (ecc *etcdConfigurationContext) Delete(key string) error {
	path := ecc.parent + "/" + key
	core.GetLogger().Debug("try to delete config key", path)

	err := ecc.client.Delete(path, nil)
	if err != nil {
		return errors.Wrapf(err, "could not delete value at %s", path)
	}
	return nil
}

// DeleteRecursive deletes a configuration key from the current context recursively.
func (ecc *etcdConfigurationContext) DeleteRecursive(key string) error {
	path := ecc.parent + "/" + key
	core.GetLogger().Debugf("try to delete config key '%s' recursive", path)

	err := ecc.client.DeleteRecursive(path)
	if err != nil {
		return errors.Wrapf(err, "could not delete value at %s", path)
	}
	return nil
}

// Exists returns true if configuration key exists in the current context
func (ecc *etcdConfigurationContext) Exists(key string) (bool, error) {
	path := ecc.parent + "/" + key
	core.GetLogger().Debugf("try to check if config key %s exists", path)

	exists, err := ecc.client.Exists(path)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check if key %s exists", path)
	}

	return exists, nil
}

// RemoveAll removes all configuration key
func (ecc *etcdConfigurationContext) RemoveAll() error {
	err := ecc.client.DeleteRecursive(ecc.parent)
	if err != nil {
		return errors.Wrapf(err, "could not remove all configuration keys from %s", ecc.parent)
	}
	return nil
}

// GetOrFalse return false and empty string when the configuration value does not exist.
// Otherwise return true and the configuration value, even when the configuration value is an empty string.
func (ecc *etcdConfigurationContext) GetOrFalse(key string) (bool, string, error) {
	exists, err := ecc.Exists(key)
	if err != nil {
		return false, "", err
	}

	if !exists {
		return false, "", nil
	}

	value, err := ecc.Get(key)
	if err != nil {
		return false, "", err
	}

	return true, value, nil
}

// Watch watches for changes of the provided key and sends the event through the channel
func (ewcc *etcdWatchConfigurationContext) Watch(ctx context.Context, key string, recursive bool, eventChannel chan *client.Response) {
	core.GetLogger().Debugf("starting watcher on key %s", key)
	ewcc.client.Watch(ctx, key, recursive, eventChannel)
}

func (ewcc *etcdWatchConfigurationContext) Get(key string) (string, error) {
	return Get("", key, ewcc.client)
}

// Get returns a configuration value from the current context, otherwise it returns an error. If the given key cannot be
// found a KeyNotFoundError is returned.
func Get(parent string, key string, client etcdClient) (string, error) {
	path := parent + "/" + key
	if parent == "" {
		if strings.HasPrefix(key, "/") {
			path = key
		}
	}

	core.GetLogger().Debug("try to get config key", path)

	value, err := client.Get(path)
	if err != nil {
		return "", errors.Wrapf(err, "could not get value %s", path)
	}
	return value, nil
}

// GetChildrenPaths returns an array of all children keys of the given key
func (ewcc *etcdWatchConfigurationContext) GetChildrenPaths(key string) ([]string, error) {
	return ewcc.client.GetChildrenPaths(key)
}
