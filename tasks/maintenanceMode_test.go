package tasks

import (
	"github.com/cloudogu/cesapp-lib/registry/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Set", "maintenance", "{\"title\": \"Maintenance\", \"text\": \"ipsum lorem\"}").Return(nil)

		// when
		err := ActivateMaintenanceMode("ipsum lorem", reg)

		// then
		require.Nil(t, err)
	})

	t.Run("should return an error on failure setting value to etcd", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Set", "maintenance", "{\"title\": \"Maintenance\", \"text\": \"ipsum lorem\"}").Return(assert.AnError)

		// when
		err := ActivateMaintenanceMode("ipsum lorem", reg)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("should return error on missing message", func(t *testing.T) {
		err := ActivateMaintenanceMode("", nil)

		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "Message text is missing")
	})

	t.Run("should return error on missing title", func(t *testing.T) {
		err := ActivateMaintenanceModeWithTitle("message", "", nil)

		require.NotNil(t, err)
		assert.Contains(t, err.Error(), "Message title is missing")
	})
}

func TestDeactivateMaintenanceMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Set", "maintenance", "{\"title\": \"Maintenance\", \"text\": \"test message\"}").Return(nil)
		globalConfig.On("Exists", "maintenance").Return(true, nil)
		globalConfig.On("Delete", "maintenance").Return(nil)
		err := ActivateMaintenanceMode("test message", reg)
		require.NoError(t, err)

		// when
		err = DeactivateMaintenanceMode(reg)

		// then
		require.Nil(t, err)
	})

	t.Run("should return an error if registry returns an error", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "maintenance").Return(false, assert.AnError)

		// when
		err := DeactivateMaintenanceMode(reg)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "assert.AnError general error for testing")
	})

	t.Run("should return no error if already activated", func(t *testing.T) {
		// given
		reg := mocks.NewRegistry(t)
		globalConfig := mocks.NewConfigurationContext(t)
		reg.On("GlobalConfig").Return(globalConfig)
		globalConfig.On("Exists", "maintenance").Return(false, nil)

		// when
		err := DeactivateMaintenanceMode(reg)

		// then
		require.Nil(t, err)
	})
}
