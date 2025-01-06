package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolume_GetClient(t *testing.T) {
	t.Run("call on volumes with no client definitions", func(t *testing.T) {
		// given
		sut := Volume{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
		}

		// when
		client, ok := sut.GetClient("testClient")

		// then
		assert.False(t, ok)
		assert.Nil(t, client)
	})
	t.Run("call on volumes with a single client definitions", func(t *testing.T) {
		// given
		type testClientParams struct {
			Type     string
			MySecret string
		}

		sut := Volume{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
			Clients: []VolumeClient{
				{Name: "testClient", Params: testClientParams{MySecret: "supersecret", Type: "myType"}},
			},
		}

		// when
		client, ok := sut.GetClient("testClient")

		// then
		assert.True(t, ok)
		require.NotNil(t, client)
		assert.Equal(t, "testClient", client.Name)

		params, ok := client.Params.(testClientParams)
		require.True(t, ok)

		assert.Equal(t, "supersecret", params.MySecret)
		assert.Equal(t, "myType", params.Type)
	})
	t.Run("call on volumes with multiple client definitions", func(t *testing.T) {
		// given
		type testClientParams struct {
			Type     string
			MySecret string
		}

		sut := Volume{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
			Clients: []VolumeClient{
				{Name: "wrongClient", Params: "wrong data"},
				{Name: "testClient", Params: testClientParams{MySecret: "supersecret", Type: "myType"}},
			},
		}

		// when
		client, ok := sut.GetClient("testClient")

		// then
		assert.True(t, ok)
		require.NotNil(t, client)
		assert.Equal(t, "testClient", client.Name)

		params, ok := client.Params.(testClientParams)
		require.True(t, ok)

		assert.Equal(t, "supersecret", params.MySecret)
		assert.Equal(t, "myType", params.Type)
	})
	t.Run("return false when requested client does not exits in multiple available clients", func(t *testing.T) {
		// given
		type testClientParams struct {
			Type     string
			MySecret string
		}

		sut := Volume{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
			Clients: []VolumeClient{
				{Name: "wrongClient", Params: "wrong data"},
				{Name: "testClient", Params: testClientParams{MySecret: "supersecret", Type: "myType"}},
			},
		}

		// when
		client, ok := sut.GetClient("againAnotherClient")

		// then
		assert.False(t, ok)
		require.Nil(t, client)
	})
}
