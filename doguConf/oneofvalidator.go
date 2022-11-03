package doguConf

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
)

// oneOfValidator checks if the given value is element of a predefined set
type oneOfValidator struct {
	values []string
}

func createOneOfValidator(descriptor core.ValidationDescriptor) EntryValidator {
	return &oneOfValidator{values: descriptor.Values}
}

// Check checks if the input is in the set of valid values
func (o *oneOfValidator) Check(input string) error {
	for _, v := range o.values {
		if v == input {
			return nil
		}
	}
	return fmt.Errorf("input should be one of %q", o.values)
}
