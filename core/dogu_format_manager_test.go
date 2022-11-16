package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ErrorMessageCannotUnmarshalDependenciesV2 = "json: cannot unmarshal string into Go struct field Dogu.Dependencies of type core.Dependency"
	ErrorMessageCannotUnmarshalDependenciesV1 = "json: cannot unmarshal object into Go struct field DoguV1.Dependencies of type string"
)

func Test_DoguFormatHandler_GetFormatProvider(t *testing.T) {
	// when
	providers := formatHandlerInstance.GetFormatProviders()

	// then
	assert.Equal(t, 2, len(providers))
}

func Test_DoguFormatHandler_RegisterFormatProvider(t *testing.T) {
	// given
	handler := DoguFormatHandler{}
	formatProvider := DoguJsonV1FormatProvider{}

	// when
	handler.RegisterFormatProvider(&formatProvider)

	// then
	assert.Equal(t, 1, len(handler.GetFormatProviders()))
	assert.Equal(t, &formatProvider, handler.GetFormatProviders()[0])
}

func TestDoguVolumeClientExpansion(t *testing.T) {
	t.Run("test read dogus from file", func(t *testing.T) {
		t.Run("Read v1 dogu content from file", func(t *testing.T) {
			// when
			dogu, _, err := ReadDoguFromFile("../resources/test/dogu-dependencies.json")

			// then
			assert.NoError(t, err)
			assert.Equal(t, "scm", dogu.Name)
			assert.Equal(t, "SCM-Manager", dogu.DisplayName)
			assert.Equal(t, "1.46", dogu.Version)
			assert.Equal(t, []VolumeClient(nil), dogu.Volumes[0].Clients)
		})
		t.Run("Read v2 dogu content from file", func(t *testing.T) {
			// given
			expectedVolumes := []Volume{{
				Name:        "data",
				Path:        "/var/lib/scm",
				Owner:       "",
				Group:       "",
				NeedsBackup: true,
				Clients: []VolumeClient{
					{Name: "myClient", Params: map[string]interface{}{"MySecret": "supersecret", "Type": "myType"}},
					{Name: "mySecondClient", Params: map[string]interface{}{"Algorithm": "myAlg", "Style": "superstyle"}},
				},
			}}

			// when
			dogu, _, err := ReadDoguFromFile("../resources/test/dogu-volume-client-expansion.json")

			// then
			assert.NoError(t, err)
			assert.Equal(t, "scm", dogu.Name)
			assert.Equal(t, "SCM-Manager", dogu.DisplayName)
			assert.Equal(t, "1.46", dogu.Version)
			assert.Equal(t, expectedVolumes, dogu.Volumes)
		})
	})
}

