package doguConf

import (
	"bufio"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry"
	"os"
	"strings"
)

type doguConfigurationContext interface {
	// Set sets a configuration value in current context
	Set(key, value string) error
	// Get returns a configuration value from the current context
	Get(key string) (string, error)
	// Exists returns true if configuration key exists in the current context
	Exists(key string) (bool, error)
	// Delete removes a configuration key and value from the current context
	Delete(key string) error
}

// DoguConfigurationEditor struct is able to edit registry configuration values of a dogu.
type DoguConfigurationEditor struct {
	ConfigurationContext doguConfigurationContext
	PublicKey            *keys.PublicKey
	Writer               FieldWriter
	Reader               FieldReader
	DeleteOnEmpty        bool
	out                  *bufio.Writer
}

// NewDoguConfigurationEditor creates a new DoguConfigurationEditor struct,
// with stdout as writer and stdin as reader for the given dogu
func NewDoguConfigurationEditor(doguConfig registry.ConfigurationContext, publicKey *keys.PublicKey) (*DoguConfigurationEditor, error) {
	return &DoguConfigurationEditor{
		ConfigurationContext: doguConfig,
		PublicKey:            publicKey,
		Reader:               &defaultFieldReader{bufio.NewReader(os.Stdin)},
		Writer:               &defaultFieldWriter{bufio.NewWriter(os.Stdout)},
		out:                  bufio.NewWriter(os.Stdout),
	}, nil
}

// FieldWriter can print a configuration field to a underlying system such as stdout.
type FieldWriter interface {
	Print(field core.ConfigurationField, currentValue string) error
}

type defaultFieldWriter struct {
	Writer *bufio.Writer
}

// Print prints various configuration parts to stdout.
func (writer *defaultFieldWriter) Print(field core.ConfigurationField, currentValue string) error {
	_, err := fmt.Fprintln(writer.Writer, field.Description)
	if err != nil {
		return fmt.Errorf("failed to write to writer: %w", err)
	}

	if field.Validation.Type != "" {
		_, err := fmt.Fprintf(writer.Writer, "Available values are: %v\n", field.Validation.Values)
		if err != nil {
			return fmt.Errorf("failed to write to writer: %w", err)
		}
	}

	_, err = fmt.Fprintf(writer.Writer, "%s (%s): ", field.Name, currentValue)
	if err != nil {
		return fmt.Errorf("failed to write to writer: %w", err)
	}

	err = writer.Writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush output writer: %w", err)
	}

	return nil
}

// FieldReader can read the input for a configuration field from an underlying system such as stdin.
type FieldReader interface {
	// Read reads the input line.
	Read() (string, error)
}

type defaultFieldReader struct {
	Reader *bufio.Reader
}

// Read reads the input for a configuration field from an underlying system.
func (reader *defaultFieldReader) Read() (string, error) {
	line, err := reader.Reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	return strings.TrimSpace(line), nil
}

func (editor *DoguConfigurationEditor) editConfigurationField(field core.ConfigurationField) error {
	currentValue, err := editor.GetCurrentValue(field)
	if err != nil {
		return err
	}

	input, err := editor.printFieldAndReadInput(field, currentValue)
	if err != nil {
		return fmt.Errorf("failed to read input for field %s: %w", field.Name, err)
	}

	err = editor.SetFieldToValue(field, input)
	if err != nil {
		return err
	}

	return nil
}

func (editor *DoguConfigurationEditor) printFieldAndReadInput(field core.ConfigurationField, currentValue string) (string, error) {
	var input string
	var err error
	var validator EntryValidator
	hasValidator := field.Validation.Type != ""
	if hasValidator {
		validator, err = CreateEntryValidator(field.Validation)
		if err != nil {
			return "", fmt.Errorf("failed to create validator: %w", err)
		}
	}

	for {
		err = editor.Writer.Print(field, currentValue)
		if err != nil {
			return "", fmt.Errorf("failed to print field %s: %w", field.Name, err)
		}

		input, err = editor.Reader.Read()
		if err != nil {
			return "", fmt.Errorf("could not read input: %w", err)
		}

		if !hasValidator || (field.Optional && input == "") {
			break
		} else if hasValidator {
			parseErr := validator.Check(input)
			if parseErr == nil {
				break
			} else {
				log.Infof("Cannot apply '%s' to field '%s'. Reason: %s \n\n", input, field.Name, parseErr)
			}
		} else {
			log.Infof("Cannot apply '%s' to field '%s'.\n\n", input, field.Name)
		}
	}

	fmt.Println()
	return input, nil
}

