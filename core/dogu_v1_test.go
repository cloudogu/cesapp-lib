package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DoguV1_createV2Copy(t *testing.T) {
	t.Run("Make deep conversion of old dogu schematic to new version", func(t *testing.T) {
		// given
		doguV1 := DoguV1{
			Name:        "test",
			Version:     "0.1",
      PublishedAt: time.Date(2024, 10, 16, 7, 33, 45, 456, time.UTC),
			DisplayName: "displayName",
			Description: "description",
			Category:    "category",
			Tags:        []string{"tag1", "tag2"},
			Logo:        "logo",
			URL:         "url",
			Image:       "image",
			ExposedPorts: []ExposedPort{{
				Type:      "tcp",
				Container: 1,
				Host:      2,
			}},
			ExposedCommands: []ExposedCommand{
				{
					Name:        "TestCommand",
					Description: "Does nothing, just for testing purpose",
					Command:     "./test.sh",
				},
			},
			Volumes: []Volume{{
				Name:        "TestVolume",
				Path:        "/super/path",
				Owner:       "test",
				Group:       "test",
				NeedsBackup: false,
			}},
			HealthCheck: HealthCheck{
				Type:       "tcp",
				State:      "test",
				Port:       8080,
				Path:       "/testpath/test",
				Parameters: nil,
			},
			HealthChecks: []HealthCheck{{
				Type:       "tcp",
				State:      "test",
				Port:       8080,
				Path:       "/testpath/test",
				Parameters: nil,
			}},
			ServiceAccounts: []ServiceAccount{{
				Type:   "test",
				Params: nil,
			}},
			Privileged: false,
			Configuration: []ConfigurationField{{
				Name:        "logging/test",
				Description: "set logging for test",
				Optional:    true,
				Encrypted:   true,
				Global:      false,
				Default:     "ERROR",
				Validation:  ValidationDescriptor{},
			}},
			Properties: Properties{"test": "44", "second": "wow"},
			EnvironmentVariables: []EnvironmentVariable{{
				Key:   "TestEnv",
				Value: "VALUE",
			}},
			Dependencies:         []string{"scm", "cas"},
			OptionalDependencies: []string{"nginx", "portainer"},
		}
		expectedDependencies := []Dependency{{
			Type: DependencyTypeDogu,
			Name: "scm",
		}, {
			Type: DependencyTypeDogu,
			Name: "cas",
		}}
		expectedOptionalDependencies := []Dependency{{
			Type: DependencyTypeDogu,
			Name: "nginx",
		}, {
			Type: DependencyTypeDogu,
			Name: "portainer",
		}}

		// when
		doguV2 := doguV1.CreateV2Copy()

		// then
		assert.Equal(t, doguV1.Name, doguV2.Name)
		assert.Equal(t, doguV1.Version, doguV2.Version)
    assert.Equal(t, doguV1.PublishedAt, doguV2.PublishedAt)
		assert.Equal(t, doguV1.DisplayName, doguV2.DisplayName)
		assert.Equal(t, doguV1.Description, doguV2.Description)
		assert.Equal(t, doguV1.Category, doguV2.Category)
		assert.Equal(t, doguV1.Tags, doguV2.Tags)
		assert.Equal(t, doguV1.Logo, doguV2.Logo)
		assert.Equal(t, doguV1.URL, doguV2.URL)
		assert.Equal(t, doguV1.Image, doguV2.Image)
		assert.Equal(t, doguV1.ExposedPorts, doguV2.ExposedPorts)
		assert.Equal(t, doguV1.ExposedCommands, doguV2.ExposedCommands)
		assert.Equal(t, doguV1.Volumes, doguV2.Volumes)
		assert.Equal(t, doguV1.HealthCheck, doguV2.HealthCheck)
		assert.Equal(t, doguV1.HealthChecks, doguV2.HealthChecks)
		assert.Equal(t, doguV1.ServiceAccounts, doguV2.ServiceAccounts)
		assert.Equal(t, doguV1.Privileged, doguV2.Privileged)
		assert.Equal(t, doguV1.Configuration, doguV2.Configuration)
		assert.Equal(t, doguV1.Properties, doguV2.Properties)
		assert.Equal(t, doguV1.EnvironmentVariables, doguV2.EnvironmentVariables)
		assert.Equal(t, expectedDependencies, doguV2.Dependencies)
		assert.Equal(t, expectedOptionalDependencies, doguV2.OptionalDependencies)
	})
}

