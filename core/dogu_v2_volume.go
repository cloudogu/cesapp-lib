package core

import "encoding/json"

// Volume defines container volumes that are created during the dogu creation or upgrade.
//
// Examples:
//   - { "Name": "data", "Path":"/usr/share/yourtool/data", "Owner":"1000", "Group":"1000", "NeedsBackup": true}
//   - { "Name": "temp", "Path":"/tmp", "Owner":"1000", "Group":"1000", "NeedsBackup": false}
type Volume struct {
	// Name identifies the volume. This field is mandatory. It must be unique in all volumes of the same dogu.
	//
	// The name syntax must comply with the file system syntax of the respective host operating system and is encouraged
	// to consist of:
	//   - lower case latin characters
	//   - special characters underscore "_", minus "-"
	//   - ciphers 0-9
	//
	// The name must not be "_private" to avoid conflicts with the dogu's private key.
	//
	// Examples:
	//   - tooldata
	//   - tool-data-0
	//
	Name string
	// Path to the directory or file where the volume will be mounted inside the dogu. This field is mandatory.
	//
	// Path may consist of several directory levels, delimited by a forward slash "/". Path must comply with the file
	// system syntax of the container operating system.
	//
	// The path must not match `/private` to avoid conflicts with the dogu's private key.
	//
	// Examples:
	//   - /usr/share/yourtool
	//   - /tmp
	//   - /usr/share/license.txt
	//
	Path string
	// Owner contains the numeric Unix UID of the user owning this volume. This field is optional.
	//
	// For security reasons it is strongly recommended to set the Owner of the volume to an unprivileged user. Please
	// note that container image must be then built in a way that the container process may own the path either by
	// user or group ownership.
	//
	// The owner syntax must consist of ciphers (0-9) only.
	//
	// Examples:
	//   - "1000" - an unprivileged user
	//   - "0" - the root user
	//
	Owner string
	// Group contains the numeric Unix GID of the group owning this volume. This field is optional.
	//
	// For security reasons it is strongly recommended to set the Group of the volume to an unprivileged Group. Please
	// note that container image must be then built in a way that the container process may own the path either by
	// user or group ownership.
	//
	// The Group syntax must consist of ciphers (0-9) only.
	//
	// Examples:
	//   - "1000" - an unprivileged group
	//   - "0" - the root group
	//
	Group string
	// NeedsBackup controls whether the Cloudogu EcoSystem backup facility backs up the whole the volume or not. This
	// field is optional. If unset, a value of `false` will be assumed.
	NeedsBackup bool
	// Clients contains a list of client-specific (t. i., the client that interprets the dogu.json) configurations for
	// the volume. This field is optional.
	//
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

// VolumeClient adds additional information for clients to create volumes.
//
// Example:
//
//	{
//	  "Name": "k8s-dogu-operator",
//	  "Params": {
//	    "Type": "configmap",
//	    "Content": {
//	      "Name": "k8s-ces-menu-json"
//	    }
//	  }
//	}
type VolumeClient struct {
	// Name identifies the client responsible to process this volume definition. This field is mandatory.
	//
	// Examples:
	//   - cesapp
	//   - k8s-dogu-operator
	Name string
	// Params contains generic data only interpretable by the client. This field is mandatory.
	Params interface{}
}
