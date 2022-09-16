//go:build integration
// +build integration

package registry_test

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
	"testing"

	"os"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
)

var reg registry.Registry

func init() {
	reg, _ = registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: createTestEndpoint(),
	})
}

func TestGetNode_inttest(t *testing.T) {
	err := reg.GlobalConfig().Set("test/key", "false")
	require.Nil(t, err)
	defer func() {
		_ = reg.GlobalConfig().DeleteRecursive("test")
	}()

	node, err := reg.GetNode()
	require.Nil(t, err)

	assert.True(t, node.IsDir)
	assert.Equal(t, "", node.FullKey)

	testkey := node.SubNodeByName("config").SubNodeByName("_global").SubNodeByName("test").SubNodeByName("key")
	assert.Equal(t, "false", testkey.Value)
	assert.Equal(t, "key", testkey.Key())
	assert.Equal(t, "/config/_global/test/key", testkey.FullKey)
}

func TestNew_inttest(t *testing.T) {
	_, err := registry.New(core.Registry{
		Type:      "consul",
		Endpoints: createTestEndpoint(),
	})
	assert.NotNil(t, err)
}

func createTestEndpoint() []string {
	etcd := os.Getenv("ETCD")
	if etcd == "" {
		etcd = "localhost"
	}
	return []string{"http://" + etcd + ":4001"}
}