func Test_newDoguV1Object(t *testing.T) {
	t.Run("Make deep conversion of old dogu schematic to new version", func(t *testing.T) {
		// given
		doguV2 := Dogu{
			Name:        "test",
			Version:     "0.1",
      PublishedAt: time.Date(2024, 10, 16, 7, 33, 45, 456, time.UTC),
			DisplayName: "displayName",
			Description: "description",
			Category:    "category",
			Tags:        []string{"tag1", "tag2"},
			Logo:        "logo",
			URL:         "url",
			Image:       "image",
			ExposedPorts: []ExposedPort{{
				Type:      "tcp",
				Container: 1,
				Host:      2,
			}},
			ExposedCommands: []ExposedCommand{
				{
					Name:        "TestCommand",
					Description: "Does nothing, just for testing purpose",
					Command:     "./test.sh",
				},
			},
			Volumes: []Volume{{
				Name:        "TestVolume",
				Path:        "/super/path",
				Owner:       "test",
				Group:       "test",
				NeedsBackup: false,
			}},
			HealthCheck: HealthCheck{
				Type:       "tcp",
				State:      "test",
				Port:       8080,
				Path:       "/testpath/test",
				Parameters: nil,
			},
			HealthChecks: []HealthCheck{{
				Type:       "tcp",
				State:      "test",
				Port:       8080,
				Path:       "/testpath/test",
				Parameters: nil,
			}},
			ServiceAccounts: []ServiceAccount{{
				Type:   "test",
				Params: nil,
			}},
			Privileged: false,
			Configuration: []ConfigurationField{{
				Name:        "logging/test",
				Description: "set logging for test",
				Optional:    true,
				Encrypted:   true,
				Global:      false,
				Default:     "ERROR",
				Validation:  ValidationDescriptor{},
			}},
			Properties: Properties{"test": "44", "second": "wow"},
			EnvironmentVariables: []EnvironmentVariable{{
				Key:   "TestEnv",
				Value: "VALUE",
			}},
			Dependencies: []Dependency{
				{Type: DependencyTypeDogu, Name: "scm"},
				{Type: DependencyTypeDogu, Name: "cas"},
			},
			OptionalDependencies: []Dependency{
				{Type: DependencyTypeDogu, Name: "nginx"},
				{Type: DependencyTypeDogu, Name: "portainer"},
			},
		}
		expectedDependencies := []string{"scm", "cas"}
		expectedOptionalDependencies := []string{"nginx", "portainer"}

		// when
		doguV1 := doguV2.CreateV1Copy()

		// then
		assert.Equal(t, doguV2.Name, doguV1.Name)
		assert.Equal(t, doguV2.Version, doguV1.Version)
    assert.Equal(t, doguV2.PublishedAt, doguV1.PublishedAt)
		assert.Equal(t, doguV2.DisplayName, doguV1.DisplayName)
		assert.Equal(t, doguV2.Description, doguV1.Description)
		assert.Equal(t, doguV2.Category, doguV1.Category)
		assert.Equal(t, doguV2.Tags, doguV1.Tags)
		assert.Equal(t, doguV2.Logo, doguV1.Logo)
		assert.Equal(t, doguV2.URL, doguV1.URL)
		assert.Equal(t, doguV2.Image, doguV1.Image)
		assert.Equal(t, doguV2.ExposedPorts, doguV1.ExposedPorts)
		assert.Equal(t, doguV2.ExposedCommands, doguV1.ExposedCommands)
		assert.Equal(t, doguV2.Volumes, doguV1.Volumes)
		assert.Equal(t, doguV2.HealthCheck, doguV1.HealthCheck)
		assert.Equal(t, doguV2.HealthChecks, doguV1.HealthChecks)
		assert.Equal(t, doguV2.ServiceAccounts, doguV1.ServiceAccounts)
		assert.Equal(t, doguV2.Privileged, doguV1.Privileged)
		assert.Equal(t, doguV2.Configuration, doguV1.Configuration)
		assert.Equal(t, doguV2.Properties, doguV1.Properties)
		assert.Equal(t, doguV2.EnvironmentVariables, doguV1.EnvironmentVariables)
		assert.Equal(t, expectedDependencies, doguV1.Dependencies)
		assert.Equal(t, expectedOptionalDependencies, doguV1.OptionalDependencies)
	})
}
