package credentials

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func TestCredentialStore(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "credential-test")
	assert.Nil(t, err)

	store, err := newSimpleStore(dir)
	assert.Nil(t, err)
	test := store.Get("test")
	assert.Nil(t, test)

	err = store.Add("test", &core.Credentials{Username: "hans", Password: "special"})
	assert.Nil(t, err)

	test = store.Get("test")
	assert.NotNil(t, test)
	assert.Equal(t, "hans", test.Username)
	assert.Equal(t, "special", test.Password)

	store, err = newSimpleStore(dir)
	assert.Nil(t, err)

	test = store.Get("test")
	assert.NotNil(t, test)
	assert.Equal(t, "hans", test.Username)
	assert.Equal(t, "special", test.Password)

	err = store.Remove("test")
	assert.Nil(t, err)

	test = store.Get("test")
	assert.Nil(t, test)

	store, err = newSimpleStore(dir)
	assert.Nil(t, err)

	test = store.Get("test")
	assert.Nil(t, test)

	err = os.RemoveAll(dir)
	assert.Nil(t, err)
}

func TestDefaultStoreName(t *testing.T) {
	assert.Equal(t, "_default", DefaultStore)
}

func TestNewStore(t *testing.T) {
	// when
	store, err := NewStore(os.TempDir())

	// then
	require.NoError(t, err)
	require.NotNil(t, store)
}

func Test_readStore(t *testing.T) {
	t.Run("failed to read store", func(t *testing.T) {
		_, err := readStore(nil, fmt.Sprintf("%s/does not exist", os.TempDir()))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read credential store")
	})

	t.Run("should return error because to short cyphertext", func(t *testing.T) {
		// given
		file := fmt.Sprintf("%s/cyphertext", os.TempDir())
		err := os.WriteFile(file, []byte("test"), 0644)
		require.NoError(t, err)

		defer func(file string) {
			err = os.Remove(file)
			require.NoError(t, err)
		}(file)

		// when
		_, err = readStore(nil, file)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ciphertext too short")
	})

	t.Run("should return error on invalid key", func(t *testing.T) {
		// given
		file := fmt.Sprintf("%s/cyphertext", os.TempDir())
		err := os.WriteFile(file, []byte("1111111111111111111111111111111111111111111111"), 0644)
		require.NoError(t, err)

		defer func(file string) {
			err = os.Remove(file)
			require.NoError(t, err)
		}(file)

		// when
		_, err = readStore([]byte("test"), file)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create cipher from secret key")
	})
}
