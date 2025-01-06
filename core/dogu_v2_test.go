package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetImageName(t *testing.T) {
	dogu := Dogu{Image: "trillian/mcmillian", Version: "1.0.0"}
	assert.Equal(t, "trillian/mcmillian:1.0.0", dogu.GetImageName())
	dogu = Dogu{Image: "trillian/mcmillian"}
	assert.Equal(t, "trillian/mcmillian", dogu.GetImageName())
}

func TestGetSimpleName(t *testing.T) {
	dogu := Dogu{Name: "trillian/mcmillian"}
	assert.Equal(t, "mcmillian", dogu.GetSimpleName())
	dogu = Dogu{Name: "trillian"}
	assert.Equal(t, "trillian", dogu.GetSimpleName())
}

func TestGetNameSpace(t *testing.T) {
	dogu := Dogu{Name: "trillian/mcmillian"}
	assert.Equal(t, "trillian", dogu.GetNamespace())
	dogu = Dogu{Name: "mcmillian"}
	assert.Equal(t, "mcmillian", dogu.GetNamespace())
}

func TestGetFullName(t *testing.T) {
	dogu := Dogu{Name: "trillian/mcmillian"}
	assert.Equal(t, "trillian/mcmillian", dogu.GetFullName())
}

func TestGetEnvironmentVariablesAsStringArray(t *testing.T) {
	sampleEnv := EnvironmentVariable{"HAPPY_MODE", "true"}
	dogu := Dogu{
		Name:                 "trillian/mcmillian",
		EnvironmentVariables: []EnvironmentVariable{sampleEnv},
	}
	assert.Equal(t, []string{"HAPPY_MODE=true"}, dogu.GetEnvironmentVariablesAsStringArray())
}

func TestGetVersion(t *testing.T) {
	dogu := Dogu{Name: "ces/hansolo", Version: "2.3-4"}
	version, err := dogu.GetVersion()
	assert.Nil(t, err)
	assert.Equal(t, 2, version.Major)
	assert.Equal(t, 3, version.Minor)
	assert.Equal(t, 0, version.Patch)
	assert.Equal(t, 0, version.Nano)
	assert.Equal(t, 4, version.Extra)
}

func TestIsEqualTo(t *testing.T) {
	dogu1 := &Dogu{Name: "ces/hansolo", Version: "2.3-4"}
	dogu2 := &Dogu{Name: "ces/hansolo", Version: "2.3-4"}

	actual, err := dogu1.IsEqualTo(dogu2)

	assert.Nil(t, err)
	assert.True(t, actual)
}

func TestIsNotEqualTo(t *testing.T) {
	dogu1 := &Dogu{Name: "ces/hansolo", Version: "2.3-3"}
	dogu2 := &Dogu{Name: "ces/hansolo", Version: "2.3-4"}

	actual, err := dogu1.IsEqualTo(dogu2)

	assert.Nil(t, err)
	assert.False(t, actual)
}

func TestGetVersionWithError(t *testing.T) {
	dogu := Dogu{Name: "ces/hansolo", Version: "2.3.a-4"}
	_, err := dogu.GetVersion()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "version")
	assert.Contains(t, err.Error(), "2.3.a-4")
	assert.Contains(t, err.Error(), "ces/hansolo")
}

func TestDoguIsNewerThanWithDifferentDogus(t *testing.T) {
	dogu := &Dogu{Name: "ces/hansolo", Version: "2"}
	other := &Dogu{Name: "ces/lea", Version: "3"}

	_, err := dogu.IsNewerThan(other)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "same name")
}

func TestDoguIsNewerThanWithDifferentNamespace(t *testing.T) {
	dogu := &Dogu{Name: "ces/hansolo", Version: "2"}
	other := &Dogu{Name: "sw/hansolo", Version: "3"}

	_, err := dogu.IsNewerThan(other)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "same name")
}

func TestDoguIsNewerThan(t *testing.T) {
	dogu := &Dogu{Name: "ces/hansolo", Version: "3"}
	other := &Dogu{Name: "ces/hansolo", Version: "2"}

	newer, _ := dogu.IsNewerThan(other)
	assert.True(t, newer)

	other.Version = "4"

	newer, _ = dogu.IsNewerThan(other)
	assert.False(t, newer)
}

