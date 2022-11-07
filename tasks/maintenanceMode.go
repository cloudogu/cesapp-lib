package tasks

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

// MaintenanceRegistryKey contains the registry key name which points to the Maintenance Mode text.
const MaintenanceRegistryKey = "maintenance"
const defaultMaintenanceTitle = "Maintenance"
const activateMaintenanceModeDescription = "Activate maintenance mode"

var log = core.GetLogger()

// ActivateMaintenanceMode activates the maintenance mode in the registry. The given message text will be presented
// to the user. The message text must not be empty.
func ActivateMaintenanceMode(messageText string, registry registry.Registry) error {
	return ActivateMaintenanceModeWithTitle(messageText, defaultMaintenanceTitle, registry)
}

// ActivateMaintenanceModeWithTitle activates the maintenance mode in the registry. The given message text will be presented
// to the user. The message text must not be empty. The given title will be presented to the user. The title must not be empty.
func ActivateMaintenanceModeWithTitle(messageText string, title string, registry registry.Registry) error {
	if messageText == "" {
		return fmt.Errorf("could not activate maintenance mode. Message text is missing")
	}
	if title == "" {
		return fmt.Errorf("could not activate maintenance mode. Message title is missing")
	}

	log.Info(activateMaintenanceModeDescription)

	json := fmt.Sprintf("{\"title\": \"%s\", "+
		"\"text\": \"%s\"}", title, messageText)

	err := registry.GlobalConfig().Set(MaintenanceRegistryKey, json)
	if err != nil {
		return fmt.Errorf("failed to activate maintenance mode: %w", err)
	}
	return nil
}

// DeactivateMaintenanceMode deactivates the maintenance mode in the registry if it exists.
// This function will not return an error if the key is not present in the registry.
func DeactivateMaintenanceMode(registry registry.Registry) error {
	exists, err := registry.GlobalConfig().Exists(MaintenanceRegistryKey)
	if err != nil {
		return fmt.Errorf("could not check if maintenance mode key exists: %w", err)
	}

	if exists {
		return registry.GlobalConfig().Delete(MaintenanceRegistryKey)
	}

	return nil
}
