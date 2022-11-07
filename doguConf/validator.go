package doguConf

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

// Validator checks if the given configurationField is valid.
type Validator interface {
	Check(core.ConfigurationField) error
}

// ConfigValidator is a Validator implementation which uses a configReader to read the values of the field.
type ConfigValidator struct {
	configReader configReader
}

type configReader interface {
	Get(string) (string, error)
	GetGlobal(string) (string, error)
}

// Check checks if a configurationField is valid.
func (c *ConfigValidator) Check(field core.ConfigurationField) error {

	var value string
	var err error

	if field.Global {
		value, err = c.configReader.GetGlobal(field.Name)
	} else {
		value, err = c.configReader.Get(field.Name)
	}

	if registry.IsKeyNotFoundError(err) && field.Optional {
		log.Debugf("Value for %s is empty but optional", field.Name)
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to get value of %s: %w", field.Name, err)
	}

	if field.Validation.Type == "" {
		log.Debugf("no validation configured for %s", field.Name)
		return nil
	}

	if field.Encrypted {
		log.Debugf("value of %s will not be validated since it is encrypted", field.Name)
		return nil
	}

	validator, err := CreateEntryValidator(field.Validation)
	if err != nil {
		return fmt.Errorf("failed to create entry validator: %w", err)
	}

	err = validator.Check(value)
	if err != nil {
		return fmt.Errorf("value %s for %s is not valid: %w", value, field.Name, err)
	}

	return nil
}
