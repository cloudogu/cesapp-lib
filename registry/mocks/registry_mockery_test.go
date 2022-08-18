package mocks_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"github.com/stretchr/testify/require"
)

func TestInitializesCorrect(t *testing.T) {
	r := mocks.CreateMockRegistry([]string{
		"myDogu1",
		"myDogu2",
	})
	reg := r.Registry
	registries := r.SubRegistries
	doguReg := r.DoguRegistry

	doguReg.On("Get", mock.Anything).Return(&core.Dogu{Name: "testdogu"}, nil)
	registries["_global"].On("Get", mock.Anything).Return("testglobalkey", nil)
	mocks.OnGet(registries["blueprints"], mock.Anything, "testbpkey", nil)
	mocks.OnGet(registries["myDogu1"], "test", "testvalue1", nil)
	mocks.OnGet(registries["myDogu2"], mock.Anything, "testvalue2", nil)

	doguRegistry := reg.DoguRegistry()
	require.NotNil(t, doguRegistry)

	testdogu, err := doguRegistry.Get("testdogu")
	require.NoError(t, err)
	require.Equal(t, "testdogu", testdogu.Name)

	globalReg := reg.GlobalConfig()
	require.NotNil(t, globalReg)
	val, err := globalReg.Get("test")
	require.NoError(t, err)
	require.Equal(t, "testglobalkey", val)

	blueprintReg := reg.BlueprintRegistry()
	require.NotNil(t, blueprintReg)
	val, err = blueprintReg.Get("test")
	require.NoError(t, err)
	require.Equal(t, "testbpkey", val)

	myDogu1 := reg.DoguConfig("myDogu1")
	require.NotNil(t, myDogu1)
	val, err = reg.DoguConfig("myDogu1").Get("test")
	require.Nil(t, err)
	require.Equal(t, "testvalue1", val)

	myDogu2 := reg.DoguConfig("myDogu2")
	require.NotNil(t, myDogu2)
	val, err = reg.DoguConfig("myDogu2").Get("test")
	require.Nil(t, err)
	require.Equal(t, "testvalue2", val)

	rootConf := reg.RootConfig()
	require.NotNil(t, rootConf)

	node := reg.GetNode()
	require.NotNil(t, node)

	reg.AssertExpectations(t)
	for _, r := range registries {
		r.AssertExpectations(t)
	}
	doguReg.AssertExpectations(t)
}

func TestCanPassNilWithoutError(t *testing.T) {
	_ = mocks.CreateMockRegistry(nil)
}

func TestOn(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnGet(registries["_global"], "myKey", "myValue", nil)
		val, err := reg.GlobalConfig().Get("myKey")
		require.NoError(t, err)
		require.Equal(t, "myValue", val)
	})

	t.Run("Get with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnGet(registries["_global"], mock.Anything, "myValue", nil)
		val, err := reg.GlobalConfig().Get("myKey")
		require.NoError(t, err)
		require.Equal(t, "myValue", val)
	})

	t.Run("Set", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSet(registries["_global"], "myKey", "myValue", nil)
		err := reg.GlobalConfig().Set("myKey", "myValue")
		require.NoError(t, err)
	})

	t.Run("Set with anything as key", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSet(registries["_global"], mock.Anything, "myValue", nil)
		err := reg.GlobalConfig().Set("myKey", "myValue")
		require.NoError(t, err)
	})

	t.Run("Set with anything as value", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSet(registries["_global"], "myKey", mock.Anything, nil)
		err := reg.GlobalConfig().Set("myKey", "myValue")
		require.NoError(t, err)
	})

	t.Run("Set with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSet(registries["_global"], mock.Anything, mock.Anything, nil)
		err := reg.GlobalConfig().Set("myKey", "myValue")
		require.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnDelete(registries["_global"], "testkey", nil)
		err := reg.GlobalConfig().Delete("testkey")
		require.Nil(t, err)
	})

	t.Run("Delete with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnDelete(registries["_global"], mock.Anything, nil)
		err := reg.GlobalConfig().Delete("testkey")
		require.Nil(t, err)
	})

	t.Run("DeleteRecursive", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnDeleteRecursive(registries["_global"], "testkey", nil)
		err := reg.GlobalConfig().DeleteRecursive("testkey")
		require.Nil(t, err)
	})

	t.Run("DeleteRecursive with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnDeleteRecursive(registries["_global"], mock.Anything, nil)
		err := reg.GlobalConfig().DeleteRecursive("testkey")
		require.Nil(t, err)
	})

	t.Run("Exists", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnExists(registries["_global"], "testkey", true, nil)
		exists, err := reg.GlobalConfig().Exists("testkey")
		require.True(t, exists)
		require.Nil(t, err)
	})
	t.Run("Exists with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnExists(registries["_global"], mock.Anything, true, nil)
		exists, err := reg.GlobalConfig().Exists("testkey")
		require.True(t, exists)
		require.Nil(t, err)
	})
	t.Run("GetOrFalse", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnGetOrFalse(registries["_global"], "testkey", "testval", true, nil)
		exists, value, err := reg.GlobalConfig().GetOrFalse("testkey")
		require.True(t, exists)
		require.Equal(t, "testval", value)
		require.Nil(t, err)
	})
	t.Run("GetOrFalse with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnGetOrFalse(registries["_global"], mock.Anything, "testval", true, nil)
		exists, value, err := reg.GlobalConfig().GetOrFalse("testkey")
		require.True(t, exists)
		require.Equal(t, "testval", value)
		require.Nil(t, err)
	})
	t.Run("Refresh", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnRefresh(registries["_global"], "testkey", 5, nil)
		err := reg.GlobalConfig().Refresh("testkey", 5)
		require.Nil(t, err)
	})
	t.Run("Refresh with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnRefresh(registries["_global"], mock.Anything, mocks.AnyLifetime, nil)
		err := reg.GlobalConfig().Refresh("testkey", 5)
		require.Nil(t, err)
	})
	t.Run("SetWithLifetime", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSetWithLifetime(registries["_global"], "testkey", "testvalue", 5, nil)
		err := reg.GlobalConfig().SetWithLifetime("testkey", "testvalue", 5)
		require.Nil(t, err)
	})
	t.Run("SetWithLifetime with key anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSetWithLifetime(registries["_global"], mock.Anything, "testvalue", 5, nil)
		err := reg.GlobalConfig().SetWithLifetime("testkey", "testvalue", 5)
		require.Nil(t, err)
	})
	t.Run("SetWithLifetime with value anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSetWithLifetime(registries["_global"], "testkey", mock.Anything, 5, nil)
		err := reg.GlobalConfig().SetWithLifetime("testkey", "testvalue", 5)
		require.Nil(t, err)
	})
	t.Run("SetWithLifetime with anything", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnSetWithLifetime(registries["_global"], mock.Anything, mock.Anything, mocks.AnyLifetime, nil)
		err := reg.GlobalConfig().SetWithLifetime("testkey", "testvalue", 5)
		require.Nil(t, err)
	})
	t.Run("RemoveAll", func(t *testing.T) {
		r := mocks.CreateMockRegistry(nil)
		reg := r.Registry
		registries := r.SubRegistries
		mocks.OnRemoveAll(registries["_global"], nil)
		err := reg.GlobalConfig().RemoveAll()
		require.Nil(t, err)
	})
}
