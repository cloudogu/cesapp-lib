package registry

import "github.com/coreos/etcd/client"

// ConfigurationContext is able to manage the configuration of a single context
type ConfigurationContext interface {
	// Set sets a configuration value in current context
	Set(key, value string) error
	// SetWithLifetime sets a configuration value in current context with the given lifetime
	SetWithLifetime(key, value string, timeToLiveInSeconds int) error
	// Refresh resets the time to live of a key
	Refresh(key string, timeToLiveInSeconds int) error
	// Get returns a configuration value from the current context
	Get(key string) (string, error)
	// GetAll returns a map of key value pairs
	GetAll() (map[string]string, error)
	// Delete removes a configuration key and value from the current context
	Delete(key string) error
	// DeleteRecursive removes a configuration key or directory from the current context
	DeleteRecursive(key string) error
	// Exists returns true if configuration key exists in the current context
	Exists(key string) (bool, error)
	// RemoveAll remove all configuration keys
	RemoveAll() error
	// GetOrFalse return false and empty string when the configuration value does not exist.
	// Otherwise, return true and the configuration value, even when the configuration value is an empty string.
	GetOrFalse(key string) (bool, string, error)
}

// WatchConfigurationContext is just able to watch and query the configuration of a single context
type WatchConfigurationContext interface {
	// Get returns a configuration value from the current context
	Get(key string) (string, error)
	// Watch watches for changes of the provided key and sends the event through the channel
	Watch(key string, recursive bool, eventChannel chan *client.Response)
}
