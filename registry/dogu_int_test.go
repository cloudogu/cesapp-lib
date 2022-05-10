//go:build integration
// +build integration

package registry_test

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func TestDogu(t *testing.T) {
	doguReg := reg.DoguRegistry()

	enabled, err := doguReg.IsEnabled("test")
	assert.Nil(t, err)
	assert.False(t, enabled)

	dogu, err := doguReg.Get("test")
	assert.NotNil(t, err)

	dogu = &core.Dogu{Name: "test", Version: "42.0"}
	err = doguReg.Register(dogu)
	assert.Nil(t, err)

	err = doguReg.Enable(dogu)
	assert.Nil(t, err)

	enabled, err = doguReg.IsEnabled("test")
	assert.Nil(t, err)
	assert.True(t, enabled)

	returnedDogu, err := doguReg.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, "42.0", returnedDogu.Version, "Version should be equal")

	err = doguReg.Unregister("test")
	assert.Nil(t, err)

	_, err = doguReg.Get("test")
	assert.NotNil(t, err)
}

func TestDogus(t *testing.T) {
	doguReg := reg.DoguRegistry()

	dogu := &core.Dogu{Name: "test-1", Version: "42.0"}
	err := doguReg.Register(dogu)
	assert.Nil(t, err)
	err = doguReg.Enable(dogu)
	assert.Nil(t, err)

	dogu = &core.Dogu{Name: "test-2", Version: "42.0"}
	err = doguReg.Register(dogu)
	assert.Nil(t, err)
	err = doguReg.Enable(dogu)
	assert.Nil(t, err)

	dogus, err := doguReg.GetAll()
	assert.Nil(t, err)
	assert.True(t, len(dogus) >= 2)

	err = doguReg.Unregister("test-1")
	assert.Nil(t, err)

	err = doguReg.Unregister("test-2")
	assert.Nil(t, err)
}