func TestReadDoguFromFile(t *testing.T) {
	t.Run("Fail with invalid input file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDoguFromFile("../resources/test/invalid.json")

		// then
		assert.Error(t, err)
		assert.Nil(t, dogus)
	})
	t.Run("Read v2 dogu content from file", func(t *testing.T) {
		// given
		expectedDependency1 := Dependency{Type: DependencyTypeDogu, Name: "cas", Version: ">=4.1.1-2"}
		expectedDependency2 := Dependency{Type: DependencyTypeDogu, Name: "ldap"}
		expectedDependency3 := Dependency{Type: DependencyTypePackage, Name: "backup-watcher", Version: "<=1.0.1"}
		expectedDependency4 := Dependency{Type: DependencyTypePackage, Name: "etcd", Version: "1.x.x-x"}
		expectedDependency5 := Dependency{Type: DependencyTypeClient, Name: "ces-setup", Version: ">=2.0.1"}
		expectedDependency6 := Dependency{Type: DependencyTypeClient, Name: "cesapp", Version: ">=2.0.1"}
		expectedDependencies := []Dependency{expectedDependency1, expectedDependency2, expectedDependency3, expectedDependency4, expectedDependency5, expectedDependency6}

		expectedServiceAccount := ServiceAccount{
			Type:        "k8s-dogu-operator",
			Kind:        "k8s",
			AccountName: "myTestAccount",
		}

		// when
		dogu, _, err := ReadDoguFromFile("../resources/test/dogu-dependencies_v2.json")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "scm", dogu.Name)
		assert.Equal(t, "SCM-Manager", dogu.DisplayName)
		assert.Equal(t, "1.46", dogu.Version)
		assert.Equal(t, expectedDependencies, dogu.Dependencies)
		assert.Equal(t, expectedServiceAccount, dogu.ServiceAccounts[0])
	})
	t.Run("Read v1 dogu content from file and upgrade dependencies accordingly", func(t *testing.T) {
		// given
		handlerOld := formatHandlerInstance
		defer func() { formatHandlerInstance = handlerOld }()

		handler := DoguFormatHandler{}
		formatHandlerInstance = &handler
		formatProvider := DoguJsonV1FormatProvider{}
		handler.RegisterFormatProvider(&formatProvider)
		expectedDependency1 := Dependency{Type: DependencyTypeDogu, Name: "cas"}
		expectedDependency2 := Dependency{Type: DependencyTypeDogu, Name: "ldap"}
		expectedDependency3 := Dependency{Type: DependencyTypeDogu, Name: "postfix"}
		expectedDependencies := []Dependency{expectedDependency1, expectedDependency2, expectedDependency3}

		// when
		dogu, _, err := ReadDoguFromFile("../resources/test/dogu-dependencies.json")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "scm", dogu.Name)
		assert.Equal(t, "SCM-Manager", dogu.DisplayName)
		assert.Equal(t, "1.46", dogu.Version)
		assert.Equal(t, expectedDependencies, dogu.Dependencies)
	})
	t.Run("read multiple dogus from files", func(t *testing.T) {
		tests := []struct {
			name                string
			filePath            string
			expectError         bool
			expectMultipleDogus bool
			expectedApiVersion  DoguApiVersion
		}{
			{
				name:                "unknown version with invalid dogu.json",
				filePath:            "../resources/test/invalid.json",
				expectError:         true,
				expectMultipleDogus: false,
				expectedApiVersion:  DoguApiVersionUnknown,
			},
			{
				name:                "v1 version with v1 dogu.json",
				filePath:            "../resources/test/scm-manager.json",
				expectError:         false,
				expectMultipleDogus: false,
				expectedApiVersion:  DoguApiV1,
			},
			{
				name:                "v2 version with v2 dogu.json",
				filePath:            "../resources/test/scm-manager_v2.json",
				expectError:         false,
				expectMultipleDogus: false,
				expectedApiVersion:  DoguApiV2,
			},
			{
				name:                "v1 version with v1 dogu.json (multiple dogus)",
				filePath:            "../resources/test/multipleDogus.json",
				expectError:         false,
				expectMultipleDogus: true,
				expectedApiVersion:  DoguApiV1,
			},
			{
				name:                "v2 version with v2 dogu.json (multiple dogus)",
				filePath:            "../resources/test/multipleDogus_v2.json",
				expectError:         false,
				expectMultipleDogus: true,
				expectedApiVersion:  DoguApiV2,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var version DoguApiVersion
				var err error
				if test.expectMultipleDogus {
					_, version, err = ReadDogusFromFile(test.filePath)
				} else {
					_, version, err = ReadDoguFromFile(test.filePath)
				}
				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, test.expectedApiVersion, version)
			})
		}
	})
}

