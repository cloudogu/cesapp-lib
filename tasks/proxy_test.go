package tasks_test

import (
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/cloudogu/cesapp-lib/tasks"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProxySettings(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(true, nil)
		globalConfig.On("Get", "proxy/enabled").Return("true", nil)
		globalConfig.On("Exists", "proxy/server").Return(true, nil)
		globalConfig.On("Get", "proxy/server").Return("proxy.cloudogu.com", nil)
		globalConfig.On("Exists", "proxy/port").Return(true, nil)
		globalConfig.On("Get", "proxy/port").Return("3864", nil)
		globalConfig.On("Exists", "proxy/username").Return(true, nil)
		globalConfig.On("Get", "proxy/username").Return("bob", nil)
		globalConfig.On("Exists", "proxy/password").Return(true, nil)
		globalConfig.On("Get", "proxy/password").Return("bob123", nil)

		proxySettings, err := tasks.CreateProxySettings(reg)
		assert.Nil(t, err)
		assert.True(t, proxySettings.Enabled)
		assert.Equal(t, "proxy.cloudogu.com", proxySettings.Server)
		assert.Equal(t, 3864, proxySettings.Port)
		assert.Equal(t, "bob", proxySettings.Username)
		assert.Equal(t, "bob123", proxySettings.Password)
	})

	t.Run("should return error on missing proxy server", func(t *testing.T) {
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(true, nil)
		globalConfig.On("Get", "proxy/enabled").Return("true", nil)
		globalConfig.On("Exists", "proxy/server").Return(true, nil)
		globalConfig.On("Get", "proxy/server").Return("proxy.cloudogu.com", nil)
		globalConfig.On("Exists", "proxy/port").Return(true, nil)
		globalConfig.On("Get", "proxy/port").Return("3864", assert.AnError)

		_, err := tasks.CreateProxySettings(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("should return error on missing username", func(t *testing.T) {
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(true, nil)
		globalConfig.On("Get", "proxy/enabled").Return("true", nil)
		globalConfig.On("Exists", "proxy/server").Return(true, nil)
		globalConfig.On("Get", "proxy/server").Return("proxy.cloudogu.com", nil)
		globalConfig.On("Exists", "proxy/port").Return(true, nil)
		globalConfig.On("Get", "proxy/port").Return("3864", nil)
		globalConfig.On("Exists", "proxy/username").Return(true, nil)
		globalConfig.On("Get", "proxy/username").Return("bob", assert.AnError)

		_, err := tasks.CreateProxySettings(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("should return error on missing password", func(t *testing.T) {
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(true, nil)
		globalConfig.On("Get", "proxy/enabled").Return("true", nil)
		globalConfig.On("Exists", "proxy/server").Return(true, nil)
		globalConfig.On("Get", "proxy/server").Return("proxy.cloudogu.com", nil)
		globalConfig.On("Exists", "proxy/port").Return(true, nil)
		globalConfig.On("Get", "proxy/port").Return("3864", nil)
		globalConfig.On("Exists", "proxy/username").Return(true, nil)
		globalConfig.On("Get", "proxy/username").Return("bob", nil)
		globalConfig.On("Exists", "proxy/password").Return(true, nil)
		globalConfig.On("Get", "proxy/password").Return("bob123", assert.AnError)

		_, err := tasks.CreateProxySettings(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("should return error on missing port", func(t *testing.T) {
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(true, nil)
		globalConfig.On("Get", "proxy/enabled").Return("true", nil)
		globalConfig.On("Exists", "proxy/server").Return(true, nil)
		globalConfig.On("Get", "proxy/server").Return("", assert.AnError)

		_, err := tasks.CreateProxySettings(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("return disabled configuration if proxy is not enabled in config", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(false, nil)

		// when
		proxySettings, err := tasks.CreateProxySettings(reg)

		// then
		assert.Nil(t, err)
		assert.False(t, proxySettings.Enabled)
	})

	t.Run("return error on query config if proxy is enabled", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "proxy/enabled").Return(false, assert.AnError)

		// when
		_, err := tasks.CreateProxySettings(reg)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})
}
