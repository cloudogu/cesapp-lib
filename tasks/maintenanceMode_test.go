package tasks

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateMaintenanceMode(t *testing.T) {
	reg := &registry.MockRegistry{}

	err := ActivateMaintenanceMode("ipsum lorem", reg)

	require.Nil(t, err)
	actual, _ := reg.GlobalConfig().Get("maintenance")
	expected := "{\"title\": \"Maintenance\", \"text\": \"ipsum lorem\"}"
	assert.Equal(t, expected, actual)
}

func TestActivateMaintenanceMode_errorOnMissingMessage(t *testing.T) {
	reg := &registry.MockRegistry{}

	err := ActivateMaintenanceMode("", reg)

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Message text is missing")
	exists, _ := reg.GlobalConfig().Exists("maintenance")
	assert.False(t, exists)
}

func TestDeactivateMaintenanceMode(t *testing.T) {
	reg := &registry.MockRegistry{}
	_ = ActivateMaintenanceMode("test message", reg)

	err := DeactivateMaintenanceMode(reg)

	require.Nil(t, err)
	exists, _ := reg.GlobalConfig().Exists("maintenance")
	assert.False(t, exists)
}

func TestDeactivateMaintenanceMode_noErrorIfAlreadyDeactivated(t *testing.T) {
	reg := &registry.MockRegistry{}

	err := DeactivateMaintenanceMode(reg)

	require.Nil(t, err)
	exists, _ := reg.GlobalConfig().Exists("maintenance")
	assert.False(t, exists)
}