func TestReadDogusFromFile(t *testing.T) {
	t.Run("Fail with invalid input file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromFile("../resources/test/invalid.json")

		// then
		assert.Error(t, err)
		assert.Nil(t, dogus)
	})

	t.Run("Read multiple dogus of v1 json content file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromFile("../resources/test/multipleDogus.json")

		// then
		assert.NoError(t, err)
		assert.Equal(t, 4, len(dogus))
		assert.Equal(t, "ldap", dogus[0].Name)
		assert.Equal(t, "jenkins", dogus[1].Name)
		assert.Equal(t, "nexus", dogus[2].Name)
		assert.Equal(t, "nginx", dogus[3].Name)
	})

	t.Run("Read multiple dogus of v2 json content file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromFile("../resources/test/multipleDogus_v2.json")

		// then
		assert.NoError(t, err)
		assert.Equal(t, 4, len(dogus))
		assert.Equal(t, "ldap", dogus[0].Name)
		assert.Equal(t, "jenkins", dogus[1].Name)
		assert.Equal(t, "nexus", dogus[2].Name)
		assert.Equal(t, "nginx", dogus[3].Name)
	})
}

func TestReadDoguFromString(t *testing.T) {
	// given
	contentInvalid, _ := GetContentOfFile("../resources/test/invalid.json")
	contentV1, _ := GetContentOfFile("../resources/test/dogu-dependencies.json")
	contentV2, _ := GetContentOfFile("../resources/test/dogu-dependencies_v2.json")

	t.Run("Fail with invalid input file", func(t *testing.T) {
		// when
		dogu, _, err := ReadDoguFromString(contentInvalid)

		// then
		assert.Error(t, err)
		assert.Nil(t, dogu)
	})

	t.Run("Read a single dogu of v1 json content file", func(t *testing.T) {
		// when
		dogu, _, err := ReadDoguFromString(contentV1)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, dogu)
	})

	t.Run("Read a single dogu of v2 json content file", func(t *testing.T) {
		// when
		dogu, _, err := ReadDoguFromString(contentV2)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, dogu)
	})
}

func TestReadDogusFromString(t *testing.T) {
	// given
	contentInvalid, _ := GetContentOfFile("../resources/test/invalid.json")
	contentV1, _ := GetContentOfFile("../resources/test/multipleDogus.json")
	contentV2, _ := GetContentOfFile("../resources/test/multipleDogus_v2.json")

	t.Run("Fail with invalid input file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromString(contentInvalid)

		// then
		assert.Error(t, err)
		assert.Nil(t, dogus)
	})

	t.Run("Read multiple dogus of v1 json content file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromString(contentV1)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, dogus)
		assert.Equal(t, 4, len(dogus))
		assert.Equal(t, "ldap", dogus[0].Name)
		assert.Equal(t, "jenkins", dogus[1].Name)
		assert.Equal(t, "nexus", dogus[2].Name)
		assert.Equal(t, "nginx", dogus[3].Name)
	})

	t.Run("Read multiple dogus of v2 json content file", func(t *testing.T) {
		// when
		dogus, _, err := ReadDogusFromString(contentV2)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, dogus)
		assert.Equal(t, 4, len(dogus))
		assert.Equal(t, "ldap", dogus[0].Name)
		assert.Equal(t, "jenkins", dogus[1].Name)
		assert.Equal(t, "nexus", dogus[2].Name)
		assert.Equal(t, "nginx", dogus[3].Name)
	})
}

