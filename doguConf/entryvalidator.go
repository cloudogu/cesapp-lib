package doguConf

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
)

const (
	// OneOfKey identifies the one-of-validator to select a single value among a list of valid values used in dogu.json.
	OneOfKey = "ONE_OF"
	// BinaryMeasurementKey identifies the binary prefix validator for integer values like "1024m" used in dogu.json.
	BinaryMeasurementKey = "BINARY_MEASUREMENT"
	// FloatPercentageHundredKey identifies the float percentage validator for float values between 0 and 100% used in dogu.json.
	FloatPercentageHundredKey = "FLOAT_PERCENTAGE_HUNDRED"
)

var entryValidatorTypes = map[string]entryValidatorCreator{
	OneOfKey:                  createOneOfValidator,
	BinaryMeasurementKey:      CreateBinaryMeasurementValidator,
	FloatPercentageHundredKey: CreateFloatPercentageValidator,
}

type entryValidatorCreator func(descriptor core.ValidationDescriptor) EntryValidator

// CreateEntryValidator creates the correct EntryValidator for the given descriptor.
func CreateEntryValidator(descriptor core.ValidationDescriptor) (EntryValidator, error) {
	if val, ok := entryValidatorTypes[descriptor.Type]; ok {
		return val(descriptor), nil
	}
	return nil, fmt.Errorf("no validator for type %s found", descriptor.Type)
}

// EntryValidator provides a Check method to check the validity of a single configuration entry.
type EntryValidator interface {
	// Check checks the given input from a configuration entry and returns nil if it is valid, otherwise an error.
	Check(input string) error
}
