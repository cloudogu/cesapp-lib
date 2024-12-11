package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
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

func TestUnmarshalProperties(t *testing.T) {
	dogu := &Dogu{}
	dogu, _, err := ReadDoguFromFile("../resources/test/unmarshalProperties.json")
	require.Nil(t, err)
	assert.Equal(t, "http://test.test", dogu.Properties["logoutUrl"])
	assert.Equal(t, "25", dogu.Properties["TestPort"])
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

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}

func routeStdoutToReplacement() (readerPipe, writerPipe *os.File) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return r, w
}

func captureOutput(fakeReaderPipe, fakeWriterPipe, originalStdout *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, fakeReaderPipe)
		outC <- buf.String()
	}()

	// back to normal state
	_ = fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}

func Test_getSimpleDoguName(t *testing.T) {
	actual := GetSimpleDoguName("official/redmine")

	assert.Equal(t, "redmine", actual)
}

func Test_getNamespace(t *testing.T) {
	actual := GetNamespace("official/redmine")

	assert.Equal(t, "official", actual)
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

func Test_sortDogus(t *testing.T) {
	unsortedDogus := []*Dogu{{Name: "Dogu1", Version: "11.22.33-1"}, {Name: "Dogu2", Version: "0.1.3-5"}, {Name: "Dogu3", Version: "0.5.3-3"}, {Name: "Dogu4", Version: "9.3.9"}}
	expectedDogus := []*Dogu{{Name: "Dogu1", Version: "11.22.33-1"}, {Name: "Dogu4", Version: "9.3.9"}, {Name: "Dogu3", Version: "0.5.3-3"}, {Name: "Dogu2", Version: "0.1.3-5"}}

	sort.Sort(ByDoguVersion(unsortedDogus))

	assert.Equal(t, unsortedDogus, expectedDogus)
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

func Test_validateSecurity(t *testing.T) {
	type args struct {
		dogu *Dogu
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"valid empty", args{&Dogu{}}, assert.NoError},
		{"valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: []Capability{AuditControl}}}}}, assert.NoError},
		{"valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Drop: []Capability{AuditControl}}}}}, assert.NoError},
		{"all possible values", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: allCapabilities, Drop: allCapabilities}}}}, assert.NoError},

		{"invalid valid add filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Add: []Capability{"err"}}}}}, assert.Error},
		{"invalid valid drop filled", args{&Dogu{Security: Security{Capabilities: Capabilities{Drop: []Capability{"err"}}}}}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.args.dogu.ValidateSecurity(), fmt.Sprintf("validateSecurity(%v)", tt.args.dogu))
		})
	}
}

func Test_validateSecurity_message(t *testing.T) {
	t.Run("should match for drop errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Drop: []Capability{"err"}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu official/dogu:1.2.3 contains an invalid security field: err is not a valid capability to be dropped")
	})
	t.Run("should match for add errors", func(t *testing.T) {
		// given
		dogu := &Dogu{Name: "official/dogu", Version: "1.2.3", Security: Security{Capabilities: Capabilities{Add: []Capability{"err"}}}}

		// when
		actual := dogu.ValidateSecurity()

		// then
		require.Error(t, actual)
		assert.ErrorContains(t, actual, "dogu official/dogu:1.2.3 contains an invalid security field: err is not a valid capability to be added")
	})
}

func TestDogu_EffectiveCapabilities(t *testing.T) {
	type fields struct {
		Security Security
	}
	tests := []struct {
		name   string
		fields fields
		want   []Capability
	}{
		{"drop all", fields{Security{Capabilities: Capabilities{Drop: []Capability{All}}}}, []Capability{}},
		{"add all", fields{Security{Capabilities: Capabilities{Add: []Capability{All}}}}, allCapabilities},
		{"drop all, add all", fields{Security{Capabilities: Capabilities{Drop: []Capability{All}, Add: []Capability{All}}}}, allCapabilities},
		{"default list", fields{Security{Capabilities: Capabilities{}}}, DefaultCapabilities},
		{"add 1 new and 1 existing caps to default list", fields{Security{Capabilities: Capabilities{Add: []Capability{Bpf, Chown}}}}, joinCapability(DefaultCapabilities, Bpf, Chown)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dogu{
				Security: tt.fields.Security,
			}
			assert.ElementsMatch(t, tt.want, d.EffectiveCapabilities(), "ListCapabilities()")
		})
	}
}

func joinCapability(capSlice []Capability, singleCaps ...Capability) []Capability {
	result := []Capability{}
	result = append(result, capSlice...)
	for _, singleCap := range singleCaps {
		if slices.Contains(capSlice, singleCap) {
			continue
		}

		result = append(result, singleCap)
	}
	return result
}
