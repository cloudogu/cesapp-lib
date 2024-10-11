package core

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSortByDependency(t *testing.T) {
	dogus := []*Dogu{}
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-009.json")
	require.NoError(t, err)
	ordered := SortDogusByDependency(dogus)
	assert.Equal(t, "gotenberg", ordered[0].GetSimpleName())
	assert.Equal(t, "postfix", ordered[1].GetSimpleName())
	assert.Equal(t, "nginx", ordered[2].GetSimpleName())
	assert.Equal(t, "cas", ordered[3].GetSimpleName())
	assert.Equal(t, "backup", ordered[4].GetSimpleName())
	assert.Equal(t, "nexus", ordered[5].GetSimpleName())
	assert.Equal(t, "redmine", ordered[6].GetSimpleName())
	assert.Equal(t, "scm", ordered[7].GetSimpleName())
	assert.Equal(t, "smeagol", ordered[8].GetSimpleName())
	assert.Equal(t, "sonar", ordered[9].GetSimpleName())
	assert.Equal(t, "usermgt", ordered[10].GetSimpleName())
}

func TestSortByDependencyWithSmallList(t *testing.T) {
	dogus := []*Dogu{}
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-002.json")
	assert.Nil(t, err)
	ordered := SortDogusByDependency(dogus)

	assert.Equal(t, "registrator", ordered[0].GetSimpleName())
	assert.Equal(t, "nginx", ordered[1].GetSimpleName())
	assert.Equal(t, "icoordinator", ordered[2].GetSimpleName())
}

func TestSortByDependencyWithTransitiveDependencies(t *testing.T) {
	// a needs b in 1.0-2 and nginx
	// b needs c in 1.0-2 and nginx
	// c needs e in 1.0-2 and nginx and optional d in 1.0-2
	// e needs nginx
	// d needs nginx
	// Dependency tree:
	// 					(opt)-> d
	// 					/		|
	// a ---> b ---> c ---> e	|
	// |	  |		 |		|	|
	// |	  |		 |		|	|
	// --------------------------->nginx
	dogus := []*Dogu{}
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-003.json")
	assert.Nil(t, err)

	ordered := SortDogusByDependency(dogus)

	assert.Equal(t, "dogud", ordered[0].GetSimpleName())
	assert.Equal(t, "dogue", ordered[1].GetSimpleName())
	assert.Equal(t, "doguc", ordered[2].GetSimpleName())
	assert.Equal(t, "dogub", ordered[3].GetSimpleName())
	assert.Equal(t, "dogua", ordered[4].GetSimpleName())
}

func TestSortByDependencyWithOptionalDependencies(t *testing.T) {
	// Dependency tree:
	// d ----> a ----> b
	//          \_    ^
	//            \   |
	//             -> c
	a := &Dogu{
		Name:                 "a",
		Dependencies:         []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
	}

	b := &Dogu{Name: "b"}

	c := &Dogu{
		Name:                 "c",
		OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
	}

	d := &Dogu{
		Name:                 "d",
		OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
	}

	dogus := SortDogusByDependency([]*Dogu{a, b, c, d})
	assert.Equal(t, "b", dogus[0].Name)
	assert.Equal(t, "c", dogus[1].Name)
	assert.Equal(t, "a", dogus[2].Name)
	assert.Equal(t, "d", dogus[3].Name)
}

func TestSortByName(t *testing.T) {
	dogus := make([]*Dogu, 0)
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-001.json")
	assert.Nil(t, err)
	ordered := SortDogusByName(dogus)
	assert.Equal(t, "cas", ordered[0].GetSimpleName())
	assert.Equal(t, "jenkins", ordered[1].GetSimpleName())
	assert.Equal(t, "ldap", ordered[2].GetSimpleName())
	assert.Equal(t, "mysql", ordered[3].GetSimpleName())
	assert.Equal(t, "nexus", ordered[4].GetSimpleName())
	assert.Equal(t, "nginx", ordered[5].GetSimpleName())
	assert.Equal(t, "postfix", ordered[6].GetSimpleName())
	assert.Equal(t, "redmine", ordered[7].GetSimpleName())
	assert.Equal(t, "registrator", ordered[8].GetSimpleName())
	assert.Equal(t, "scm", ordered[9].GetSimpleName())
	assert.Equal(t, "sonar", ordered[10].GetSimpleName())
	assert.Equal(t, "usermgt", ordered[11].GetSimpleName())
}

