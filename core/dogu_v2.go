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
	// Name identifies the volume.
	Name string
	// Path to the directory or file where the volume will be mounted inside the dogu.
	Path string
	// Owner contains the uid of the user owning this volume.
	Owner string
	// Group contains the gid of the group owning this volume.
	Group string
	// NeedsBackup defines if backups need to be created for the volume.
	NeedsBackup bool
	// Clients adds client-specific configurations for the volume.
	Clients []VolumeClient `json:"Clients,omitempty"`
}

// GetClient retrieves a client with a given name and return a pointer to it. If a client does not exist a nil pointer
// and false are returned.
func (v *Volume) GetClient(clientName string) (*VolumeClient, bool) {
	if v.Clients == nil {
		return nil, false
	}

	for i := range v.Clients {
		if v.Clients[i].Name == clientName {
			return &v.Clients[i], true
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
	// Type contains the name of the service on which the account should be created.
	Type string
	// Params contains additional arguments necessary for the service account creation.
	Params []string
	// Kind defines the kind of service on which the account should be created, e.g. `dogu` or `k8s`.
	// Reading this property and creating a corresponding service account is up to the client.
	// If empty, a default value of `dogu` should be assumed.
	Kind string `json:"Kind,omitempty"`
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
	// Name contains the dogu's full qualified name which consists of the dogu namespace and the dogu simple name,
	// delimited by a single forward slash "/". This field is mandatory.
	//
	// The dogu namespace allows to regulate access to dogus in that namespace. There are three reserved dogu
	// namespaces: The namespaces `official` and `k8s` are open to all users without any further costs. In contrast to
	// that is the namespace `premium` is open to subscription users, only.
	//
	// The namespace syntax is encouraged to consist of:
	//   - lower case latin characters
	//   - special characters underscore "_", minus "-"
	//   - ciphers 0-9
	//   - an overall length of less than 200 characters
	//
	// The dogu simple name allows to address in multiple ways. The simple name will be the part of the URI of the
	// Cloudogu EcoSystem to address a URI part (if the dogu provides an exposed UI). Also, the simple name will be used
	// to address the dogu after the installation process (f. i. to start, stop or remove a dogu), or to address
	// generated resources that belong to the dogu.
	//
	// The simple name syntax must be an DNS-compatible identifier and is encouraged to consist of
	//   - lower case latin characters
	//   - special characters underscore "_", minus "-"
	//   - ciphers 0-9
	//   - an overall length of less than 20 characters
	//
	// It is recommended to use the same full qualified dogu name within the dogu's Dockerfile as environment variable
	// `NAME`.
	//
	// Examples:
	//   - official/redmine
	//   - premium/confluence
	//   - foo-1/bar-2
	//
	Name string
	// Version defines the actual version of the dogu. This field is mandatory.
	//
	// The version follows the format from semantic versioning and additionally is split in two parts.
	// The application version and the dogu version.
	//
	// An example would be 1.7.8-1 or 2.2.0-4. The first part of the version (e.g. 1.7.8) represents the
	// version of the application (e.g. the nginx version in the nginx dogu). The second part represents the version
	// of the dogu and for an initial release it should start at 1 (e.g. 1.7.8-1).
	//
	// If the application does not change but e.g. there are changes in the startup script of the dogu the new version
	// should be 1.7.8-2. If the application itself changes (e.g. there is a nginx upgrade) the new version should be
	// 1.7.9-1. Notice that in this case the version of the dogu should be set to 1 again.
	//
	// Whereas the dogu struct is the core place for the version and is used by the cesapp in various processes like
	// installation and release the version should also be placed as a label in the dockerfile from the dogu.
	//
	// Example versions in the dogu.json:
	//  1.7.8-1
	//  2.2.0-4
	//
	// Recommended example in the Dockerfile:
	//  LABEL maintainer="hello@cloudogu.com" \
	//    NAME="official/nginx" \
	//    VERSION="1.23.2-1"
	//
	Version string
	// DisplayName is the name of the dogu which is used in ui frontends to represent the dogu. This field is mandatory.
	//
	// Usages:
	// In the setup of the ecosystem the display name of the dogu is used to select it for installation.
	//
	// For dogus with a web ui an important location is the warp menu where you can navigate with a click of the
	// display name to the dogu.
	//
	// Another location is the textual output of tools like the cesapp or the k8s-dogu-operator where the name is used
	// in commands like list upgradeable dogus.
	//
	// The display name may consist of
	//   - lower and upper case latin characters where the first is upper case
	//   - special characters minus "-", ampersand "&"
	//   - ciphers 0-9
	//   - an overall length of less than 30 characters
	//
	// Examples:
	//  Jenkins CI
	//  Backup & Restore
	//  SCM-Manager
	//  Smeagol
	//
	DisplayName string
	// Description describes in a few words what the dogu is and maybe do. This field is mandatory.
	//
	// It is used in the setup of the ecosystem in the dogu selection.
	// Therefor the description should give an uninformed user a brief hint what the dogu is
	// and maybe the function the dogu fulfills.
	//
	// The description may consist of
	//   - lower and upper case latin characters where the first is upper case
	//   - special characters minus "-", ampersand "&"
	//   - ciphers 0-9
	//   - an overall length of less than 30 words
	//
	// Examples:
	//  "Jenkins Continuous Integration Server"
	//  "MySQL - Relational database"
	//  "The Nexus Repository is like the local warehouse where all the parts and finished goods used in your
	//  software supply chain are stored and distributed."
	//
	Description string
	// Category organizes the dogus in three categories. This field is mandatory.
	//
	// These categories are fixed and must be either:
	//
	// "Development Apps" - For regular dogus which should be used by a regular user of the ecosystem,
	// "Administration Apps" - For dogus which should be used by a user with administration rights, or
	// "Base" - For dogus which are important for the overall system.
	//
	// The categories "Development Apps" and "Administration Apps" are represented in the warp menu to order the dogus.
	//
	// Example dogus for each category:
	//  "Development Apps": Redmine, SCM-Manager, Jenkins
	//  "Administration Apps": Backup & Restore, User Management
	//  "Base": Nginx, Registrator, OpenLDAP
	//
	Category string
	// Tags is a slice of one-word-tags which are in connection with the dogu. This field is optional.
	//
	// If the dogu should be displayed in the warp menu the tag "warp" is necessary.
	// Actually other tags won't be processed.
	//
	// Examples for e.g. Jenkins:
	//  {"warp", "build", "ci", "cd"}
	//
	Tags []string
	// Logo originally represented a URI to a web picture depicting the dogu tool. This field is optional.
	//
	// Deprecated: The Cloudogu EcoSystem does not facilitate the logo URI. It is a candidate for removal.
	// Other options of representing a tool or application can be:
	//   - embed the logo in the dogu's Git repository (if public)
	//   - provide the logo in to dogu UI (if the dogu provides one)
	//
	Logo string
	// URL may link the website to the original tool vendor. This field is optional. Like Logo, the Cloudogu EcoSystem
	// does not facilitate this information. Anyhow, in a public dogu repository in which a dogu vendor re-packages a
	// third party application the URL may point users to resources of the original tool vendor.
	//
	// Examples:
	//   - https://github.com/cloudogu/usermgt
	//   - https://www.atlassian.com/software/jira
	//
	URL string
	// Image contains a reference to the [OCI container] image which packages the dogu application. This field is
	// mandatory. The image must not contain image tags, like the image version or "latest" (use for the field Version
	// for this information instead). The image registry part of this field must point to "registry.cloudogu.com".
	//
	// It is good practice to apply the same name to the image repository as from the Name field in order to enable
	// access strategies as well as to avoid storage conflicts.
	//
	// Examples for official/redmine:
	//   - registry.cloudogu.com/official/redmine
	//
	// [OCI container]: https://opencontainers.org/
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
