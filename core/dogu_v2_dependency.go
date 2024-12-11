package core

// Contains the different kind of types supported by dogu dependencies
const (
	// DependencyTypeDogu identifies a dogu dependency towards another dogu.
	DependencyTypeDogu = "dogu"
	// DependencyTypeClient identifies a dogu dependency towards a dogu.json-processing client like cesapp or
	// k8s-dogu-operator.
	DependencyTypeClient = "client"
	// DependencyTypePackage identifies a dogu dependency towards an operating system package.
	DependencyTypePackage = "package"
)

// Dependency describes the quality of a dogu dependency towards another entity.
//
// Examples:
//
//	{
//	  "type": "dogu",
//	  "name": "postgresql"
//	}
//
//	{
//	  "name": "postgresql"
//	}
//
//	{
//	  "type": "client",
//	  "name": "k8s-dogu-operator",
//	  "version": ">=0.16.0"
//	}
//
//	{
//	  "type": "package",
//	  "name": "cesappd",
//	  "version": ">=3.2.0"
//	}
type Dependency struct {
	// Type identifies the entity on which the dogu depends. This field is optional.
	// If unset, a value of `dogu` is then assumed.
	//
	// Valid values are one of these: "dogu", "client", "package".
	//
	// A type of "dogu" references another dogu which must be present and running
	// during the dependency check.
	//
	// A type of "client" references the client which processes this dogu's
	// "dogu.json". Several client dependencies of a different client type can be
	// used f. i. to prohibit the processing of a certain client.
	//
	// A type of "package" references a necessary operating system package that must
	// be present during the dependency check.
	//
	// Examples:
	//  - "dogu"
	//  - "client"
	//  - "package"
	//
	Type string `json:"type"`
	// Name identifies the entity selected by Type. This field is mandatory. If the Type selects another dogu, Name
	// must use the simple dogu name (f. e. "postgres"), not the full qualified dogu name (not "official/postgres").
	//
	// Examples:
	//  - "postgresql"
	//  - "k8s-dogu-operator"
	//  - "cesappd"
	//
	Name string `json:"name"`
	// Version selects the version of entity selected by Type. This field is optional. If unset, any version of the
	// selected entity will be accepted during the dependency check.
	//
	// Version accepts different version styles and compare operators.
	//
	// Examples:
	//
	//  - ">=4.1.1-2" - select the entity version greater than or equal to version 4.1.1-2
	//  - "<=1.0.1" - select the entity version less than or equal to version 1.0.1
	//  - "1.2.3.4" - select exactly the version 1.2.3.4
	//
	// With a non-existing version it is possible to negate a dependency.
	//
	// Example:
	//
	//   - "<=0.0.0" - prohibit the selected entity being present
	//
	Version string `json:"version"`
}
