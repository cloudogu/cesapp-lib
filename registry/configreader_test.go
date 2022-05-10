package registry_test

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	configuration := newConfigurationContext()
	configuration.Set("hello", "world")

	configReader := registry.NewConfigurationReader(configuration)
	value, err := configReader.GetString("hello")
	assert.Nil(t, err)
	assert.Equal(t, "world", value)
}

func TestGetNonExistingString(t *testing.T) {
	configuration := newConfigurationContext()

	configReader := registry.NewConfigurationReader(configuration)
	value, err := configReader.GetString("hello")
	assert.Nil(t, err)
	assert.Equal(t, "", value)
}

func TestGetBool(t *testing.T) {
	configuration := newConfigurationContext()
	configuration.Set("enabled", "true")
	configuration.Set("disabled", "false")

	configReader := registry.NewConfigurationReader(configuration)
	enabled, err := configReader.GetBool("enabled")
	assert.Nil(t, err)
	assert.True(t, enabled)

	disabled, err := configReader.GetBool("disabled")
	assert.Nil(t, err)
	assert.False(t, disabled)

	nonexisting, err := configReader.GetBool("nonexisting")
	assert.Nil(t, err)
	assert.False(t, nonexisting)
}

func TestGetBoolWithInvalidValue(t *testing.T) {
	configuration := newConfigurationContext()
	configuration.Set("enabled", "xyz")

	configReader := registry.NewConfigurationReader(configuration)
	_, err := configReader.GetBool("enabled")
	assert.NotNil(t, err)
}

func TestGetInt(t *testing.T) {
	configuration := newConfigurationContext()
	configuration.Set("fourtytwo", "42")

	configReader := registry.NewConfigurationReader(configuration)
	fourtytwo, err := configReader.GetInt("fourtytwo")
	assert.Nil(t, err)
	assert.Equal(t, 42, fourtytwo)

	nonexisting, err := configReader.GetInt("nonexisting")
	assert.Nil(t, err)
	assert.Equal(t, 0, nonexisting)
}

func TestGetIntWithInvalidValue(t *testing.T) {
	configuration := newConfigurationContext()
	configuration.Set("nan", "abc")

	configReader := registry.NewConfigurationReader(configuration)
	_, err := configReader.GetInt("nan")
	assert.NotNil(t, err)
}

func newConfigurationContext() registry.ConfigurationContext {
	registry := &registry.MockRegistry{}
	return registry.GlobalConfig()
}
