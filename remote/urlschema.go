package remote

import (
	"fmt"
	"strings"
)

// URLSchema creates the url for the remote backend
type URLSchema interface {
	// Create returns the url for pushing dogu.json files to the backend
	Create(name string) string
	// Get returns the url for the latest version of the dogu
	Get(name string) string
	// GetVersion returns the url for one specific version of the dogu.json
	GetVersion(name string, version string) string
	// GetAll returns the url for fetching all dogu descriptors
	GetAll() string
	// GetVersionsOf returns the url for fetching all versions of a dogu
	GetVersionsOf(name string) string
}

// FormatURLSchema allows the definition of the url schema with string formats. Each format gets the base url plus the
// the function arguments as format parameters.
type FormatURLSchema struct {
	Endpoint            string
	CreateFormat        string
	GetFormat           string
	GetVersionFormat    string
	GetAllFormat        string
	GetVersionsOfFormat string
}

// Create returns the url for pushing dogu.json files to the backend
func (schema *FormatURLSchema) Create(name string) string {
	return fmt.Sprintf(schema.CreateFormat, schema.Endpoint, name)
}

// Get returns the url for the latest version of the dogu
func (schema *FormatURLSchema) Get(name string) string {
	return fmt.Sprintf(schema.GetFormat, schema.Endpoint, name)
}

// GetVersion returns the url for one specific version of the dogu.json
func (schema *FormatURLSchema) GetVersion(name string, version string) string {
	return fmt.Sprintf(schema.GetVersionFormat, schema.Endpoint, name, version)
}

// GetAll returns the url for fetching all dogu descriptors
func (schema *FormatURLSchema) GetAll() string {
	return fmt.Sprintf(schema.GetAllFormat, schema.Endpoint)
}

// GetVersionsOf returns the url for fetching all versions of a dogu
func (schema *FormatURLSchema) GetVersionsOf(name string) string {
	return fmt.Sprintf(schema.GetVersionsOfFormat, schema.Endpoint, name)
}

// NewDefaultURLSchema returns the url schema of the active backend.
func NewDefaultURLSchema(endpoint string) URLSchema {
	return &FormatURLSchema{
		Endpoint:            trimEndingSlash(endpoint),
		CreateFormat:        "%s/dogus/%s",
		GetFormat:           "%s/dogus/%s",
		GetVersionFormat:    "%s/dogus/%s/%s",
		GetAllFormat:        "%s/dogus/",
		GetVersionsOfFormat: "%s/dogus/%s/_versions",
	}
}

// NewIndexURLSchema returns the url schema which is required if the dogu descriptors are mirrored to a webserver.
func NewIndexURLSchema(endpoint string) URLSchema {
	return &FormatURLSchema{
		Endpoint:            trimEndingSlash(endpoint),
		CreateFormat:        "%s/%s/index.json",
		GetFormat:           "%s/%s/index.json",
		GetVersionFormat:    "%s/%s/%s/index.json",
		GetAllFormat:        "%s/index.json",
		GetVersionsOfFormat: "%s/%s/_versions.json",
	}
}

// NewURLSchemaByName returns the url schema for the given name or nil if no url schema exists with this name.
func NewURLSchemaByName(name string, endpoint string) URLSchema {
	switch name {
	case "":
		return NewDefaultURLSchema(endpoint)
	case "default":
		return NewDefaultURLSchema(endpoint)
	case "index":
		return NewIndexURLSchema(endpoint)
	}
	return nil
}

func trimEndingSlash(endpoint string) string {
	return strings.TrimSuffix(endpoint, "/")
}
