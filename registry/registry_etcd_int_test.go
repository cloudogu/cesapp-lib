//go:build integration
// +build integration

package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_mapEtcdNodeToRegistryNode_inttest(t *testing.T) {
	// create etcd address, for local execution and on ci
	etcd := os.Getenv("ETCD")
	if etcd == "" {
		etcd = "localhost"
	}

	client, err := createEtcdClient(core.Registry{
		Type: "etcd",
		Endpoints: []string{
			"http://" + etcd + ":4001",
		},
	})

	defer func() {
		_ = client.DeleteRecursive("/dir_test")
		_ = client.DeleteRecursive("/config/_global")
	}()

	_, err = client.Set("/dir_test/key1/subkey1", "val1", nil)
	require.Nil(t, err)

	_, err = client.Set("/dir_test/key1/subkey2", "val2", nil)
	require.Nil(t, err)

	_, err = client.Set("/dir_test/key2", "val3", nil)
	require.Nil(t, err)

	_, err = client.Set("/config/_global/key", "val4", nil)
	require.Nil(t, err)

	node, err := client.getMainNode()
	require.NoError(t, err)

	result := mapEtcdNodeToRegistryNode(node)

	assert.GreaterOrEqual(t, len(result.SubNodes), 1)
	assert.Len(t, result.SubNodeByName("dir_test").SubNodes, 2)
	assert.Len(t, result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodes, 2)
	assert.Len(t, result.SubNodeByName("dir_test").SubNodeByName("key2").SubNodes, 0)

	t.Run("global config is included", func(t *testing.T) {
		assert.Equal(t, "/config/_global", result.SubNodeByName("config").SubNodeByName("_global").FullKey)
		assert.Equal(t, "/config/_global/key", result.SubNodeByName("config").SubNodeByName("_global").SubNodeByName("key").FullKey)
	})

	t.Run("result is correctly setup", func(t *testing.T) {
		assert.Equal(t, "", result.Key())
		assert.Equal(t, "dir_test", result.SubNodeByName("dir_test").Key())
		assert.Equal(t, "dir_test", result.SubNodeByName("dir_test").Key())
	})

	t.Run("keys are set correctly", func(t *testing.T) {
		assert.Equal(t, "key1", result.SubNodeByName("dir_test").SubNodeByName("key1").Key())
		assert.Equal(t, "/dir_test/key1", result.SubNodeByName("dir_test").SubNodeByName("key1").FullKey)
		assert.Equal(t, "/dir_test/key1/subkey1", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey1").FullKey)
		assert.Equal(t, "subkey1", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey1").Key())
		assert.Equal(t, "val1", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey1").Value)
		assert.Equal(t, "/dir_test/key1/subkey2", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey2").FullKey)
		assert.Equal(t, "subkey2", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey2").Key())
		assert.Equal(t, "val2", result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey2").Value)
		assert.Equal(t, "/dir_test/key2", result.SubNodeByName("dir_test").SubNodeByName("key2").FullKey)
		assert.Equal(t, "key2", result.SubNodeByName("dir_test").SubNodeByName("key2").Key())
		assert.Equal(t, "val3", result.SubNodeByName("dir_test").SubNodeByName("key2").Value)
	})

	t.Run("IsDir returns correct Value", func(t *testing.T) {
		assert.True(t, result.IsDir)
		assert.True(t, result.SubNodeByName("dir_test").IsDir)
		assert.True(t, result.SubNodeByName("dir_test").SubNodeByName("key1").IsDir)
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey1").IsDir)
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey2").IsDir)
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key2").IsDir)
	})

	t.Run("HasSubNodes works", func(t *testing.T) {
		assert.True(t, result.HasSubNodes())
		assert.True(t, result.SubNodeByName("dir_test").HasSubNodes())
		assert.True(t, result.SubNodeByName("dir_test").SubNodeByName("key1").HasSubNodes())
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey1").HasSubNodes())
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key1").SubNodeByName("subkey2").HasSubNodes())
		assert.False(t, result.SubNodeByName("dir_test").SubNodeByName("key2").HasSubNodes())
	})
}