// EditConfiguration prints registry keys to writer and read values from reader.
func (editor *DoguConfigurationEditor) EditConfiguration(fields []core.ConfigurationField) error {
	for _, field := range fields {
		if !field.Global {
			err := editor.editConfigurationField(field)
			if err != nil {
				return err
			}
		} else {
			log.Debug("skip global field", field.Name)
		}
	}
	return nil
}

// GetCurrentValue returns a value for a given ConfigurationField if it exists, otherwise it returns an error.
func (editor *DoguConfigurationEditor) GetCurrentValue(field core.ConfigurationField) (string, error) {
	exists, err := editor.ConfigurationContext.Exists(field.Name)
	if err != nil {
		return "", fmt.Errorf("failed to check if key %s exists: %w", field.Name, err)
	}

	if exists {
		var value string
		// do not return encrypted value
		if field.Encrypted {
			value = "_encrypted_"
		} else {
			value, err = editor.ConfigurationContext.Get(field.Name)
			if err != nil {
				return "", fmt.Errorf("failed to get key %s from registry: %w", field.Name, err)
			}
		}
		return value, nil
	}

	return "", nil
}

// SetFieldToValue set the Field as value into the editor.
func (editor *DoguConfigurationEditor) SetFieldToValue(field core.ConfigurationField, value string) error {
	var err error
	if value == "" {
		err = editor.handleEmptyFieldInput(field)
	} else if field.Encrypted {
		err = editor.setEncryptedFieldValue(field, value)
	} else {
		err = editor.setFieldValue(field, value)
	}
	return err
}

func (editor *DoguConfigurationEditor) setEncryptedFieldValue(field core.ConfigurationField, value string) error {
	log.Debugf("encrypt field %s", field.Name)
	encryptedValue, err := editor.encryptValue(field, value)
	if err != nil {
		return err
	}

	return editor.setFieldValue(field, encryptedValue)
}

func (editor *DoguConfigurationEditor) encryptValue(field core.ConfigurationField, value string) (string, error) {
	if editor.PublicKey == nil {
		return "", fmt.Errorf("no public key for encryption of field %s found", field.Name)
	}

	encryptedValue, err := editor.PublicKey.Encrypt(value)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt field %s: %w", field.Name, err)
	}
	return encryptedValue, nil
}

func (editor *DoguConfigurationEditor) setFieldValue(field core.ConfigurationField, value string) error {
	log.Debugf("set new value for field %s", field.Name)
	err := editor.ConfigurationContext.Set(field.Name, value)
	if err != nil {
		return fmt.Errorf("failed to set value for field %s: %w", field.Name, err)
	}
	return nil
}

func (editor *DoguConfigurationEditor) handleEmptyFieldInput(field core.ConfigurationField) error {
	if editor.DeleteOnEmpty {
		log.Debugf("delete value for %s, because input was empty and DeleteOnEmpty option was used", field.Name)
		err := editor.DeleteField(field)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("received empty input for key %s, do not change value", field.Name)
	}
	return nil
}

func (editor *DoguConfigurationEditor) DeleteField(field core.ConfigurationField) error {
	exists, err := editor.ConfigurationContext.Exists(field.Name)
	if err != nil {
		return fmt.Errorf("failed to check, if field %s exists: %w", field.Name, err)
	}
	if exists {
		err := editor.ConfigurationContext.Delete(field.Name)
		if err != nil {
			return fmt.Errorf("failed to delete field %s: %w", field.Name, err)
		}
	} else {
		log.Debugf("skip deletion of non existing field %s", field.Name)
	}
	return nil
}

// HasConfiguration returns true if the dogu has configuration fields applied, otherwise false.
func HasConfiguration(dogu *core.Dogu) bool {
	return len(dogu.Configuration) > 0
}
