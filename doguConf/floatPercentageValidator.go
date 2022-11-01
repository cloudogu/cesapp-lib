package doguConf

import (
	"regexp"
	"strconv"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

const (
	// validFloatPercentageRegex must contain at least a single preceding and succeeding digit.
	// The preceding digits are restricted to three digits while the succeeding digits are restricted to two.
	validFloatPercentageRegex = "^(\\d{1,3}).(\\d{1,2})$"
)

// FloatPercentageValidator checks if the given value matches a percentage in float form, f. i. '55.50' for 55.5%.
// FloatPercentageValidator validates only the actual config value and does not evaluate the field
// ValidationDescriptor.Values
// It is best to not configure the ValidationDescriptor.Values field when using this validator.
type FloatPercentageValidator struct {
	validationMatcher *regexp.Regexp
}

// CreateFloatPercentageValidator creates a float percentage measurement validator that validates a single config value.
func CreateFloatPercentageValidator(_ core.ValidationDescriptor) EntryValidator {
	validationMatcher, err := regexp.Compile(validFloatPercentageRegex)
	if err != nil {
		log.Errorf("error during compiling float percentage validator: %s", err.Error())
	}

	return &FloatPercentageValidator{
		validationMatcher: validationMatcher,
	}
}

// Check checks the correct regex and validates that the actual value is between 0%-100%
func (fpv *FloatPercentageValidator) Check(input string) error {
	matches := fpv.validationMatcher.MatchString(input)
	if !matches {
		return errors.Errorf("input '%s' should be a float number with at least one decimals place (f. ex. 22.5 for 22.5 percent); regex validation failed", input)
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.Errorf("input '%s' should be a float number with at least one decimals place (f. ex. 22.5 for 22.5 percent); parsing value failed", input)
	}

	if value < 0 || value > 100 {
		return errors.Errorf("input '%s' should be between 0 and 100; invalid percentage", input)
	}

	return nil
}
