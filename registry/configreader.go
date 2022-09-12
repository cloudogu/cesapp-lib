package registry

import (
	"strconv"

	"github.com/pkg/errors"
)

// NewConfigurationReader creates a new ConfigurationReader
func NewConfigurationReader(configuration ConfigurationContext) *ConfigurationReader {
	return &ConfigurationReader{configuration}
}

// ConfigurationReader is a simple abstraction for reading and converting registry Value
type ConfigurationReader struct {
	Configuration ConfigurationContext
}

// GetBool reads the configuration Value and converts it to a boolean. If the key could not be found, the function
// will return false.
func (configReader *ConfigurationReader) GetBool(key string) (bool, error) {
	stringValue, err := configReader.GetString(key)
	if err != nil {
		return false, err
	}

	if stringValue == "" {
		return false, nil
	}

	value, err := strconv.ParseBool(stringValue)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse %s as bool from registry key %s", stringValue, key)
	}

	return value, nil
}

// GetInt reads the configuration Value and converts it to an integer. If the key could not be found, the function
// will return 0.
func (configReader *ConfigurationReader) GetInt(key string) (int, error) {
	stringValue, err := configReader.GetString(key)
	if err != nil {
		return -1, err
	}

	if stringValue == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return -1, errors.Wrapf(err, "failed to parse %s as integer from registry key %s", stringValue, key)
	}

	return value, nil
}

// GetString reads a string from registry
func (configReader *ConfigurationReader) GetString(key string) (string, error) {
	exists, err := configReader.Configuration.Exists(key)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check if registry key %s exists", key)
	}

	if !exists {
		return "", nil
	}

	value, err := configReader.Configuration.Get(key)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %s from registry", key)
	}

	return value, nil
}
