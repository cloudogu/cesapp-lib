package core

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSortDogusByDependencyWithError(t *testing.T) {
	t.Run("should return ordered dogu slice with many dogus", func(t *testing.T) {
		dogus := []*Dogu{}
		dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-004.json")
		require.NoError(t, err)
		var doguNames []string
		for _, dogu := range dogus {
			doguNames = append(doguNames, dogu.GetSimpleName())
		}
		ordered, err := SortDogusByDependencyWithError(dogus)
		require.NoError(t, err)

		var installedDogus []string
		assert.Equal(t, len(dogus), len(ordered))
		for _, orderedDogu := range ordered {
			for _, dependency := range orderedDogu.GetDependenciesOfType(DependencyTypeDogu) {
				assert.Contains(t, installedDogus, dependency.Name,
					"%s installed before dependency: %s", orderedDogu.GetSimpleName(), dependency.Name)
			}

			for _, optionalDependency := range orderedDogu.GetOptionalDependenciesOfType(DependencyTypeDogu) {
				if stringSliceContains(doguNames, optionalDependency.Name) {
					assert.Contains(t, installedDogus, optionalDependency.Name,
						"%s installed before dependency: %s", orderedDogu.GetSimpleName(), optionalDependency.Name)
				}
			}
			installedDogus = append(installedDogus, orderedDogu.GetSimpleName())
		}
	})

	t.Run("should return ordered dogu slice with small list of dogus", func(t *testing.T) {
		dogus := []*Dogu{}
		dogus, _, err := ReadDogusFromFile("../resources/test/dogu-sort-002.json")
		assert.Nil(t, err)
		ordered, err := SortDogusByDependencyWithError(dogus)
		require.NoError(t, err)

		assert.Equal(t, "registrator", ordered[0].GetSimpleName())
		assert.Equal(t, "nginx", ordered[1].GetSimpleName())
		assert.Equal(t, "icoordinator", ordered[2].GetSimpleName())
	})

	t.Run("should return ordered dogu slice with transitive dogu dependencies", func(t *testing.T) {
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

		ordered, err := SortDogusByDependencyWithError(dogus)
		require.NoError(t, err)

		var orderedNames []string
		for _, orderedDogu := range ordered {
			orderedNames = append(orderedNames, orderedDogu.GetSimpleName())
		}

		// Dogu E and Dogu D can both be installed first and the sort algorithm is not deterministic.
		assert.Contains(t, orderedNames[0:2], "dogue")
		assert.Contains(t, orderedNames[0:2], "dogud")
		assert.Equal(t, "doguc", orderedNames[2])
		assert.Equal(t, "dogub", orderedNames[3])
		assert.Equal(t, "dogua", orderedNames[4])
	})

	t.Run("should return ordered dogu slice with optional dogu dependencies", func(t *testing.T) {
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

		dogus, err := SortDogusByDependencyWithError([]*Dogu{a, b, c, d})
		require.NoError(t, err)
		assert.Equal(t, "b", dogus[0].Name)
		assert.Equal(t, "c", dogus[1].Name)
		assert.Equal(t, "a", dogus[2].Name)
		assert.Equal(t, "d", dogus[3].Name)
	})

	t.Run("should return error when cyclic dependencies", func(t *testing.T) {
		// Dependency tree:
		// a <--- b
		//  \_    ^
		//    \   |
		//     -> c
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		b := &Dogu{
			Name:         "b",
			Dependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
		}

		c := &Dogu{
			Name:                 "c",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		}

		dogus, err := SortDogusByDependencyWithError([]*Dogu{a, b, c})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sort by dependency failed")
		assert.ErrorContains(t, err, "error in sorting dogus by dependency")
		assert.Nil(t, dogus)
	})

	t.Run("should return dogu if only irrelevant optional dependencies are set", func(t *testing.T) {
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		dogus, err := SortDogusByDependencyWithError([]*Dogu{a})
		assert.NoError(t, err)
		assert.Len(t, dogus, 1)
		assert.Equal(t, "a", dogus[0].Name)
	})

}

func stringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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

func TestSortDogusByInvertedDependencyWithError(t *testing.T) {
	t.Run("should return empty slice", func(t *testing.T) {
		actual, err := SortDogusByInvertedDependencyWithError([]*Dogu{})
		require.NoError(t, err)

		require.Empty(t, actual)
	})

	t.Run("should dogus sorted by least important dogu first", func(t *testing.T) {
		dogus := []*Dogu{}
		dogus, _, _ = ReadDogusFromFile("../resources/test/dogu-sort-002.json")

		actual, err := SortDogusByInvertedDependencyWithError(dogus)
		require.NoError(t, err)

		assert.Equal(t, "icoordinator", actual[0].GetSimpleName())
		assert.Equal(t, "nginx", actual[1].GetSimpleName())
		assert.Equal(t, "registrator", actual[2].GetSimpleName())
	})

	t.Run("should return dogus sorted by least important dogus first including optional dogus", func(t *testing.T) {
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

		dogus, err := SortDogusByInvertedDependencyWithError([]*Dogu{a, b, c, d})
		require.NoError(t, err)

		assert.Equal(t, "d", dogus[0].Name)
		assert.Equal(t, "a", dogus[1].Name)
		assert.Equal(t, "c", dogus[2].Name)
		assert.Equal(t, "b", dogus[3].Name)
	})

	t.Run("should return error when cyclic dependencies", func(t *testing.T) {
		// Dependency tree:
		// a <--- b
		//  \_    ^
		//    \   |
		//     -> c
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		b := &Dogu{
			Name:         "b",
			Dependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
		}

		c := &Dogu{
			Name:                 "c",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		}

		dogus, err := SortDogusByInvertedDependencyWithError([]*Dogu{a, b, c})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sort by dependency failed")
		assert.ErrorContains(t, err, "error in sorting dogus by inverted dependency")
		assert.Nil(t, dogus)
	})
}

func TestSortDogusByDependency(t *testing.T) {
	t.Run("should return nil slice and no error when cyclic dependencies", func(t *testing.T) {
		// Dependency tree:
		// a <--- b
		//  \_    ^
		//    \   |
		//     -> c
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		b := &Dogu{
			Name:         "b",
			Dependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
		}

		c := &Dogu{
			Name:                 "c",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		}

		//goland:noinspection GoDeprecation
		dogus := SortDogusByDependency([]*Dogu{a, b, c})
		assert.Nil(t, dogus)
	})

}

func TestSortDogusByInvertedDependency(t *testing.T) {
	t.Run("should return nil slice and no error when cyclic dependencies", func(t *testing.T) {
		// Dependency tree:
		// a <--- b
		//  \_    ^
		//    \   |
		//     -> c
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		b := &Dogu{
			Name:         "b",
			Dependencies: []Dependency{{Type: DependencyTypeDogu, Name: "a"}},
		}

		c := &Dogu{
			Name:                 "c",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "b"}},
		}

		//goland:noinspection GoDeprecation
		dogus := SortDogusByInvertedDependency([]*Dogu{a, b, c})
		assert.Nil(t, dogus)
	})

	t.Run("should return dogu if only irrelevant optional dependencies are set", func(t *testing.T) {
		a := &Dogu{
			Name:                 "a",
			OptionalDependencies: []Dependency{{Type: DependencyTypeDogu, Name: "c"}},
		}

		dogus, err := SortDogusByInvertedDependencyWithError([]*Dogu{a})
		assert.NoError(t, err)
		assert.Len(t, dogus, 1)
		assert.Equal(t, "a", dogus[0].Name)
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