func TestWriteDoguToFile(t *testing.T) {
	t.Run("Ensure simple content is begin written", func(t *testing.T) {
		// given
		file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
		path := file.Name()
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI"}
		err := WriteDoguToFile(path, &dogu)
		require.NoError(t, err)

		// then
		d, _, err := ReadDoguFromFile(path)
		require.NoError(t, err)
		assert.Equal(t, "jenkins", d.Name, "Name should be equal")
		assert.Equal(t, "1.625.2", d.Version, "Version should be equal")
		assert.Equal(t, "Jenkins CI", d.DisplayName, "DisplayName should be equal")
	})
	t.Run("Ensure only v1 format is begin written to the file. This transforms extended dependencies into their simple form", func(t *testing.T) {
		// given
		file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
		path := file.Name()
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		dependency1 := Dependency{Type: DependencyTypeDogu, Name: "testDogu", Version: "2.1.2"}
		dependency2 := Dependency{Type: DependencyTypeClient, Name: "testClient", Version: "2.1.2"}
		dependency3 := Dependency{Type: DependencyTypeClient, Name: "testPackage", Version: "2.1.2"}
		dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Dependencies: []Dependency{dependency1, dependency2, dependency3}}
		expectedDependency := []Dependency{{Type: DependencyTypeDogu, Name: "testDogu", Version: "2.1.2"}, {Type: "client", Name: "testClient", Version: "2.1.2"}, {Type: "client", Name: "testPackage", Version: "2.1.2"}}
		err := WriteDoguToFile(path, &dogu)
		require.NoError(t, err)

		// then
		d, _, err := ReadDoguFromFile(path)
		require.NoError(t, err)
		assert.Equal(t, "jenkins", d.Name, "Name should be equal")
		assert.Equal(t, "1.625.2", d.Version, "Version should be equal")
		assert.Equal(t, "Jenkins CI", d.DisplayName, "DisplayName should be equal")
		assert.Equal(t, expectedDependency, d.Dependencies)
	})
}

func TestWriteDoguToFileWithFormat(t *testing.T) {
	t.Run("Explicitly write the dogu as V2 format to preserve the advanced dependencies", func(t *testing.T) {
		// given
		file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
		path := file.Name()
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		dependency1 := Dependency{Type: DependencyTypeDogu, Name: "testDogu", Version: "2.1.2"}
		dependency2 := Dependency{Type: DependencyTypeClient, Name: "testClient", Version: "2.1.2"}
		dependency3 := Dependency{Type: DependencyTypeClient, Name: "testPackage", Version: "2.1.2"}
		dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Dependencies: []Dependency{dependency1, dependency2, dependency3}}
		err := WriteDoguToFileWithFormat(path, &dogu, &DoguJsonV2FormatProvider{})
		require.NoError(t, err)

		// then
		d, _, err := ReadDoguFromFile(path)
		require.NoError(t, err)
		assert.Equal(t, "jenkins", d.Name, "Name should be equal")
		assert.Equal(t, "1.625.2", d.Version, "Version should be equal")
		assert.Equal(t, "Jenkins CI", d.DisplayName, "DisplayName should be equal")
		assert.Equal(t, dogu.Dependencies, d.Dependencies)
	})
	t.Run("Write extended volume definitons with the v2 format", func(t *testing.T) {
		// given
		file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
		path := file.Name()
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		expectedVolumes := []Volume{{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
			Clients: []VolumeClient{
				{Name: "myClient", Params: map[string]interface{}{"MySecret": "supersecret", "Type": "myType"}},
				{Name: "mySecondClient", Params: map[string]interface{}{"Algorithm": "myAlg", "Style": "superstyle"}},
			},
		}}

		dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Volumes: expectedVolumes}
		err := WriteDoguToFileWithFormat(path, &dogu, &DoguJsonV2FormatProvider{})
		require.NoError(t, err)

		// then
		d, _, err := ReadDoguFromFile(path)
		require.NoError(t, err)
		assert.Equal(t, "jenkins", d.Name, "Name should be equal")
		assert.Equal(t, "1.625.2", d.Version, "Version should be equal")
		assert.Equal(t, "Jenkins CI", d.DisplayName, "DisplayName should be equal")
		assert.Equal(t, expectedVolumes, d.Volumes)
	})
	t.Run("Write volume definitons without volume clients with the v2 format", func(t *testing.T) {
		// given
		file, _ := ioutil.TempFile(os.TempDir(), "cesapp-")
		path := file.Name()
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		expectedVolumes := []Volume{{
			Name:        "data",
			Path:        "/var/lib/scm",
			Owner:       "",
			Group:       "",
			NeedsBackup: true,
		}}

		dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Volumes: expectedVolumes}
		err := WriteDoguToFileWithFormat(path, &dogu, &DoguJsonV2FormatProvider{})
		require.NoError(t, err)

		// then
		d, _, err := ReadDoguFromFile(path)
		require.NoError(t, err)
		assert.Equal(t, "jenkins", d.Name, "Name should be equal")
		assert.Equal(t, "1.625.2", d.Version, "Version should be equal")
		assert.Equal(t, "Jenkins CI", d.DisplayName, "DisplayName should be equal")
		assert.Equal(t, expectedVolumes, d.Volumes)
	})
}

