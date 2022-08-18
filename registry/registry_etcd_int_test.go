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

func Test_mapEtcdNodeToRegistryNode(t *testing.T) {
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
	}()

	_, err = client.Set("/dir_test/key1/subkey1", "val1", nil)
	require.Nil(t, err)

	_, err = client.Set("/dir_test/key1/subkey2", "val2", nil)
	require.Nil(t, err)

	_, err = client.Set("/dir_test/key2", "val3", nil)
	require.Nil(t, err)

	node, err := client.getMainNode()
	require.NoError(t, err)

	result := mapEtcdNodeToRegistryNode(node, nil)

	assert.Nil(t, result.GetParent())

	assert.GreaterOrEqual(t, len(result.GetSubNodes()), 1)
	assert.Len(t, result.GetSubNode("dir_test").GetSubNodes(), 2)
	assert.Len(t, result.GetSubNode("dir_test").GetSubNode("key1").GetSubNodes(), 2)
	assert.Len(t, result.GetSubNode("dir_test").GetSubNode("key2").GetSubNodes(), 0)

	t.Run("result is correctly setup", func(t *testing.T) {
		assert.Equal(t, "", result.GetKey())
		assert.Equal(t, "dir_test", result.GetSubNode("dir_test").GetKey())
		assert.Equal(t, "dir_test", result.GetSubNode("dir_test").GetKey())
		assert.Equal(t, "", result.GetSubNode("dir_test").GetParent().GetKey())
	})

	t.Run("keys are set correctly", func(t *testing.T) {
		assert.Equal(t, "key1", result.GetSubNode("dir_test").GetSubNode("key1").GetKey())
		assert.Equal(t, "/dir_test/key1", result.GetSubNode("dir_test").GetSubNode("key1").GetFullKey())
		assert.Equal(t, "/dir_test/key1/subkey1", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").GetFullKey())
		assert.Equal(t, "subkey1", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").GetKey())
		assert.Equal(t, "val1", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").GetValue())
		assert.Equal(t, "/dir_test/key1/subkey2", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey2").GetFullKey())
		assert.Equal(t, "subkey2", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey2").GetKey())
		assert.Equal(t, "val2", result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey2").GetValue())
		assert.Equal(t, "/dir_test/key2", result.GetSubNode("dir_test").GetSubNode("key2").GetFullKey())
		assert.Equal(t, "key2", result.GetSubNode("dir_test").GetSubNode("key2").GetKey())
		assert.Equal(t, "val3", result.GetSubNode("dir_test").GetSubNode("key2").GetValue())
	})

	t.Run("IsDir returns correct value", func(t *testing.T) {
		assert.True(t, result.IsDir())
		assert.True(t, result.GetSubNode("dir_test").IsDir())
		assert.True(t, result.GetSubNode("dir_test").GetSubNode("key1").IsDir())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").IsDir())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey2").IsDir())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key2").IsDir())
	})

	t.Run("HasSubNodes works", func(t *testing.T) {
		assert.True(t, result.HasSubNodes())
		assert.True(t, result.GetSubNode("dir_test").HasSubNodes())
		assert.True(t, result.GetSubNode("dir_test").GetSubNode("key1").HasSubNodes())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").HasSubNodes())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey2").HasSubNodes())
		assert.False(t, result.GetSubNode("dir_test").GetSubNode("key2").HasSubNodes())
	})

	t.Run("can follow parent-child flow", func(t *testing.T) {
		deepestNode := result.GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1")
		assert.Equal(t, "/dir_test/key1/subkey1", deepestNode.GetFullKey())
		assert.Equal(t, "/dir_test/key1", deepestNode.GetParent().GetFullKey())
		assert.Equal(t, "key1", deepestNode.GetParent().GetKey())
		assert.Equal(t, "/dir_test", deepestNode.GetParent().GetParent().GetFullKey())
		assert.Equal(t, "dir_test", deepestNode.GetParent().GetParent().GetKey())
		assert.Equal(t, "", deepestNode.GetParent().GetParent().GetParent().GetFullKey())
		assert.Equal(t, "", deepestNode.GetParent().GetParent().GetParent().GetKey())
		assert.Equal(t,
			"/dir_test/key1/subkey1",
			deepestNode.GetParent().GetParent().GetParent().GetSubNode("dir_test").GetSubNode("key1").GetSubNode("subkey1").GetFullKey())
	})

}
