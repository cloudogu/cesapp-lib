package core

import (
	"fmt"
	"strings"

	"encoding/json"
)

// VolumeClient adds additional information for clients to create volumes.
type VolumeClient struct {
	// Name defines the actual client responsible to process this volume definition.
	Name string
	// Params contains generic data only known by the client.
	Params interface{}
}

// Volume is the volume struct of a dogu and will be used to define docker
// volumes
type Volume struct {
	Name        string
	Path        string
	Owner       string
	Group       string
	NeedsBackup bool
	Clients     []VolumeClient `json:"clients,omitempty"`
}

// GetClient retrieves a client with a given name and return a pointer to it. If a client does not exist a nil pointer
// and false are returned.
func (v *Volume) GetClient(clientName string) (*VolumeClient, bool) {
	for _, client := range v.Clients {
		if client.Name == clientName {
			return &client, true
		}
	}

	return nil, false
}

// UnmarshalJSON sets the default value for NeedsBackup. We are preventing an infinite loop by using a local Alias type
// to call json.Unmarshal again
func (v *Volume) UnmarshalJSON(data []byte) error {
	type Alias Volume
	a := &struct {
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	v.NeedsBackup = true
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	return nil
}

// HealthCheck struct will be used to do readyness and health checks for the
// final container
type HealthCheck struct {
	Type       string            // health check type tcp or state
	State      string            // expected state for state health check, default is ready
	Port       int               // port for tcp state check
	Path       string            // path for http check
	Parameters map[string]string // key value pairs for check specific parameters
}

// ExposedPort struct is used to define ports which are exported to the host
type ExposedPort struct {
	// Type contains the protocol type over which the container communicates (f. i. 'tcp').
	Type string
	// Container contains the mapped port on side of the container.
	Container int
	// Host contains the mapped port on side of the host.
	Host int
}

// GetType returns type of expose port, the default is tcp
func (ep *ExposedPort) GetType() string {
	if ep.Type == "" {
		return "tcp"
	}
	return ep.Type
}

// ExposedCommand struct represents a command which can be executed inside the
// dogu
type ExposedCommand struct {
	Name        string
	Description string
	Command     string
}

// Names for ExposedCommands correspond with actual dogu descriptor instance values. Do not change because these come
// with side effects.
const (
	ExposedCommandServiceAccountCreate = "service-account-create"
	ExposedCommandServiceAccountRemove = "service-account-remove"
	ExposedCommandBackupConsumer       = "backup-consumer"
	ExposedCommandPreBackup            = "pre-backup"
	ExposedCommandPostBackup           = "post-backup"
	ExposedCommandPostUpgrade          = "post-upgrade"
	ExposedCommandPreUpgrade           = "pre-upgrade"
	ExposedCommandUpgradeNotification  = "upgrade-notification"
)

// EnvironmentVariable struct represents custom parameters that can change
// the behaviour of a dogu build process
type EnvironmentVariable struct {
	Key   string
	Value string
}

// String returns a string representation of this EnvironmentVariable
func (env EnvironmentVariable) String() string {
	// Formatting of an EnvironmentVariable "ENV1=VALUE1"
	return env.Key + "=" + env.Value
}

// ServiceAccount struct can be used to get access to a other dogu.
type ServiceAccount struct {
	Type   string
	Params []string
}

// ConfigurationField describes a field of the dogu configuration which is stored in the registry.
type ConfigurationField struct {
	// Name contains the name of the key. It must not be empty. It must not contain leading or trailing slashes "/", but
	// it may contain directory keys delimited with slashes within the name.
	Name string
	// Description should mention the context and purpose of the config field in human readable format.
	Description string
	// Optional allows to have this config field unset.
	Optional bool
	// Encrypted marks this config field to contain a sensitive value that will be encrypted with the dogu's private key.
	Encrypted bool
	// Global marks this config field to contain a value that is available for all dogus.
	Global bool
	// Default defines a default value that may be evaluated if no value was configured, or the vallue is empty or even invalid.
	Default string
	// Validation configures a Validator that will be used to validate this config field.
	Validation ValidationDescriptor
}

// ValidationDescriptor describes how to determine if a config value is valid.
type ValidationDescriptor struct {
	// Type contains the name of the config value validator.
	Type string
	// Values may contain values that aid the selected validator. It is up to the selected validator whether this field is mandatory, optional, or unused.
	Values []string
}

// Properties describes generic properties of the dogu.
type Properties map[string]string

// Contains the different kind of types supported by dogu dependencies
const (
	DependencyTypeDogu    = "dogu"
	DependencyTypeClient  = "client"
	DependencyTypePackage = "package"
)

// Dependency contains the dependencies of the application and all their necessary information.
type Dependency struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Dogu defines an application for the CES. A dogu defines the image and
// meta information for the resulting container.
type Dogu struct {
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
	Dependencies         []Dependency
	OptionalDependencies []Dependency
}

// GetFullName returns the name of the dogu including its namespace
func (d *Dogu) GetFullName() string {
	return d.Name
}

// GetSimpleName returns the name of the dogu without the namespace
func (d *Dogu) GetSimpleName() string {
	parts := strings.Split(d.Name, "/")
	return parts[len(parts)-1]
}

// GetNamespace returns the namespace of the dogu without the name
func (d *Dogu) GetNamespace() string {
	parts := strings.Split(d.Name, "/")
	return parts[0]
}

// GetImageName returns the name of the docker image
func (d *Dogu) GetImageName() string {
	imageName := d.Image
	if d.Version != "" {
		imageName += ":" + d.Version
	}
	return imageName
}

// GetRegistryServerURI returns the name of the docker registry which is used by this dogu
func (d *Dogu) GetRegistryServerURI() string {
	return strings.TrimSuffix(d.Image, "/"+d.Name)
}

// GetAllDependenciesOfType returns all dependencies in accordance to the given dependency type.
func (d *Dogu) GetAllDependenciesOfType(dependencyType string) []Dependency {
	dependenciesAsNameList := make([]Dependency, 0)
	for _, dep := range d.Dependencies {
		if dep.Type == dependencyType {
			dependenciesAsNameList = append(dependenciesAsNameList, dep)
		}
	}
	for _, depOpt := range d.OptionalDependencies {
		if depOpt.Type == dependencyType {
			dependenciesAsNameList = append(dependenciesAsNameList, depOpt)
		}
	}
	return dependenciesAsNameList
}

// GetDependenciesOfType returns all dependencies in accordance to the given dependency type.
func (d *Dogu) GetDependenciesOfType(dependencyType string) []Dependency {
	dependenciesAsNameList := make([]Dependency, 0)
	deps := d.Dependencies
	for _, dep := range deps {
		if dep.Type == dependencyType {
			dependenciesAsNameList = append(dependenciesAsNameList, dep)
		}
	}
	return dependenciesAsNameList
}

// GetOptionalDependenciesOfType returns all optional dependencies in accordance to the given dependency type.
func (d *Dogu) GetOptionalDependenciesOfType(dependencyType string) []Dependency {
	dependenciesAsNameList := make([]Dependency, 0)
	optDeps := d.OptionalDependencies
	for _, depOpt := range optDeps {
		if depOpt.Type == dependencyType {
			dependenciesAsNameList = append(dependenciesAsNameList, depOpt)
		}
	}
	return dependenciesAsNameList
}

// HasExposedCommand checks if the dogu is a provider of a given command name. Example see constants like ExposedCommandServiceAccountCreate
func (d *Dogu) HasExposedCommand(commandName string) bool {
	for _, command := range d.ExposedCommands {
		if command.Name == commandName {
			return true
		}
	}
	return false
}

// GetExposedCommand returns a ExposedCommand for a given command name if it exists. Otherwise it returns nil.
// To test if a dogu has a command with a given command name use the HasExposedCommand method:
//
//	if dogu.HasExposedCommand(commandName) { doSomething() }
func (d *Dogu) GetExposedCommand(commandName string) *ExposedCommand {
	for _, command := range d.ExposedCommands {
		if command.Name == commandName {
			return &command
		}
	}
	return nil
}

// GetEnvironmentVariables returns environment variables formatted as array of key value objects
func (d *Dogu) GetEnvironmentVariables() []EnvironmentVariable {
	return d.EnvironmentVariables
}

// GetEnvironmentVariablesAsStringArray returns environment variables formatted as string array
func (d *Dogu) GetEnvironmentVariablesAsStringArray() []string {
	var environmentVariables []string
	for _, environmentVariable := range d.EnvironmentVariables {
		environmentVariables = append(environmentVariables, environmentVariable.String())
	}
	return environmentVariables
}

// DependsOn returns true if the dogu has a hard dependency to the given dogu
func (d *Dogu) DependsOn(name string) bool {
	dependencies := d.GetDependenciesOfType(DependencyTypeDogu)
	if dependencies == nil {
		return false
	}

	for _, dependency := range dependencies {
		if dependency.Name == name {
			return true
		}
	}

	return false
}

// GetVersion parses the dogu version and returns a struct which can be used to compare versions
func (d *Dogu) GetVersion() (Version, error) {
	version, err := ParseVersion(d.Version)
	if err != nil {
		return version, fmt.Errorf("failed to parse version %s of dogu %s: %w", d.Version, d.Name, err)
	}
	return version, nil
}

// IsEqualTo returns true if the other dogu has the same name and version.
func (d *Dogu) IsEqualTo(otherDogu *Dogu) (bool, error) {
	if d.Name != otherDogu.Name {
		return false, fmt.Errorf("only dogus with the same name can be compared")
	}

	version, err := d.GetVersion()
	if err != nil {
		return false, err
	}

	otherVersion, err := otherDogu.GetVersion()
	if err != nil {
		return false, err
	}

	return version.IsEqualTo(otherVersion), nil
}

// IsNewerThan returns true if the other dogu has the same name and has a higher version
func (d *Dogu) IsNewerThan(otherDogu *Dogu) (bool, error) {
	if d.Name != otherDogu.Name {
		return false, fmt.Errorf("only dogus with the same name can be compared")
	}

	version, err := d.GetVersion()
	if err != nil {
		return false, err
	}

	otherVersion, err := otherDogu.GetVersion()
	if err != nil {
		return false, err
	}

	return version.IsNewerThan(otherVersion), nil
}

// GetSimpleDoguName returns the dogu name without its namespace.
func GetSimpleDoguName(fullDoguName string) string {
	dogu := Dogu{Name: fullDoguName}
	return dogu.GetSimpleName()
}

// GetNamespace returns a dogu's namespace.
func GetNamespace(fullDoguName string) string {
	dogu := Dogu{Name: fullDoguName}
	return dogu.GetNamespace()
}

// CreateV1Copy converts this dogu object into a deep-copied DoguV1 object (for legacy reasons).
func (d *Dogu) CreateV1Copy() DoguV1 {
	dogu := DoguV1{}
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

	var dependencies []string
	for _, dependency := range d.GetDependenciesOfType(DependencyTypeDogu) {
		// the old format only allows dogu dependencies
		dependencies = append(dependencies, dependency.Name)
	}
	dogu.Dependencies = dependencies

	var optionalDependencies []string
	for _, dependency := range d.GetOptionalDependenciesOfType(DependencyTypeDogu) {
		// the old format only allows dogu dependencies
		optionalDependencies = append(optionalDependencies, dependency.Name)
	}
	dogu.OptionalDependencies = optionalDependencies

	return dogu
}

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
	return dogu, err
}

