package core

// ConfigurationField describes a single dogu configuration field which is stored in the Cloudogu EcoSystem registry.
type ConfigurationField struct {
	// Name contains the name of the configuration key. The field is mandatory. It
	// must not contain leading or trailing slashes "/", but it may contain
	// directory keys delimited with slashes "/" within the name.
	//
	// The Name syntax is encouraged to consist of:
	//   - lower case latin characters
	//   - special characters underscore "_"
	//   - ciphers 0-9
	//
	// Example:
	//   - feedback_url
	//   - logging/root
	//
	Name string
	// Description offers context and purpose of the configuration field in human-readable format. This field is
	// optional, yet highly recommended to be set.
	//
	// Example:
	//   - "Set the root log level to one of ERROR, WARN, INFO, DEBUG or TRACE. Default is INFO"
	//   - "URL of the feedback service"
	//
	Description string
	// Optional allows to have this configuration field unset, otherwise a value must be set. This field is optional.
	// If unset, a value of `false` will be assumed.
	//
	// Example:
	//   - true
	//
	Optional bool
	// Encrypted marks this configuration field to contain a sensitive value that will be encrypted with the dogu's
	// public key. This field is optional. If unset, a value of `false` will be assumed.
	//
	// Example:
	//   - true
	//
	Encrypted bool
	// Global marks this configuration field to contain a value that is available for all dogus. This field is optional.
	// If unset, a value of `false` will be assumed.
	//
	// Example:
	//   - true
	//
	Global bool
	// Default defines a default value that may be evaluated if no value was configured, or the value is empty or even
	// invalid. This field is optional.
	//
	// Example:
	//   - "WARN"
	//   - "true"
	//   - "https://scm-manager.org/plugins"
	//
	Default string
	// Validation configures a validator that will be used to mark invalid or
	// out-of-range values for this configuration field. This field is optional.
	//
	// Example:
	//  "Validation": {
	//     "Type": "ONE_OF", // only allows one of these two values
	//     "Values": [
	//       "value 1",
	//       "value 2",
	//     ]
	//   }
	//
	//  "Validation": {
	//     "Type": "FLOAT_PERCENTAGE_HUNDRED" // valid values range between 0.0 and 100.0
	//  }
	//
	//  "Validation": {
	//    "Type": "BINARY_MEASUREMENT" // only allows suffixed integer values measured in byte, kibibyte, mebibyte, gibibyte
	//   }
	//
	//
	Validation ValidationDescriptor
	// IsDirectory marks this configuration field to be a directory. Multiple keys can be stored in this directory.
	// This field is optional. If unset, a value of `false` will be assumed.
	//
	// Example:
	//   - false
	//
	IsDirectory bool
}

// ValidationDescriptor describes how to determine if a config value is valid.
type ValidationDescriptor struct {
	// Type contains the name of the config value validator. This field is mandatory. Valid types are:
	//
	//   - ONE_OF
	//   - BINARY_MEASUREMENT
	//   - FLOAT_PERCENTAGE_HUNDRED
	//
	Type string
	// Values may contain values that aid the selected validator. The values may or
	// may not be optional, depending on the Type being used.
	// It is up to the selected validator whether this field is mandatory, optional,
	// or unused.
	Values []string
}
