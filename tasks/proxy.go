package tasks

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

// CreateProxySettings reads the proxy related settings from registry.
func CreateProxySettings(reg registry.Registry) (core.ProxySettings, error) {
	settings := core.ProxySettings{}

	configReader := registry.NewConfigurationReader(reg.GlobalConfig())
	enabled, err := configReader.GetBool("proxy/enabled")
	if err != nil {
		return settings, err
	}

	if enabled {
		settings.Enabled = true

		server, err := configReader.GetString("proxy/server")
		if err != nil {
			return settings, err
		}
		settings.Server = server

		port, err := configReader.GetInt("proxy/port")
		if err != nil {
			return settings, err
		}
		settings.Port = port

		username, err := configReader.GetString("proxy/username")
		if err != nil {
			return settings, err
		}
		settings.Username = username

		password, err := configReader.GetString("proxy/password")
		if err != nil {
			return settings, err
		}
		settings.Password = password
	}

	return settings, nil
}
