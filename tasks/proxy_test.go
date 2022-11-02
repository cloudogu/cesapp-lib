package tasks_test

import (
	"github.com/cloudogu/cesapp-lib/tasks"
	"testing"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
)

func TestCreateProxySettings(t *testing.T) {
	registry := &registry.MockRegistry{}

	registry.GlobalConfig().Set("proxy/enabled", "true")
	registry.GlobalConfig().Set("proxy/server", "proxy.cloudogu.com")
	registry.GlobalConfig().Set("proxy/port", "3864")
	registry.GlobalConfig().Set("proxy/username", "bob")
	registry.GlobalConfig().Set("proxy/password", "bob123")

	proxySettings, err := tasks.CreateProxySettings(registry)
	assert.Nil(t, err)
	assert.True(t, proxySettings.Enabled)
	assert.Equal(t, "proxy.cloudogu.com", proxySettings.Server)
	assert.Equal(t, 3864, proxySettings.Port)
	assert.Equal(t, "bob", proxySettings.Username)
	assert.Equal(t, "bob123", proxySettings.Password)
}

func TestCreateProxySettingsWithoutConfiguration(t *testing.T) {
	registry := &registry.MockRegistry{}
	proxySettings, err := tasks.CreateProxySettings(registry)
	assert.Nil(t, err)
	assert.False(t, proxySettings.Enabled)
}
