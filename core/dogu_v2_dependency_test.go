package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependsOn(t *testing.T) {
	dogu := Dogu{
		Name: "hansolo",
		Dependencies: []Dependency{
			{Type: DependencyTypeDogu, Name: "a"},
			{Type: DependencyTypeDogu, Name: "b"},
		},
		OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
	}
	assert.True(t, dogu.DependsOn("a"))
	assert.True(t, dogu.DependsOn("b"))
	assert.False(t, dogu.DependsOn("c"))
	assert.False(t, dogu.DependsOn("d"))
}

func TestDogu_Dependencies(t *testing.T) {
	// given
	depDogu1 := Dependency{Type: DependencyTypeDogu, Name: "DoguTest"}
	depDogu2 := Dependency{Type: DependencyTypeDogu, Name: "DoguTest2"}
	depClient1 := Dependency{Type: DependencyTypeClient, Name: "ClientTest1"}
	depClient2 := Dependency{Type: DependencyTypeClient, Name: "ClientTest2"}
	depPackage1 := Dependency{Type: DependencyTypePackage, Name: "PackageTest1"}
	depPackage2 := Dependency{Type: DependencyTypePackage, Name: "PackageTest2"}

	t.Run("GetAllDependenciesOfType - Return only empty slice when dependencies are empty", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name: "namespace/dogu1",
		}

		// when
		dependenciesDogu := dogu1.GetAllDependenciesOfType(DependencyTypeDogu)

		// then
		assert.Empty(t, dependenciesDogu)
	})

	t.Run("GetDependenciesOfType - Return only empty slice when dependencies are empty", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name: "namespace/dogu1",
		}

		// when
		dependenciesDogu := dogu1.GetDependenciesOfType(DependencyTypeDogu)

		// then
		assert.Empty(t, dependenciesDogu)
	})

	t.Run("GetOptionalDependenciesOfType - Return only empty slice when dependencies are empty", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name: "namespace/dogu1",
		}

		// when
		dependenciesDogu := dogu1.GetOptionalDependenciesOfType(DependencyTypeDogu)

		// then
		assert.Empty(t, dependenciesDogu)
	})

	t.Run("GetAllDependenciesOfType - Return all dependencies with their correct types", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name:                 "namespace/dogu1",
			Dependencies:         []Dependency{depDogu1, depClient1, depPackage1},
			OptionalDependencies: []Dependency{depDogu2, depClient2, depPackage2},
		}

		// when
		dependenciesDogu := dogu1.GetAllDependenciesOfType(DependencyTypeDogu)
		dependenciesClient := dogu1.GetAllDependenciesOfType(DependencyTypeClient)
		dependenciesPackage := dogu1.GetAllDependenciesOfType(DependencyTypePackage)

		// then
		assert.NotEmpty(t, dependenciesDogu)
		assert.Len(t, dependenciesDogu, 2)
		assert.Equal(t, depDogu1, dependenciesDogu[0])
		assert.Equal(t, depDogu2, dependenciesDogu[1])
		assert.Len(t, dependenciesClient, 2)
		assert.Equal(t, depClient1, dependenciesClient[0])
		assert.Equal(t, depClient2, dependenciesClient[1])
		assert.Len(t, dependenciesPackage, 2)
		assert.Equal(t, depPackage1, dependenciesPackage[0])
		assert.Equal(t, depPackage2, dependenciesPackage[1])
	})

	t.Run("GetDependenciesOfType - Return only required dependencies with their correct types", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name:                 "namespace/dogu1",
			Dependencies:         []Dependency{depDogu1, depClient1, depPackage1},
			OptionalDependencies: []Dependency{depDogu2, depClient2, depPackage2},
		}

		// when
		dependenciesDogu := dogu1.GetDependenciesOfType(DependencyTypeDogu)
		dependenciesClient := dogu1.GetDependenciesOfType(DependencyTypeClient)
		dependenciesPackage := dogu1.GetDependenciesOfType(DependencyTypePackage)

		// then
		assert.Len(t, dependenciesDogu, 1)
		assert.Equal(t, depDogu1, dependenciesDogu[0])
		assert.Len(t, dependenciesClient, 1)
		assert.Equal(t, depClient1, dependenciesClient[0])
		assert.Len(t, dependenciesPackage, 1)
		assert.Equal(t, depPackage1, dependenciesPackage[0])
	})

	t.Run("GetOptionalDependenciesOfType - Return all optional dependencies with their correct types", func(t *testing.T) {
		// given
		dogu1 := &Dogu{
			Name:                 "namespace/dogu1",
			Dependencies:         []Dependency{depDogu1, depClient1, depPackage1},
			OptionalDependencies: []Dependency{depDogu2, depClient2, depPackage2},
		}

		// when
		dependenciesDogu := dogu1.GetOptionalDependenciesOfType(DependencyTypeDogu)
		dependenciesClient := dogu1.GetOptionalDependenciesOfType(DependencyTypeClient)
		dependenciesPackage := dogu1.GetOptionalDependenciesOfType(DependencyTypePackage)

		// then
		assert.Len(t, dependenciesDogu, 1)
		assert.Equal(t, depDogu2, dependenciesDogu[0])
		assert.Len(t, dependenciesClient, 1)
		assert.Equal(t, depClient2, dependenciesClient[0])
		assert.Len(t, dependenciesPackage, 1)
		assert.Equal(t, depPackage2, dependenciesPackage[0])
	})
}
