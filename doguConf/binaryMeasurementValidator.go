package doguConf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
)

var log = core.GetLogger()

const (
	// ValidBinaryUnits contains available binary measurement units that are manifolds of 1024.
	ValidBinaryUnits = "bkmg"
	// validBinaryUnitsRegex must contain at least a single digit up to 19 digits (because of MaxInt64) followed by a
	// binary measurement unit.
	validBinaryUnitsRegex = "^(\\d{1,19})([" + ValidBinaryUnits + "])$"
)

// BinaryMeasurementValidator checks if the given value matches a binary measurement, f. i. '1024k'.
// BinaryMeasurementValidator validates only the actual config value and does not evaluate the field
//	ValidationDescriptor.Values
// It is best to not configure the ValidationDescriptor.Values field when using this validator.
type BinaryMeasurementValidator struct {
	validationMatcher *regexp.Regexp
}

// CreateBinaryMeasurementValidator creates a binary measurement validator that validates a single config value.
func CreateBinaryMeasurementValidator(core.ValidationDescriptor) EntryValidator {
	validationMatcher, err := regexp.Compile(validBinaryUnitsRegex)
	if err != nil {
		log.Errorf("error during compiling binary measurement validator: %s", err.Error())
	}

	return &BinaryMeasurementValidator{
		validationMatcher: validationMatcher,
	}
}

// Check checks if the input matches a zero or positive integer followed by a binary measurement, namely: b, k, m, g.
func (bpv *BinaryMeasurementValidator) Check(input string) error {
	matches := bpv.validationMatcher.MatchString(input)
	if matches {
		return nil
	}

	allowedBinaryUnits := strings.Join(strings.Split(ValidBinaryUnits, ""), ",")
	return fmt.Errorf("input '%s' should be an integer with a binary measurement (f. ex. 2k for 2048 bytes); valid units are: %s", input, allowedBinaryUnits)
}

// SplitValueAndUnit returns the input's integer part as well as the binary unit. Be sure to call Check() beforehand.
func (bpv *BinaryMeasurementValidator) SplitValueAndUnit(input string) (string, string) {
	fullStringAndDigitGroup := bpv.validationMatcher.FindStringSubmatch(input)
	valueGroup := 1
	unitGroup := 2
	return fullStringAndDigitGroup[valueGroup], fullStringAndDigitGroup[unitGroup]
}
