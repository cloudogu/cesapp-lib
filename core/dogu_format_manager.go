package core

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

const (
	DoguApiVersionUnknown DoguApiVersion = 0
	DoguApiV1             DoguApiVersion = 1
	DoguApiV2             DoguApiVersion = 2
)

type DoguFormatProvider interface {
	// ReadDoguFromString reads a dogu from the given filePath and returns an instance of a dogu.
	ReadDoguFromString(content string) (*Dogu, error)
	// ReadDogusFromString reads multiple dogus from the given filePath and returns a slice of dogu instances.
	ReadDogusFromString(content string) ([]*Dogu, error)
	// WriteDoguToString converts the given dogu to string representation and returns it.
	WriteDoguToString(dogu *Dogu) (string, error)
	// WriteDogusToString converts the given dogus to string representation and returns it.
	WriteDogusToString(dogu []*Dogu) (string, error)
	// GetVersion returns the api version which is used for this format provider.
	GetVersion() DoguApiVersion
}

// DoguFormatHandler is responsible for reading and writing dogus in different formats.
type DoguFormatHandler struct {
	Providers []DoguFormatProvider
}

type DoguApiVersion int

// formatHandlerInstance is a singleton instance of a DoguFormatHandler.
var formatHandlerInstance *DoguFormatHandler

func init() {
	// initialize static instance on load
	formatHandlerInstance = &DoguFormatHandler{}
	// register format providers
	formatHandlerInstance.RegisterFormatProvider(&DoguJsonV2FormatProvider{})
	formatHandlerInstance.RegisterFormatProvider(&DoguJsonV1FormatProvider{})
}

func (d *DoguFormatHandler) RegisterFormatProvider(provider DoguFormatProvider) {
	d.Providers = append(d.Providers, provider)
}

func (d *DoguFormatHandler) GetFormatProviders() []DoguFormatProvider {
	return d.Providers
}

func ReadDoguFromFile(filePath string) (*Dogu, DoguApiVersion, error) {
	fileContent, err := GetContentOfFile(filePath)
	if err != nil {
		return nil, DoguApiVersionUnknown, errors.Wrap(err, "Cannot read dogu from invalid file")
	}

	return ReadDoguFromString(fileContent)
}

func ReadDogusFromFile(filePath string) ([]*Dogu, DoguApiVersion, error) {
	fileContent, err := GetContentOfFile(filePath)
	if err != nil {
		return nil, DoguApiVersionUnknown, errors.Wrap(err, "Cannot read dogus from invalid file")
	}

	return ReadDogusFromString(fileContent)
}

func ReadDoguFromString(content string) (*Dogu, DoguApiVersion, error) {
	var firstError error
	for _, provider := range formatHandlerInstance.Providers {
		dogu, err := provider.ReadDoguFromString(content)
		if err != nil && firstError == nil {
			// only save the first error, i.e., the error for the newest format
			firstError = err
		} else if err == nil {
			return dogu, provider.GetVersion(), err
		}
	}
	return nil, DoguApiVersionUnknown, firstError
}

func ReadDogusFromString(content string) ([]*Dogu, DoguApiVersion, error) {
	var firstError error
	for _, provider := range formatHandlerInstance.Providers {
		dogus, err := provider.ReadDogusFromString(content)
		if err != nil && firstError == nil {
			// only save the first error, i.e., the error for the newest format
			firstError = err
		} else if err == nil {
			return dogus, provider.GetVersion(), err
		}
	}

	return nil, DoguApiVersionUnknown, firstError
}

// WriteDoguToFile writes the dogu to the given file. Uses the default format (first registered).
func WriteDoguToFile(filePath string, dogu *Dogu) error {
	return WriteDoguToFileWithFormat(filePath, dogu, &DoguJsonV2FormatProvider{})
}

// WriteDogusToFile writes all dogus to the given file. Uses the default format (first registered).
func WriteDogusToFile(filePath string, dogus []*Dogu) error {
	return WriteDogusToFileWithFormat(filePath, dogus, &DoguJsonV2FormatProvider{})
}

// WriteDoguToString writes the dogu and return the string representation of the specified format.
func WriteDoguToString(dogu *Dogu) (string, error) {
	provider := &DoguJsonV2FormatProvider{}
	data, err := provider.WriteDoguToString(dogu)
	if err != nil {
		return "", err
	}
	return data, err
}

// WriteDogusToString writes all dogus and return the string representation of the specified format.
func WriteDogusToString(dogus []*Dogu) (string, error) {
	provider := &DoguJsonV2FormatProvider{}
	data, err := provider.WriteDogusToString(dogus)
	if err != nil {
		return "", err
	}
	return data, err
}

// WriteDoguToFileWithFormat writes the dogu to the given file using a specified format.
func WriteDoguToFileWithFormat(filePath string, dogu *Dogu, formatProvider DoguFormatProvider) error {
	data, err := formatProvider.WriteDoguToString(dogu)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// WriteDogusToFileWithFormat writes all dogus to the given file using a specified format.
func WriteDogusToFileWithFormat(filePath string, dogus []*Dogu, formatProvider DoguFormatProvider) error {
	data, err := formatProvider.WriteDogusToString(dogus)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// GetContentOfFile all content of the file at the given path and returns it as string
func GetContentOfFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), err
}
