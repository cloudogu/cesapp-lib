package credentials

import (
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
