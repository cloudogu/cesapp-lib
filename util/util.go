package util

import (
	"encoding/json"
	"github.com/cloudogu/cesapp-lib/core"
	"io"
	"io/ioutil"
	"os"
)

// ReadJSONFile reads a json file from the given path and mappes the fields to
// the given struct
func ReadJSONFile(structure interface{}, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &structure)
}

// GetContentOfFile all content of the file at the given path and returns it as string
func GetContentOfFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// WriteJSONFile writes a structure to a json file
func WriteJSONFile(structure interface{}, path string) error {
	data, err := json.Marshal(structure)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// WriteJSONFileIndented writes a structure indented to a json file
func WriteJSONFileIndented(structure interface{}, path string) error {
	data, err := json.MarshalIndent(structure, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// Exists returns true if the path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CloseButLogError closes IO closables and logs an error
func CloseButLogError(closer io.Closer, callContext string) {
	err := closer.Close()
	if err != nil {
		core.GetLogger().Errorf("could not close io.closer during %s. If the closing happened on a reading source this is probably a recoverable condition. %s",
			callContext, err.Error())
	}
}

// Contains checks if item is in the slice
func Contains(slice []string, item string) bool {
	if slice == nil {
		return false
	}
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

//Reverse returns a reversed string.
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
