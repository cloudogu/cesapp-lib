package config

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"strings"
)

var log = core.GetLogger()

type configReader interface {
	// Get takes a config key and returns a value if it exists (including the empty string). Otherwise, an
	// error will be returned (including errors concerning the absence of the config key).
	Get(key string) (value string, err error)
}

// RegistryKeyConfigValidator is a Validator implementation which uses a configReader to read the values of the field
type RegistryKeyConfigValidator struct{}

// CreateRegistryKeyConfigValidator creates a new validator that checks registry values.
func CreateRegistryKeyConfigValidator() *RegistryKeyConfigValidator {
	return &RegistryKeyConfigValidator{}
}

// CheckRegistryKey checks if a registry key match their registry value counterparts.
func (ekcv *RegistryKeyConfigValidator) CheckRegistryKey(value string, path string, config configReader, allowEmpty bool) error {
	err := validateRegistryKey(path, value, config, allowEmpty)
	if err != nil {
		return err
	}

	return nil
}

func validateRegistryKey(path string, value string, config configReader, allowEmptyKey bool) error {
	regValue, err := config.Get(path)
	if err != nil {
		if allowEmptyKey && isKeyNotFoundErr(err) {
			log.Warningf("Registry does not contain a key '%s'.", path)
			return nil
		}

		return fmt.Errorf("failed to get global config value for '%s' : %w", path, err)
	}

	if allowEmptyKey && regValue == "" {
		log.Warningf("Registry contains an empty value for '%s'.", path)
		return nil
	}

	if value != regValue {
		return fmt.Errorf("configuration (%s) is not equal to registry configuration (%s) for key '%s'", value,
			regValue, path)
	}

	return nil
}

func isKeyNotFoundErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "100: Key not found")
}