func TestHasExposedCommand_shouldFindCommand(t *testing.T) {
	dogu := &Dogu{ExposedCommands: []ExposedCommand{{Name: "doNothing"}, {Name: "service-account-create"}}}
	assert.True(t, dogu.HasExposedCommand(ExposedCommandServiceAccountCreate))
}

func TestHasExposedCommand_shouldNotFindCommand(t *testing.T) {
	dogu := &Dogu{ExposedCommands: []ExposedCommand{{Name: "doNothing"}}}
	assert.False(t, dogu.HasExposedCommand(ExposedCommandServiceAccountCreate))
}

func TestGetExposedCommand_shouldReturnCommand(t *testing.T) {
	command1 := ExposedCommand{
		Name:        ExposedCommandPostUpgrade,
		Description: "a description",
	}
	command2 := ExposedCommand{
		Name:        ExposedCommandServiceAccountCreate,
		Description: "another description",
	}
	sut := Dogu{Name: "doguA", ExposedCommands: []ExposedCommand{command1, command2}}

	actualCommand := sut.GetExposedCommand("service-account-create")

	require.Equal(t, &command2, actualCommand)
}

func TestGetExposedCommand_shouldReturnNil(t *testing.T) {
	command1 := ExposedCommand{
		Name:        ExposedCommandPostUpgrade,
		Description: "a description",
	}
	command2 := ExposedCommand{
		Name:        ExposedCommandServiceAccountCreate,
		Description: "another description",
	}
	sut := Dogu{Name: "doguA", ExposedCommands: []ExposedCommand{command1, command2}}

	actualCommand := sut.GetExposedCommand("somethingTotallyDifferent")

	require.Nil(t, actualCommand)
}

func TestIsServiceAccountProviderNoExposedCommand(t *testing.T) {
	dogu := &Dogu{ExposedCommands: []ExposedCommand{{Name: "doNothing"}}}
	assert.False(t, dogu.HasExposedCommand(ExposedCommandServiceAccountCreate))
}

func TestGetRegistryName(t *testing.T) {
	t.Run("with dogu nexus", func(t *testing.T) {
		dogu := &Dogu{
			Name:  "official/nexus",
			Image: "registry.cloudogu.com/official/nexus",
		}
		require.Equal(t, "registry.cloudogu.com", dogu.GetRegistryServerURI())
	})
	t.Run("with dogu jenkins", func(t *testing.T) {
		dogu := &Dogu{
			Name:  "official/jenkins",
			Image: "registry.cloudogu.com/official/jenkins",
		}
		require.Equal(t, "registry.cloudogu.com", dogu.GetRegistryServerURI())
	})
	t.Run("with registry containing slashes", func(t *testing.T) {
		dogu := &Dogu{
			Name:  "official/jenkins",
			Image: "registry/cloudogu/com/official/jenkins",
		}
		require.Equal(t, "registry/cloudogu/com", dogu.GetRegistryServerURI())
	})
}

func Test_getSimpleDoguName(t *testing.T) {
	actual := GetSimpleDoguName("official/redmine")

	assert.Equal(t, "redmine", actual)
}

func Test_getNamespace(t *testing.T) {
	actual := GetNamespace("official/redmine")

	assert.Equal(t, "official", actual)
}

func TestContainsDogu(t *testing.T) {
	dogua := createTestDogu(t, "dogua", "1.0.0-1")
	doguav2 := createTestDogu(t, "dogua", "1.0.0-2")
	dogub := createTestDogu(t, "dogub", "1.0.0-1")
	doguc := createTestDogu(t, "doguc", "1.0.0-1")
	dogud := createTestDogu(t, "dogud", "1.0.0-1")
	dogue := createTestDogu(t, "dogue", "1.0.0-1")

	dogus := []*Dogu{dogua, doguav2, dogub, doguc, dogud}

	assert.True(t, ContainsDoguWithName(dogus, dogua.Name))
	assert.True(t, ContainsDoguWithName(dogus, dogub.Name))
	assert.True(t, ContainsDoguWithName(dogus, doguc.Name))
	assert.True(t, ContainsDoguWithName(dogus, dogud.Name))
	assert.False(t, ContainsDoguWithName(dogus, dogue.Name))
	assert.False(t, ContainsDoguWithName(dogus, "official/doguf"))
	assert.False(t, ContainsDoguWithName(dogus, "something else"))
}

func createTestDogu(t *testing.T, name string, version string) *Dogu {
	t.Helper()
	return &Dogu{Name: "official/" + name, Version: version}
}