// ReadDogusFromString reads multiple dogus from a string and returns the API v2 representation.
func (d *DoguJsonV2FormatProvider) ReadDogusFromString(content string) ([]*Dogu, error) {
	var dogus []*Dogu
	err := json.Unmarshal([]byte(content), &dogus)
	return dogus, err
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

// ByDoguVersion implements sort.Interface for []Dogu to Dogus by their versions
type ByDoguVersion []*Dogu

// Len is the number of elements in the collection.
func (doguVersions ByDoguVersion) Len() int {
	return len(doguVersions)
}

// Swap swaps the elements with indexes i and j.
func (doguVersions ByDoguVersion) Swap(i, j int) {
	doguVersions[i], doguVersions[j] = doguVersions[j], doguVersions[i]
}

// Less reports whether the element with index i should sort before the element with index j.
func (doguVersions ByDoguVersion) Less(i, j int) bool {
	v1, err := ParseVersion(doguVersions[i].Version)
	if err != nil {
		GetLogger().Errorf("connot parse version %s for comparison", doguVersions[i].Version)
	}
	v2, err := ParseVersion(doguVersions[j].Version)
	if err != nil {
		GetLogger().Errorf("connot parse version %s for comparison", doguVersions[j].Version)
	}

	isNewer := v1.IsNewerThan(v2)
	return isNewer
}

// ContainsDoguWithName checks if a dogu is contained in a slice by comparing the full name (including namespace)
func ContainsDoguWithName(dogus []*Dogu, name string) bool {
	for _, dogu := range dogus {
		if dogu.Name == name {
			return true
		}
	}

	return false
}
