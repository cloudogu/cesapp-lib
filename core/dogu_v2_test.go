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

func TestTestDogu(t *testing.T) {
	var mydogu core.Dogu
	err := json.Unmarshal([]byte(testDogu(t)), &mydogu)

	require.NoError(t, err)
	assert.NotNil(t, mydogu)
	assert.Equal(t, "asdf", mydogu.Dependencies[0].Name)
}
func TestTestDogu2(t *testing.T) {
	var mydogu core.Dogu
	err := json.Unmarshal([]byte(testDogu(t)), &mydogu)

	require.NoError(t, err)
	assert.NotNil(t, mydogu)

	dogubytes, err := json.Marshal(&mydogu)
	require.NoError(t, err)
	assert.Equal(t, "{hallo}", string(dogubytes))
}
func TestTestDogu2Bson(t *testing.T) {
	var mydogu core.Dogu
	err := json.Unmarshal([]byte(testDogu(t)), &mydogu)

	require.NoError(t, err)
	assert.NotNil(t, mydogu)

	// save into mongodb here ;)
	dogubytes, err := bson.Marshal(&mydogu)
	require.NoError(t, err)

	err = bson.Unmarshal(dogubytes, &mydogu)
	require.NoError(t, err)
	assert.True(t, time.Time(mydogu.PublishedAt).IsZero())
}

func testDogu(t *testing.T) string {
	t.Helper()

	return `  {
    "_id": {
      "$oid": "666ab089e93bd78f3282c6d5"
    },
    "name": "testing/jenkins",
    "version": "2.440.2-999",
    "displayname": "Jenkins CI",
    "description": "Jenkins Continuous Integration Server",
    "category": "Development Apps",
    "tags": [
      "warp",
      "build",
      "ci",
      "cd"
    ],
    "logo": "https://cloudogu.com/images/dogus/jenkins.png",
    "url": "https://jenkins-ci.org",
    "image": "registry.cloudogu.com/testing/jenkins",
    "exposedports": null,
    "exposedcommands": [
      {
        "name": "upgrade-notification",
        "description": "",
        "command": "/upgrade-notification.sh"
      },
      {
        "name": "pre-upgrade",
        "description": "",
        "command": "/pre-upgrade.sh"
      }
    ],
    "volumes": [
      {
        "name": "data",
        "path": "/var/lib/jenkins",
        "owner": "1000",
        "group": "1000",
        "needsbackup": true,
        "clients": null
      },
      {
        "name": "custom.init.groovy.d",
        "path": "/var/lib/custom.init.groovy.d",
        "owner": "1000",
        "group": "1000",
        "needsbackup": true,
        "clients": null
      },
      {
        "name": "tmp",
        "path": "/tmp",
        "owner": "1000",
        "group": "1000",
        "needsbackup": false,
        "clients": null
      }
    ],
    "healthcheck": {
      "type": "",
      "state": "",
      "port": 0,
      "path": "",
      "parameters": null
    },
    "healthchecks": [
      {
        "type": "tcp",
        "state": "",
        "port": 8080,
        "path": "",
        "parameters": null
      },
      {
        "type": "state",
        "state": "",
        "port": 0,
        "path": "",
        "parameters": null
      }
    ],
    "serviceaccounts": null,
    "privileged": false,
    "configuration": [
      {
        "name": "additional.plugins",
        "description": "Comma separated list of plugin names to install on start",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "",
          "values": null
        }
      },
      {
        "name": "container_config/memory_limit",
        "description": "Limits the container's memory usage. Use a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte).",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/swap_limit",
        "description": "Limits the container's swap memory usage. Use zero or a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte). 0 will disable swapping.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/java_max_ram_percentage",
        "description": "Limits the heap stack size of the Jenkins process to the configured percentage of the available physical memory when the container has more than approx. 250 MB of memory available. Is only considered when a memory_limit is set. Use a valid float value with decimals between 0 and 100 (f. ex. 55.0 for 55%). Default value for Jenkins: 25%",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "25.0",
        "validation": {
          "type": "FLOAT_PERCENTAGE_HUNDRED",
          "values": null
        }
      },
      {
        "name": "container_config/java_min_ram_percentage",
        "description": "Limits the heap stack size of the Jenkins process to the configured percentage of the available physical memory when the container has less than approx. 250 MB of memory available. Is only considered when a memory_limit is set. Use a valid float value with decimals between 0 and 100 (f. ex. 55.0 for 55%). Default value for Jenkins: 50%",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "50.0",
        "validation": {
          "type": "FLOAT_PERCENTAGE_HUNDRED",
          "values": null
        }
      },
      {
        "name": "container_config/memory_limit",
        "description": "Limits the container's memory usage. Use a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte).",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/memory_request",
        "description": "Requests the container's minimal memory requirement. Use a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte).",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "2g",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/swap_limit",
        "description": "Limits the container's swap memory usage. Use zero or a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte). 0 will disable swapping.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/cpu_core_limit",
        "description": "Limits the container's CPU core usage. Use a positive floating value describing a fraction of 1 CPU core. When you define a value of '0.5', you are requesting half as much CPU time compared to if you asked for '1.0' CPU.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "",
          "values": null
        }
      },
      {
        "name": "container_config/cpu_core_request",
        "description": "Requests the container's minimal CPU core requirement. Use a positive floating value describing a fraction of 1 CPU core. When you define a value of '0.5', you are requesting half as much CPU time compared to if you asked for '1.0' CPU.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "1.0",
        "validation": {
          "type": "",
          "values": null
        }
      },
      {
        "name": "container_config/storage_limit",
        "description": "Limits the container's ephemeral storage usage. Use a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte).",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "container_config/storage_request",
        "description": "Requests the container's minimal ephemeral storage requirement. Use a positive integer value followed by one of these units [b,k,m,g] (byte, kibibyte, mebibyte, gibibyte).",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "",
        "validation": {
          "type": "BINARY_MEASUREMENT",
          "values": null
        }
      },
      {
        "name": "additional_java_args",
        "description": "Additional args that are passed to the jenkins process.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "UNSET",
        "validation": {
          "type": "",
          "values": null
        }
      },
      {
        "name": "logging/root",
        "description": "Set the root log level to one of ERROR, WARN, INFO, DEBUG.",
        "optional": true,
        "encrypted": false,
        "global": false,
        "default": "INFO",
        "validation": {
          "type": "ONE_OF",
          "values": [
            "WARN",
            "DEBUG",
            "INFO",
            "ERROR"
          ]
        }
      }
    ],
    "properties": null,
    "environmentvariables": null,
    "dependencies": [
      {
        "type": "dogu",
        "name": "cas",
        "version": ""
      },
      {
        "type": "dogu",
        "name": "nginx",
        "version": ""
      },
      {
        "type": "dogu",
        "name": "postfix",
        "version": ""
      }
    ],
    "optionaldependencies": null,
    "publishedAt": ""
  }
`
}
