//go:build integration
// +build integration

package registry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
)

func TestGlobalConfig_inttest(t *testing.T) {
	testConfigurationContext(t, reg.GlobalConfig())
}

func TestDoguConfig_inttest(t *testing.T) {
	testConfigurationContext(t, reg.DoguConfig("unit-test-1"))
}

func TestDoguConfigRemoveAll_inttest(t *testing.T) {
	cc := reg.DoguConfig("unit-test-2")

	err := cc.Set("key-1", "Value-1")
	assert.Nil(t, err)

	err = cc.Set("key-2", "Value-2")
	assert.Nil(t, err)

	err = cc.Set("key-3", "Value-3")
	assert.Nil(t, err)

	err = cc.RemoveAll()
	assert.Nil(t, err)

	exists, err := cc.Exists("key-1")
	assert.Nil(t, err)
	assert.False(t, exists)

	exists, err = cc.Exists("key-2")
	assert.Nil(t, err)
	assert.False(t, exists)

	exists, err = cc.Exists("key-3")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestEtcdDoguConfigGetAll_inttest(t *testing.T) {
	cc := reg.DoguConfig("unit-test-3")
	defer cc.RemoveAll()

	err := cc.Set("key-1", "Value-1")
	assert.Nil(t, err)

	err = cc.Set("key-2", "Value-2")
	assert.Nil(t, err)

	err = cc.Set("key-3", "Value-3")
	assert.Nil(t, err)

	err = cc.Set("keys/4", "Value-4")
	assert.Nil(t, err)

	keyValuePairs, err := cc.GetAll()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(keyValuePairs))
	assert.Equal(t, "Value-1", keyValuePairs["key-1"])
	assert.Equal(t, "Value-2", keyValuePairs["key-2"])
	assert.Equal(t, "Value-3", keyValuePairs["key-3"])
	assert.Equal(t, "Value-4", keyValuePairs["keys/4"])

	err = cc.RemoveAll()
	assert.Nil(t, err)
}

func testConfigurationContext(t *testing.T, cc registry.ConfigurationContext) {
	t.Helper()
	defer cc.RemoveAll()

	exists, err := cc.Exists("key-1")
	assert.Nil(t, err)
	assert.False(t, exists)

	exists2, value2, err := cc.GetOrFalse("key-1")
	assert.False(t, exists2)
	assert.Empty(t, value2)
	assert.NoError(t, err)

	err = cc.Set("key-1", "Value-1")
	assert.Nil(t, err)

	err = cc.Set("dir1/key1", "Value-1")
	assert.Nil(t, err)

	exists, err = cc.Exists("key-1")
	assert.Nil(t, err)
	assert.True(t, exists)

	value, err := cc.Get("key-1")
	assert.Nil(t, err)
	assert.Equal(t, "Value-1", value)

	exists2, value2, err = cc.GetOrFalse("key-1")
	assert.True(t, exists2)
	assert.Equal(t, "Value-1", value2)
	assert.NoError(t, err)

	err = cc.Delete("key-1")
	assert.Nil(t, err)

	err = cc.Delete("dir1")
	assert.Error(t, err)

	err = cc.DeleteRecursive("dir1")
	assert.NoError(t, err)

	exists, err = cc.Exists("key-1")
	assert.Nil(t, err)
	assert.False(t, exists)

	// Test ttl
	ttl := 5
	require.Nil(t, err)

	err = cc.SetWithLifetime("key-2", "Value-2", ttl)
	assert.Nil(t, err)

	exists, err = cc.Exists("key-2")
	assert.Nil(t, err)
	assert.True(t, exists)

	value, err = cc.Get("key-2")
	assert.Nil(t, err)
	assert.Equal(t, "Value-2", value)

	// Refresh to have the maximum ttl
	err = cc.Refresh("key-2", ttl)
	require.Nil(t, err)

	// Wait ttl-2 seconds
	refreshWaitDuration := ttl - 2
	refreshWaitDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", refreshWaitDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)

	err = cc.Refresh("key-2", ttl)
	require.Nil(t, err)

	// Wait again ttl-2 seconds and make sure that the Value still exists
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)
	value, err = cc.Get("key-2")
	require.Nil(t, err)
	require.Equal(t, "Value-2", value)

	// Wait until expiration
	expireDuration := ttl + 1
	expireDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", expireDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", expireDuration)
	time.Sleep(expireDurationParsed)

	exists, err = cc.Exists("key-2")
	require.Nil(t, err)
	require.False(t, exists)
}
