package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegistryKeyConfigValidator_CheckRegistryKey(t *testing.T) {
	t.Run("should successfully execute registry validation", func(t *testing.T) {
		// given
		value := "pkcs1v15"
		path := "key_provider"
		globalReg := new(mockGlobalRegistry)
		globalReg.On("Get", path).Return("pkcs1v15", nil)
		sut := &RegistryKeyConfigValidator{}

		// when
		err := sut.CheckRegistryKey(value, path, globalReg, false)

		// then
		require.NoError(t, err)
	})
	t.Run("should execute registry validation", func(t *testing.T) {
		// given
		path := "key_provider"
		globalReg := new(mockGlobalRegistry)
		globalReg.On("Get", path).Return("", nil)
		sut := &RegistryKeyConfigValidator{}

		// when
		err := sut.CheckRegistryKey("fileConfig", path, globalReg, true)

		// then
		require.NoError(t, err)
	})
	t.Run("should execute registry validation allow empty", func(t *testing.T) {
		// given
		value := "pkcs1v15"
		path := "key_provider"
		globalReg := new(mockGlobalRegistry)
		globalReg.On("Get", path).Return("", errors.New("100: Key not found"))
		sut := CreateRegistryKeyConfigValidator()

		// when
		err := sut.CheckRegistryKey(value, path, globalReg, true)

		// then
		require.NoError(t, err)
	})
	t.Run("should return any registry validation error", func(t *testing.T) {
		// given
		value := "pkcs1v15"
		path := "key_provider"
		globalReg := new(mockGlobalRegistry)
		globalReg.On("Get", path).Return("oaesp", nil)
		sut := CreateRegistryKeyConfigValidator()

		// when
		err := sut.CheckRegistryKey(value, path, globalReg, false)

		// then
		require.Error(t, err)
	})
}

func Test_validateRegistryKey(t *testing.T) {
	t.Run("should return without error", func(t *testing.T) {
		// given
		globalReg := new(mockGlobalRegistry)
		path := "key_provider"
		globalReg.On("Get", path).Return("pkcs1v15", nil)

		// when
		err := validateRegistryKey(path, "pkcs1v15", globalReg, false)

		// then
		require.NoError(t, err)
		globalReg.AssertExpectations(t)
	})
	t.Run("should return error on difference", func(t *testing.T) {
		// given
		globalReg := new(mockGlobalRegistry)
		path := "key_provider"
		globalReg.On("Get", path).Return("", assert.AnError)

		// when
		err := validateRegistryKey(path, "oaesp", globalReg, false)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get global config value for 'key_provider'")
		globalReg.AssertExpectations(t)
	})
}

type mockGlobalRegistry struct {
	mock.Mock
}

func (m *mockGlobalRegistry) Get(key string) (value string, err error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}
