package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("hello").Return(true, nil)
	configuration.EXPECT().Get("hello").Return("world", nil)

	configReader := NewConfigurationReader(configuration)
	value, err := configReader.GetString("hello")
	assert.Nil(t, err)
	assert.Equal(t, "world", value)
}

func TestGetNonExistingString(t *testing.T) {
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("hello").Return(true, nil)
	configuration.EXPECT().Get("hello").Return("", nil)

	configReader := NewConfigurationReader(configuration)
	value, err := configReader.GetString("hello")
	assert.Nil(t, err)
	assert.Equal(t, "", value)
}

func TestGetBool(t *testing.T) {
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("enabled").Return(true, nil)
	configuration.EXPECT().Get("enabled").Return("true", nil)
	configuration.EXPECT().Exists("disabled").Return(true, nil)
	configuration.EXPECT().Get("disabled").Return("false", nil)
	configuration.EXPECT().Exists("nonexisting").Return(true, nil)
	configuration.EXPECT().Get("nonexisting").Return("false", nil)

	configReader := NewConfigurationReader(configuration)
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
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("enabled").Return(true, nil)
	configuration.EXPECT().Get("enabled").Return("", assert.AnError)

	configReader := NewConfigurationReader(configuration)
	_, err := configReader.GetBool("enabled")
	assert.NotNil(t, err)
}

func TestGetInt(t *testing.T) {
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("fourtytwo").Return(true, nil).Once()
	configuration.EXPECT().Get("fourtytwo").Return("42", nil).Once()
	configuration.EXPECT().Exists("nonexisting").Return(false, nil).Once()

	configReader := NewConfigurationReader(configuration)
	fourtytwo, err := configReader.GetInt("fourtytwo")
	assert.Nil(t, err)
	assert.Equal(t, 42, fourtytwo)

	nonexisting, err := configReader.GetInt("nonexisting")
	assert.Nil(t, err)
	assert.Equal(t, 0, nonexisting)
}

func TestGetIntWithInvalidValue(t *testing.T) {
	configuration := NewMockConfigurationContext(t)
	configuration.EXPECT().Exists("nan").Return(true, nil).Once()
	configuration.EXPECT().Get("nan").Return("abc", nil).Once()

	configReader := NewConfigurationReader(configuration)
	_, err := configReader.GetInt("nan")
	assert.NotNil(t, err)
}