func TestSortDogusByInvertedDependency(t *testing.T) {
	t.Run("should return empty slice", func(t *testing.T) {
		actual := SortDogusByInvertedDependency([]*Dogu{})

		require.Empty(t, actual)
	})

	t.Run("should dogus sorted by least important dogu first", func(t *testing.T) {
		dogus := []*Dogu{}
		dogus, _, _ = ReadDogusFromFile("../resources/test/dogu-sort-002.json")

		actual := SortDogusByInvertedDependency(dogus)

		assert.Equal(t, "icoordinator", actual[0].GetSimpleName())
		assert.Equal(t, "nginx", actual[1].GetSimpleName())
		assert.Equal(t, "registrator", actual[2].GetSimpleName())
	})

	t.Run("should return dogus sorted by by least important dogus first including optional dogus", func(t *testing.T) {
		// Dependency tree:
		// d ----> a ----> b
		//          \_    ^
		//            \   |
		//             -> c
		a := &Dogu{
			Name:                 "a",
			Dependencies:         []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		b := &Dogu{Name: "b"}

		c := &Dogu{
			Name:                 "c",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		}

		d := &Dogu{
			Name:                 "d",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
		}

		dogus := SortDogusByInvertedDependency([]*Dogu{a, b, c, d})
		assert.Equal(t, "d", dogus[0].Name)
		assert.Equal(t, "a", dogus[1].Name)
		assert.Equal(t, "c", dogus[2].Name)
		assert.Equal(t, "b", dogus[3].Name)
	})
}

func TestTransformsDependencyListToDoguList(t *testing.T) {
	dogus := []*Dogu{}
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-003.json")
	require.Nil(t, err)

	sorter := sortByDependency{dogus}
	result := sorter.dependenciesToDogus([]Dependency{
		{
			Type: "dogu",
			Name: "dogua",
		},
		{
			Type: "dogu",
			Name: "dogub",
		},
		{
			Type: "dogu",
			Name: "doguc",
		},
		{
			Type: "dogu",
			Name: "dogud",
		},
		{
			Type: "dogu",
			Name: "dogue",
		},
	})

	assert.Len(t, result, 5)
	assert.Equal(t, result[0].Name, "testing/dogua")
	assert.Equal(t, result[1].Name, "testing/dogub")
	assert.Equal(t, result[2].Name, "testing/doguc")
	assert.Equal(t, result[3].Name, "testing/dogud")
	assert.Equal(t, result[4].Name, "testing/dogue")
}

func TestGetDependenciesRecursive(t *testing.T) {
	dogus := []*Dogu{}
	dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-003.json")
	require.Nil(t, err)
	sorter := sortByDependency{dogus}

	tests := []struct {
		name                 string
		dogu                 *Dogu
		expectedDependencies []Dependency
	}{
		{
			name: "dependencies for dogua",
			dogu: getDoguBySimpleName(t, dogus, "dogua"),
			expectedDependencies: []Dependency{
				createDependency(t, "nginx"),
				createDependency(t, "dogub"),
				createDependency(t, "doguc"),
				createDependency(t, "dogud"),
				createDependency(t, "dogue"),
			},
		},
		{
			name: "dependencies for dogub",
			dogu: getDoguBySimpleName(t, dogus, "dogub"),
			expectedDependencies: []Dependency{
				createDependency(t, "nginx"),
				createDependency(t, "doguc"),
				createDependency(t, "dogud"),
				createDependency(t, "dogue"),
			},
		},
		{
			name: "dependencies for doguc",
			dogu: getDoguBySimpleName(t, dogus, "doguc"),
			expectedDependencies: []Dependency{
				createDependency(t, "nginx"),
				createDependency(t, "dogud"),
				createDependency(t, "dogue"),
			},
		},
		{
			name: "dependencies for dogud",
			dogu: getDoguBySimpleName(t, dogus, "dogud"),
			expectedDependencies: []Dependency{
				createDependency(t, "nginx"),
			},
		},
		{
			name: "dependencies for dogue",
			dogu: getDoguBySimpleName(t, dogus, "dogue"),
			expectedDependencies: []Dependency{
				createDependency(t, "nginx"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dependencies := sorter.getAllDoguDependenciesRecursive(test.dogu)
			assert.Len(t, dependencies, len(test.expectedDependencies))
			for _, expectedDependency := range test.expectedDependencies {
				assert.True(t, contains(dependencies, expectedDependency.Name))
			}
		})
	}
}

func createDependency(t *testing.T, name string) Dependency {
	t.Helper()
	return Dependency{
		Type: DependencyTypeDogu,
		Name: name,
	}
}

func getDoguBySimpleName(t *testing.T, dogus []*Dogu, doguName string) *Dogu {
	t.Helper()
	for _, dogu := range dogus {
		if dogu.GetSimpleName() == doguName {
			return dogu
		}
	}

	return nil
}
