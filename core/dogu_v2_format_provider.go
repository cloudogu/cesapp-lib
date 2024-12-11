package core

import "encoding/json"

// DoguJsonV2FormatProvider provides methods to format Dogu results compatible to v2 API.
type DoguJsonV2FormatProvider struct{}

// GetVersion returns DoguApiV2 for this implementation.
func (d *DoguJsonV2FormatProvider) GetVersion() DoguApiVersion {
	return DoguApiV2
}

// ReadDoguFromString reads a dogu from a string and returns the API v2 representation.
func (d *DoguJsonV2FormatProvider) ReadDoguFromString(content string) (*Dogu, error) {
	var dogu *Dogu
	err := json.Unmarshal([]byte(content), &dogu)
	if err != nil {
		return nil, err
	}

	return dogu, nil
}

// ReadDogusFromString reads multiple dogus from a string and returns the API v2 representation.
func (d *DoguJsonV2FormatProvider) ReadDogusFromString(content string) ([]*Dogu, error) {
	var dogus []*Dogu
	err := json.Unmarshal([]byte(content), &dogus)
	if err != nil {
		return nil, err
	}

	return dogus, nil
}

// WriteDoguToString receives a single dogu and returns the API v2 representation.
func (d *DoguJsonV2FormatProvider) WriteDoguToString(dogu *Dogu) (string, error) {
	data, err := json.Marshal(dogu)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// WriteDogusToString receives a list of dogus and returns the API v2 representation.
func (d *DoguJsonV2FormatProvider) WriteDogusToString(dogu []*Dogu) (string, error) {
	data, err := json.Marshal(dogu)
	if err != nil {
		return "", err
	}
	return string(data), err
}