type formatTest struct {
	name     string
	provider DoguFormatProvider
	wantErr  bool
	errorMsg string
	expected string
}

func assertTestError(t *testing.T, err error, test formatTest) {
	require.Error(t, err)
	assert.Contains(t, err.Error(), test.errorMsg)
}

func Test_DoguFormatProvider_ReadDoguFromString(t *testing.T) {
	contentV1, _ := GetContentOfFile("../resources/test/scm-manager.json")
	contentV2, _ := GetContentOfFile("../resources/test/scm-manager_v2.json")
	contentVolumes, _ := GetContentOfFile("../resources/test/dogu-volumes.json")
	contentDogu, _ := GetContentOfFile("../resources/test/conf-dogu/dogu.json")
	expectedDependencies := []Dependency{
		{Type: DependencyTypeDogu, Name: "cas"},
	}

	t.Run("Reading from string with v1 content", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: true, errorMsg: ErrorMessageCannotUnmarshalDependenciesV2},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogu, err := test.provider.ReadDoguFromString(contentV1)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					//then
					require.NoError(t, err)
					assert.Equal(t, "scm", dogu.Name, "Name should be equal")
					assert.Equal(t, "1.46", dogu.Version, "Version should be equal")
					assert.Equal(t, "https://www.scm-manager.org", dogu.URL, "Version should be equal")
					assert.Equal(t, expectedDependencies, dogu.Dependencies, "Dependency should be equal")
					assert.Equal(t, []Volume{{"data", "/var/lib/scm", "", "", true, nil}}, dogu.Volumes, "Volumes should be equal")
					assert.Equal(t, HealthCheck{Type: "tcp", Port: 8080}, dogu.HealthCheck, "HealthCheck should be equal")
				}
			})
		}
	})
	t.Run("Reading from string with v2 content", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: true, errorMsg: ErrorMessageCannotUnmarshalDependenciesV1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogu, err := test.provider.ReadDoguFromString(contentV2)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					//then
					require.NoError(t, err)
					assert.Equal(t, "scm", dogu.Name, "Name should be equal")
					assert.Equal(t, "1.46", dogu.Version, "Version should be equal")
					assert.Equal(t, "https://www.scm-manager.org", dogu.URL, "Version should be equal")
					assert.Equal(t, expectedDependencies, dogu.Dependencies, "Dependency should be equal")
					assert.Equal(t, []Volume{{"data", "/var/lib/scm", "", "", true, nil}}, dogu.Volumes, "Volumes should be equal")
					assert.Equal(t, HealthCheck{Type: "tcp", Port: 8080}, dogu.HealthCheck, "HealthCheck should be equal")
				}
			})
		}
	})
	t.Run("Test 'NeedsBackup' property by reading from file", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogu, err := test.provider.ReadDoguFromString(contentVolumes)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					//then
					require.NoError(t, err)
					assert.Equal(t, Volume{Name: "withBackup", NeedsBackup: true}, dogu.Volumes[0], "Volumes should be equal")
					assert.Equal(t, Volume{Name: "withoutBackup", NeedsBackup: false}, dogu.Volumes[1], "Volumes should be equal")
					assert.Equal(t, Volume{Name: "withDefault", NeedsBackup: true}, dogu.Volumes[2], "Volumes should be equal")
				}
			})
		}
	})
	t.Run("Test 'Configuration' entries by reading from file", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogu, err := test.provider.ReadDoguFromString(contentDogu)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					// then
					require.NoError(t, err)
					assert.Equal(t, "title", dogu.Configuration[0].Name)
					assert.Equal(t, "Title of index page", dogu.Configuration[0].Description)
					assert.False(t, dogu.Configuration[0].Optional)
					assert.True(t, dogu.Configuration[1].Optional)
				}
			})
		}
	})
}

