package core

import (
	"encoding/json"
	"time"
)

// DoguV1 defines an application for the CES. A dogu defines the image and meta information for
// the resulting container.
//
// Deprecated: This schema is deprecated as it does not contain advanced expressions
// for dependencies. Please use core.DoguV2 instead.
type DoguV1 struct {
	Name                 string
	Version              string
	PublishedAt          time.Time
	DisplayName          string
	Description          string
	Category             string
	Tags                 []string
	Logo                 string
	URL                  string
	Image                string
	ExposedPorts         []ExposedPort
	ExposedCommands      []ExposedCommand
	Volumes              []Volume
	HealthCheck          HealthCheck // deprecated use HealthChecks
	HealthChecks         []HealthCheck
	ServiceAccounts      []ServiceAccount
	Privileged           bool
	Configuration        []ConfigurationField
	Properties           Properties
	EnvironmentVariables []EnvironmentVariable
	Dependencies         []string
	OptionalDependencies []string
}

// CreateV2Copy creates a deep DoguV2 copy from an existing DoguV1 object.
func (d *DoguV1) CreateV2Copy() Dogu {
	dogu := Dogu{
		Name:                 d.Name,
		Version:              d.Version,
		PublishedAt:          d.PublishedAt,
		DisplayName:          d.DisplayName,
		Description:          d.Description,
		Category:             d.Category,
		Tags:                 d.Tags,
		Logo:                 d.Logo,
		URL:                  d.URL,
		Image:                d.Image,
		ExposedPorts:         d.ExposedPorts,
		ExposedCommands:      d.ExposedCommands,
		Volumes:              d.Volumes,
		HealthCheck:          d.HealthCheck,
		HealthChecks:         d.HealthChecks,
		ServiceAccounts:      d.ServiceAccounts,
		Privileged:           d.Privileged,
		Configuration:        d.Configuration,
		Properties:           d.Properties,
		EnvironmentVariables: d.EnvironmentVariables,
	}

	// upgrade dependencies to new version
	// the old schema only contained dogus as dependencies
	var dependencies []Dependency
	for _, dependencyOld := range d.Dependencies {
		dependencyNew := Dependency{
			Type: DependencyTypeDogu,
			Name: dependencyOld,
		}
		dependencies = append(dependencies, dependencyNew)
	}
	dogu.Dependencies = dependencies

	// upgrade optional dependencies to new version
	// the old schema only contained dogus as optional dependencies
	var optionalDependencies []Dependency
	for _, dependencyOld := range d.OptionalDependencies {
		dependencyNew := Dependency{
			Type: DependencyTypeDogu,
			Name: dependencyOld,
		}
		optionalDependencies = append(optionalDependencies, dependencyNew)
	}
	dogu.OptionalDependencies = optionalDependencies

	return dogu
}

// DoguJsonV1FormatProvider provides methods to format Dogu results compatible to v1 API.
type DoguJsonV1FormatProvider struct{}

// GetVersion returns DoguApiV1 for this implementation.
func (d *DoguJsonV1FormatProvider) GetVersion() DoguApiVersion {
	return DoguApiV1
}

// ReadDoguFromString reads a dogu from a string and returns the API v1 representation.
func (d *DoguJsonV1FormatProvider) ReadDoguFromString(content string) (*Dogu, error) {
	var doguV1 *DoguV1
	err := json.Unmarshal([]byte(content), &doguV1)
	if err != nil {
		return nil, err
	}
	dogu := doguV1.CreateV2Copy()
	return &dogu, err
}

// ReadDogusFromString reads multiple dogus from a string and returns the API v1 representation.
func (d *DoguJsonV1FormatProvider) ReadDogusFromString(content string) ([]*Dogu, error) {
	var dogusV1 []*DoguV1
	err := json.Unmarshal([]byte(content), &dogusV1)
	if err != nil {
		return nil, err
	}

	var dogus []*Dogu
	for _, doguV1 := range dogusV1 {
		doguV2 := doguV1.CreateV2Copy()

		dogus = append(dogus, &doguV2)
	}

	return dogus, err
}

// WriteDoguToString receives a single dogu and returns the API v1 representation.
func (d *DoguJsonV1FormatProvider) WriteDoguToString(doguV2 *Dogu) (string, error) {
	dogu := doguV2.CreateV1Copy()

	data, err := json.Marshal(dogu)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// WriteDogusToString receives a list of dogus and returns the API v1 representation.
func (d *DoguJsonV1FormatProvider) WriteDogusToString(doguV2List []*Dogu) (string, error) {
	var doguV1List []*DoguV1
	for _, doguV2 := range doguV2List {
		doguV1 := doguV2.CreateV1Copy()
		doguV1List = append(doguV1List, &doguV1)
	}

	data, err := json.Marshal(doguV1List)
	if err != nil {
		return "", err
	}
	return string(data), err
}
