package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProxySettings_CreateURL(t *testing.T) {
	settings := ProxySettings{Enabled: true, Server: "proxy.cloudogu.com", Port: 3182}
	assert.Equal(t, "http://proxy.cloudogu.com:3182", settings.CreateURL())
}