func Test_DoguFormatProvider_ReadDogusFromString(t *testing.T) {
	contentMultipleV1, _ := GetContentOfFile("../resources/test/multipleDogus.json")
	contentMultipleV2, _ := GetContentOfFile("../resources/test/multipleDogus_v2.json")

	t.Run("Reading multiples dogus from string with v1 content", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: true, errorMsg: ErrorMessageCannotUnmarshalDependenciesV2},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogus, err := test.provider.ReadDogusFromString(contentMultipleV1)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					//then
					require.NoError(t, err)
					assert.Equal(t, 4, len(dogus))
					assert.Equal(t, "ldap", dogus[0].Name)
					assert.Nil(t, dogus[0].Dependencies)
					assert.Equal(t, 0, len(dogus[0].Dependencies))

					assert.Equal(t, "jenkins", dogus[1].Name)
					assert.NotNil(t, dogus[1].Dependencies)
					assert.Equal(t, 3, len(dogus[1].Dependencies))

					assert.Equal(t, "nexus", dogus[2].Name)
					assert.NotNil(t, dogus[2].Dependencies)
					assert.Equal(t, 3, len(dogus[2].Dependencies))

					assert.Equal(t, "nginx", dogus[3].Name)
					assert.NotNil(t, dogus[3].Dependencies)
					assert.Equal(t, 1, len(dogus[3].Dependencies))
				}
			})
		}
	})
	t.Run("Reading from string with v2 content", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: true, errorMsg: ErrorMessageCannotUnmarshalDependenciesV1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				dogus, err := test.provider.ReadDogusFromString(contentMultipleV2)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					//then
					require.NoError(t, err)
					assert.Equal(t, 4, len(dogus))
					assert.Equal(t, "ldap", dogus[0].Name)
					assert.Nil(t, dogus[0].Dependencies)
					assert.Equal(t, 0, len(dogus[0].Dependencies))

					assert.Equal(t, "jenkins", dogus[1].Name)
					assert.NotNil(t, dogus[1].Dependencies)
					assert.Equal(t, 3, len(dogus[1].Dependencies))

					assert.Equal(t, "nexus", dogus[2].Name)
					assert.NotNil(t, dogus[2].Dependencies)
					assert.Equal(t, 3, len(dogus[2].Dependencies))

					assert.Equal(t, "nginx", dogus[3].Name)
					assert.NotNil(t, dogus[3].Dependencies)
					assert.Equal(t, 1, len(dogus[3].Dependencies))
				}
			})
		}
	})
}

func Test_DoguFormatProvider_WriteDoguToString(t *testing.T) {
	expectedDependencies := []Dependency{
		{Type: DependencyTypeDogu, Name: "cas"},
	}
	dogu := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Dependencies: expectedDependencies}
	expectedRepresentationV1 := "{\"Name\":\"jenkins\",\"Version\":\"1.625.2\",\"DisplayName\":\"Jenkins CI\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[\"cas\"],\"OptionalDependencies\":null}"
	expectedRepresentationV2 := "{\"Name\":\"jenkins\",\"Version\":\"1.625.2\",\"DisplayName\":\"Jenkins CI\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[{\"type\":\"dogu\",\"name\":\"cas\",\"version\":\"\"}],\"OptionalDependencies\":null}"

	t.Run("Convert given dogu object into the string representation", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false, expected: expectedRepresentationV2},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false, expected: expectedRepresentationV1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				stringRepresentation, err := test.provider.WriteDoguToString(&dogu)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					// then
					require.NoError(t, err)
					assert.Equal(t, test.expected, stringRepresentation)
					assert.Contains(t, stringRepresentation, dogu.Name)
					assert.Contains(t, stringRepresentation, dogu.Version)
					assert.Contains(t, stringRepresentation, dogu.DisplayName)
					assert.Contains(t, stringRepresentation, dogu.Dependencies[0].Name)
				}
			})
		}
	})
}

