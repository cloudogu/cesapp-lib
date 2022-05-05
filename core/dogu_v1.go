package core

import (
	"encoding/json"
)

// DoguV1 defines an application for the CES. A dogu defines the image and meta information for
// the resulting container. This schema is deprecated as it does not contain advanced expressions
// for dependencies.
type DoguV1 struct {
	Name                 string
	Version              string
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

func (d *DoguV1) CreateV2Copy() Dogu {
	dogu := Dogu{}
	dogu.Name = d.Name
	dogu.Version = d.Version
	dogu.DisplayName = d.DisplayName
	dogu.Description = d.Description
	dogu.Category = d.Category
	dogu.Tags = d.Tags
	dogu.Logo = d.Logo
	dogu.URL = d.URL
	dogu.Image = d.Image
	dogu.ExposedPorts = d.ExposedPorts
	dogu.ExposedCommands = d.ExposedCommands
	dogu.Volumes = d.Volumes
	dogu.HealthCheck = d.HealthCheck
	dogu.HealthChecks = d.HealthChecks
	dogu.ServiceAccounts = d.ServiceAccounts
	dogu.Privileged = d.Privileged
	dogu.Configuration = d.Configuration
	dogu.Properties = d.Properties
	dogu.EnvironmentVariables = d.EnvironmentVariables

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

type DoguJsonV1FormatProvider struct{}

func (d *DoguJsonV1FormatProvider) GetVersion() DoguApiVersion {
	return DoguApiV1
}

func (d *DoguJsonV1FormatProvider) ReadDoguFromString(content string) (*Dogu, error) {
	var doguV1 *DoguV1
	err := json.Unmarshal([]byte(content), &doguV1)
	if err != nil {
		return nil, err
	}
	dogu := doguV1.CreateV2Copy()
	return &dogu, err
}

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

func (d *DoguJsonV1FormatProvider) WriteDoguToString(doguV2 *Dogu) (string, error) {
	dogu := doguV2.CreateV1Copy()

	data, err := json.Marshal(dogu)
	if err != nil {
		return "", err
	}
	return string(data), err
}

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
