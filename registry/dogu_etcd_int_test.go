//go:build integration
// +build integration

package registry

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllReturnsV1AsWellAsV2_inttest(t *testing.T) {
	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	v1Format := &core.DoguJsonV1FormatProvider{}
	v2Format := &core.DoguJsonV2FormatProvider{}

	cl, err := newResilientEtcdClient([]string{server.URL}, core.RetryPolicy{Interval: 100})
	require.Nil(t, err)

	defer func() {
		cleanupRegistry(t, cl)
	}()

	v1reg := &etcdDoguRegistry{
		"dogu",
		cl,
		v1Format,
	}

	v2reg := &etcdDoguRegistry{
		"dogu_v2",
		cl,
		v2Format,
	}

	dogua := createTestDogu(t, "dogua", "1.0.0-1")
	dogub := createTestDogu(t, "dogub", "1.0.0-1")
	doguc := createTestDogu(t, "doguc", "1.0.0-1")
	dogud := createTestDogu(t, "dogud", "1.0.0-1")
	dogue := createTestDogu(t, "dogue", "1.0.0-1")

	dogusv1 := []*core.Dogu{dogua, dogub, doguc, dogud, dogue}
	dogusv2 := []*core.Dogu{dogub, doguc, dogud}

	registerAll(t, dogusv1, v1reg)
	registerAll(t, dogusv2, v2reg)

	combinedReg := newCombinedEtcdDoguRegistry(cl, "dogu", "dogu_v2")

	allDogusResult, err := combinedReg.GetAll()
	assert.Nil(t, err)
	assert.Len(t, allDogusResult, 5)
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogua"))
	assert.Contains(t, allDogusResult, dogub)
	assert.Contains(t, allDogusResult, doguc)
	assert.Contains(t, allDogusResult, dogud)
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogue"))
}

func TestCombinedEtcdDoguRegistryWritesToV1AndV2_inttest(t *testing.T) {
	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	v1Format := &core.DoguJsonV1FormatProvider{}
	v2Format := &core.DoguJsonV2FormatProvider{}

	cl, err := newResilientEtcdClient([]string{server.URL}, core.RetryPolicy{Interval: 100})
	require.Nil(t, err)

	defer func() {
		cleanupRegistry(t, cl)
	}()

	v1reg := &etcdDoguRegistry{
		"dogu",
		cl,
		v1Format,
	}

	v2reg := &etcdDoguRegistry{
		"dogu_v2",
		cl,
		v2Format,
	}

	dogua := createTestDogu(t, "dogua", "1.0.0-1")
	doguav2 := createTestDogu(t, "dogua", "1.0.0-2")
	dogub := createTestDogu(t, "dogub", "1.0.0-1")
	doguc := createTestDogu(t, "doguc", "1.0.0-1")
	dogud := createTestDogu(t, "dogud", "1.0.0-1")
	dogue := createTestDogu(t, "dogue", "1.0.0-1")

	allDogus := []*core.Dogu{dogua, doguav2, dogub, doguc, dogud, dogue}

	combinedReg := newCombinedEtcdDoguRegistry(cl, "dogu", "dogu_v2")
	registerAll(t, allDogus, combinedReg)

	allDogusResult, err := v1reg.GetAll()
	assert.Nil(t, err)
	assert.Len(t, allDogusResult, 5)
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogua"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogub"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/doguc"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogud"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogue"))

	allDogusResult, err = v2reg.GetAll()
	assert.Nil(t, err)
	assert.Len(t, allDogusResult, 5)
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogua"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogub"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/doguc"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogud"))
	assert.True(t, core.ContainsDoguWithName(allDogusResult, "official/dogue"))
}

func registerAll(t *testing.T, dogus []*core.Dogu, registry DoguRegistry) {
	t.Helper()
	for _, dogu := range dogus {
		err := registry.Register(dogu)
		assert.Nil(t, err)
		err = registry.Enable(dogu)
		assert.Nil(t, err)
	}
}

func createTestDogu(t *testing.T, name string, version string) *core.Dogu {
	t.Helper()

	return &core.Dogu{Name: "official/" + name, Version: version, Dependencies: []core.Dependency{{
		Type:    "dogu",
		Name:    name + "dep",
		Version: version,
	}}}

}

func cleanupRegistry(t *testing.T, cl *resilentEtcdClient) {
	exists, err := cl.Exists("dogu")
	assert.Nil(t, err)

	if exists {
		err = cl.DeleteRecursive("dogu")
		assert.Nil(t, err)
	}

	exists, err = cl.Exists("dogu_v2")
	assert.Nil(t, err)

	if exists {
		err = cl.DeleteRecursive("dogu_v2")
		assert.Nil(t, err)
	}
}