func Test_DoguFormatProvider_WriteDogusToString(t *testing.T) {
	expectedDependencies := []Dependency{
		{Type: DependencyTypeDogu, Name: "cas"},
	}
	dogu1 := Dogu{Name: "jenkins", Version: "1.625.2", DisplayName: "Jenkins CI", Dependencies: expectedDependencies}
	dogu2 := Dogu{Name: "scm", Version: "2.625.2", DisplayName: "Scm Manager", Dependencies: expectedDependencies}
	dogus := []*Dogu{&dogu1, &dogu2}
	expectedRepresentationV1 := "[{\"Name\":\"jenkins\",\"Version\":\"1.625.2\",\"DisplayName\":\"Jenkins CI\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[\"cas\"],\"OptionalDependencies\":null},{\"Name\":\"scm\",\"Version\":\"2.625.2\",\"DisplayName\":\"Scm Manager\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[\"cas\"],\"OptionalDependencies\":null}]"
	expectedRepresentationV2 := "[{\"Name\":\"jenkins\",\"Version\":\"1.625.2\",\"DisplayName\":\"Jenkins CI\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[{\"type\":\"dogu\",\"name\":\"cas\",\"version\":\"\"}],\"OptionalDependencies\":null},{\"Name\":\"scm\",\"Version\":\"2.625.2\",\"DisplayName\":\"Scm Manager\",\"Description\":\"\",\"Category\":\"\",\"Tags\":null,\"Logo\":\"\",\"URL\":\"\",\"Image\":\"\",\"ExposedPorts\":null,\"ExposedCommands\":null,\"Volumes\":null,\"HealthCheck\":{\"Type\":\"\",\"State\":\"\",\"Port\":0,\"Path\":\"\",\"Parameters\":null},\"HealthChecks\":null,\"ServiceAccounts\":null,\"Privileged\":false,\"Configuration\":null,\"Properties\":null,\"EnvironmentVariables\":null,\"Dependencies\":[{\"type\":\"dogu\",\"name\":\"cas\",\"version\":\"\"}],\"OptionalDependencies\":null}]"

	t.Run("Convert given dogu object into the string representation", func(t *testing.T) {
		tests := []formatTest{
			{name: "V2", provider: &DoguJsonV2FormatProvider{}, wantErr: false, expected: expectedRepresentationV2},
			{name: "V1", provider: &DoguJsonV1FormatProvider{}, wantErr: false, expected: expectedRepresentationV1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// when
				stringRepresentation, err := test.provider.WriteDogusToString(dogus)
				if test.wantErr {
					assertTestError(t, err, test)
				} else {
					// then
					require.NoError(t, err)
					assert.Equal(t, test.expected, stringRepresentation)
					assert.Contains(t, stringRepresentation, dogus[0].Name)
					assert.Contains(t, stringRepresentation, dogus[1].Name)
					assert.Contains(t, stringRepresentation, dogus[0].Version)
					assert.Contains(t, stringRepresentation, dogus[1].Version)
					assert.Contains(t, stringRepresentation, dogus[0].DisplayName)
					assert.Contains(t, stringRepresentation, dogus[1].DisplayName)
					assert.Contains(t, stringRepresentation, dogus[0].Dependencies[0].Name)
					assert.Contains(t, stringRepresentation, dogus[1].Dependencies[0].Name)
				}
			})
		}
	})
}
