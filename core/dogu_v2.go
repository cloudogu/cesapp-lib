package core

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

// Names for ExposedCommands correspond with actual dogu descriptor instance values. Do not change because these come
// with side effects.
const (
	// ExposedCommandServiceAccountCreate identifies a name of a core.ExposedCommand which produces a service account
	// for other dogus.
	ExposedCommandServiceAccountCreate = "service-account-create"

	// ExposedCommandServiceAccountRemove identifies a name of a core.ExposedCommand which removes a service accounts.
	ExposedCommandServiceAccountRemove = "service-account-remove"

	// ExposedCommandPreUpgrade identifies a name of a core.ExposedCommand which executes before a dogu upgrade.
	ExposedCommandPreUpgrade = "pre-upgrade"

	// ExposedCommandPostUpgrade identifies a name of a core.ExposedCommand which executes after a dogu upgrade.
	ExposedCommandPostUpgrade = "post-upgrade"

	// ExposedCommandUpgradeNotification identifies a name of a core.ExposedCommand which informs a user of important,
	// upcoming changes (possibly with a call for action) that should be acknowledged by the administrator.
	ExposedCommandUpgradeNotification = "upgrade-notification"

	// ExposedCommandBackupConsumer identifies a name of a core.ExposedCommand which backs up a service account
	// consumer dogu's data during the back-up.
	ExposedCommandBackupConsumer = "backup-consumer"

	// ExposedCommandPreBackup
	//
	// Deprecated: This field is no longer used. There is no substitute.
	ExposedCommandPreBackup = "pre-backup"

	// ExposedCommandPostBackup
	//
	// Deprecated: This field is no longer used. There is no substitute.
	ExposedCommandPostBackup = "post-backup"
)

// HealthCheck provide readiness and health checks for the dogu container.
type HealthCheck struct {
	// Type specifies the nature of the health check. This field is mandatory. It can be either tcp, http or state.
	//
	// For Type "tcp" the given Port needs to be open to be healthy.
	//
	// For Type "http" the service needs to return a status code between >= 200 and < 300 to be healthy.
	// Port and Path are used to reach the service.
	//
	// For Type "state" the /state/<dogu> key in the Cloudogu EcoSystem registry gets checked.
	// This key is written by the installation process.
	// A 'ready' value in etcd means that the dogu is healthy.
	Type string
	// State contains the expected health check state of Type state. This field is optional, even for health checks of
	// the type "state". The default is "ready".
	State string
	// Port is the tcp-port for health checks of Type tcp and http. This field is mandatory for the health check of type
	// "tcp" and "http"
	Port int
	// Path is the Http-Path for health checks of Type "http". This field is mandatory for the health check of type
	// "http". The default is '/health'.
	Path string
	// Parameters may contain key-value pairs for check specific parameters.
	//
	// Deprecated: is not in use.
	Parameters map[string]string
}

// ExposedPort struct is used to define ports which are exported to the host.
//
// Example:
//
//	{ "Type": "tcp", "Container": "2222", "Host":"2222" }
//	{ "Container": "2222", "Host":"2222" }
type ExposedPort struct {
	// Type contains the protocol type over which the container communicates (f. i. 'tcp'). This field is optional (the
	// value of `tcp` is then assumed).
	//
	// Example:
	//   - tcp
	//   - udp
	//   - sctp
	//
	Type string
	// Container contains the mapped port on side of the container. This field is mandatory. Usual port range
	// limitations apply.
	//
	// Examples:
	//   - 80
	//   - 8080
	//   - 65535
	//
	Container int
	// Host contains the mapped port on side of the host. This field is mandatory. Usual port range limitations apply.
	//
	// Examples:
	//   - 80
	//   - 8080
	//   - 65535
	//
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
	// Name identifies the exposed command. This field is mandatory. It must be unique in all commands of the same dogu.
	//
	// The name must consist of:
	//   - lower case latin characters
	//   - special characters underscore "_", minus "-"
	//   - ciphers 0-9
	//
	// Examples:
	//   - service-account-create
	//   - plugin-delete
	//
	Name string
	// Description describes the exposed command in a few words. This field is optional.
	//
	// Example:
	//  - Creates a new service account
	//  - Delete a plugin
	//
	Description string
	// Command identifies the script to be executed for this command. This field is mandatory.
	//
	// The name syntax must comply with the file system syntax of the respective host operating system and is encouraged
	// to consist of:
	//   - lower case latin characters
	//   - special characters underscore "_", minus "-"
	//   - ciphers 0-9
	//
	// Examples:
	//   - /resources/create-sa.sh
	//   - /resources/deletePlugin.sh
	//
	Command string
}

// EnvironmentVariable struct represents custom parameters that can change
// the behaviour of a dogu build process
type EnvironmentVariable struct {
	Key   string
	Value string
}

// String returns a string representation of this [EnvironmentVariable]
func (env EnvironmentVariable) String() string {
	// Formatting of an EnvironmentVariable "ENV1=VALUE1"
	return env.Key + "=" + env.Value
}

// ServiceAccount struct can be used to get access to another dogu.
//
// Example:
//
//	{
//	 "Type": "k8s-dogu-operator",
//	 "Kind": "k8s"
//	}
type ServiceAccount struct {
	// Type contains the name of the service on which the account should be created. This field is mandatory.
	//
	// Example:
	//   - postgresql
	//   - your-dogu
	//
	Type string
	// Params contains additional arguments necessary for the service account creation. The optionality of this field
	// depends on the desired service account producer dogu. Please consult the service account producer dogu in
	// question.
	Params []string
	// Kind defines the kind of service on which the account should be created, e.g. `dogu` or `k8s`. This field is
	// optional. If empty, a default value of `dogu` should be assumed.
	//
	// Reading this property and creating a corresponding service account is up to the client.
	//
	Kind string `json:"Kind,omitempty"`
}

// Properties describes generic properties of the dogu which are evaluated by a client like cesapp or k8s-dogu-operator.
//
// Example:
//   - { "key": "value" }
type Properties map[string]string

// Dogu describes properties of a containerized application for the Cloudogu EcoSystem. Besides the meta information and
// the [OCI container image], Dogu describes all necessities for automatic container instantiation, f. i. volumes,
// dependencies towards other dogus, and much more.
//
// Example:
//
//	{
//	 "Name": "official/newdogu",
//	 "Version": "1.0.0-1",
//	 "DisplayName": "My new Dogu",
//	 "Description": "Newdogu is a test application",
//	 "Category": "Development Apps",
//	 "Tags": ["warp"],
//	 "Url": "https://www.company.com/newdogu",
//	 "Image": "registry.cloudogu.com/namespace/newdogu",
//	 "Dependencies": [
//	   {
//	     "type":"dogu",
//	     "name":"nginx"
//	   }
//	 ],
//	 "Volumes": [
//	   {
//	     "Name": "temp",
//	     "Path":"/tmp",
//	     "Owner":"1000",
//	     "Group":"1000",
//	     "NeedsBackup": false
//	   }
//	 ],
//	 "HealthChecks": [
//	   {
//	     "Type": "tcp",
//	     "Port": 8080
//	   }
//	 ]
//	}
//
// [OCI container image]: https://opencontainers.org/
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
	//  - 1.7.8-1
	//  - 2.2.0-4
	//
	// Recommended example in the Dockerfile:
	//  LABEL maintainer="hello@cloudogu.com" \
	//  NAME="official/nginx" \
	//  VERSION="1.23.2-1"
	//
	Version string
	// PublishedAt is the date and time when the dogu was created.
	//
	// This field does not need to be filled in by the developer. The dogu.cloudogu.com service automatically replaces
	// the content with the current time when the new dogu is created.
	//
	// Examples:
	//   - 2024-10-16T07:49:34.738Z
	//   - 2019-05-03T13:31:48.612Z
	//
	PublishedAt time.Time
	// DisplayName is the name of the dogu which is used in UI frontends to represent the dogu. This field is mandatory.
	//
	// Usages:
	// In the setup of the ecosystem the display name of the dogu is used to select it for installation.
	//
	// The display name is used in the warp menu to let the user navigate to the web ui of the dogu.
	//
	// Another usage is the textual output of tools like the cesapp or the k8s-dogu-operator where the name is used
	// in commands like list upgradeable dogus.
	//
	// The display name may consist of:
	//   - lower and upper case latin characters where the first is upper case
	//   - any special characters
	//   - ciphers 0-9
	//
	// DisplayName is encouraged to consist of less than 30 characters because it is displayed in the Cloudogu EcoSystem
	// warp menu.
	//
	// Examples:
	//  - Jenkins CI
	//  - Backup & Restore
	//  - SCM-Manager
	//  - Smeagol
	//
	DisplayName string
	// Description contains a short explanation, what the dogu does. This field is mandatory.
	//
	// It is used in the setup of the ecosystem in the dogu selection.
	// Therefore, the description should give an uninformed user a brief hint what the dogu is
	// and maybe the function the dogu fulfills.
	//
	// The description may consist of:
	//   - lower and upper case latin characters where the first is upper case
	//   - any special characters
	//   - ciphers 0-9
	//
	// Description is encouraged to consist a readable sentence which explains shortly the dogu's main topic.
	//
	// Examples:
	//  - Jenkins Continuous Integration Server
	//  - MySQL - Relational database
	//  - The Nexus Repository is like the local warehouse where all the parts and finished goods used in your
	//     software supply chain are stored and distributed.
	//
	Description string
	// Category organizes the dogus in three categories. This field is mandatory.
	//
	// These categories are fixed and must be either:
	//
	//  - `Development Apps` - For regular dogus which should be used by a regular user of the ecosystem,
	//  - `Administration Apps` - For dogus which should be used by a user with administration rights
	//  - `Base` - For dogus which are important for the overall system.
	//
	// The categories "Development Apps" and "Administration Apps" are represented in the warp menu to order the dogus.
	//
	// Example dogus for each category:
	//  - `Development Apps`: Redmine, SCM-Manager, Jenkins
	//  - `Administration Apps`: Backup & Restore, User Management
	//  - `Base`: Nginx, Registrator, OpenLDAP
	//
	Category string
	// Tags contains a list of one-word-tags which are in connection with the dogu. This field is optional.
	//
	// If the dogu should be displayed in the warp menu the tag "warp" is necessary.
	// Other tags won't be automatically processed.
	//
	// Examples for e.g. Jenkins:
	//  {"warp", "build", "ci", "cd"}
	//
	Tags []string
	// Logo represents a URI to a web picture depicting the dogu tool. This field is optional.
	//
	// Deprecated: The Cloudogu EcoSystem does not facilitate the logo URI. It is a candidate for removal.
	// Other options of representing a tool or application can be:
	//   - embed the logo in the dogu's Git repository (if public)
	//   - provide the logo in to dogu UI (if the dogu provides one)
	//
	Logo string
	// URL links the website to the original tool vendor. This field is optional.
	//
	// The Cloudogu EcoSystem does not facilitate this information.
	// Anyhow, in a public dogu repository in which a dogu vendor re-packages
	// a third party application the URL may point users to resources of the original tool vendor.
	//
	// Examples:
	//   - https://github.com/cloudogu/usermgt
	//   - https://www.atlassian.com/software/jira
	//
	URL string
	// Image links to the [OCI container] image which packages the dogu application. This field is mandatory.
	//
	// The image must not contain image tags, like the image version or "latest" (use the field Version
	// for this information instead). The image registry part of this field must point to `registry.cloudogu.com`.
	//
	// It is good practice to apply the same name to the image repository as from the Name field in order to enable
	// access strategies as well as to avoid storage conflicts.
	//
	// Examples for official/redmine:
	//   - registry.cloudogu.com/official/redmine
	//
	// [OCI container]: https://opencontainers.org/
	//
	Image string
	// ExposedPorts contains a list of [ExposedPort]. This field is optional.
	//
	// Exposed ports provide a way to route traffic from outside the Cloudogu EcoSystem towards the dogu. Dogus that
	// provide a GUI must also expose a port that accepts requests and returns responses. Non-GUI dogus can also expose
	// ports if the dogu provides services to a consumer (f. i. when the dogu wants to provide an API for a CLI tool).
	//
	// Examples:
	//
	//   [ { "Type": "tcp", "Container": "2222", "Host":"2222" } ]
	//
	ExposedPorts []ExposedPort
	// ExposedCommands defines actions of type [ExposedCommand] which can be executed in different phases of
	// the dogu lifecycle automatically triggered by a dogu client like the cesapp
	// or the k8s-dogu-operator or manually from an administrative user. This
	// field is optional.
	//
	// Usually, commands which are automatically triggered by a dogu-client are
	// those in upgrade or service account processes. These commands are:
	// pre-upgrade, post-upgrade, upgrade-notification, service-account-create,
	// service-account-remove
	//
	// pre-upgrade:
	// This command will be executed during an early stage of an
	// upgrade process from a dogu. A dogu client will mount the pre-upgrade script
	// from the new dogu version in the container of the old still running dogu and
	// executes it. It is mainly used to prepare data migrations  (e.g. export a
	// database). Core of the script should be the comparison between the version
	// and to determine if a migration is needed. For this, the dogu client will
	// call the script with the old version as the first and the new version as the
	// second parameter. In addition, it is recommended to set states like
	// "upgrading" or "pre-upgrade done" in the script. This can be very useful
	// because you can use it in the post-upgrade or the regular startup for a
	// waiting functionality.
	//
	// post-upgrade:
	// This command will be executed after a regular dogu upgrade.
	// Like in the pre-upgrade the old dogu version is passed as the first and the
	// new dogu version is passed as the second parameter. They should be used to
	// determine if an action is needed. Keep in mind that in this time the regular
	// startup script of your new container will be executed as well. Use a state in
	// the etcd to handle a wait functionality in the startup. If the post-upgrade
	// ends reset this state and start the regular container.
	//
	// upgrade-notification:
	// If the upgrade process works, e.g. with really sensitive data an upgrade
	// notification script can be implemented to inform the user about the risk of
	// this process and e.g. give a backup instruction. Before an upgrade the dogu
	// client executes the notification script, prints all upgrade steps and asks
	// the user if he wants to proceed.
	//
	// service-account-create:
	// The service-account-create command is used in the
	// installation process of a dogu. A service account in the Cloudogu EcoSystem
	// represents credentials to authorize another dogu, like a database or an
	// authentication server. The service-account-create command has to be
	// implemented in the service account producer dogu which produces the
	// credentials. If a service account consuming dogu will be installed and
	// requires a service account (e.g. for postgresql) the dogu client will call
	// the service-account-create script in the postgresql dogu with the service
	// name as the first parameters and custom ones as additional parameters. See
	// core.ServiceAccounts for how to define custom parameters. With this
	// information the script should create a service account and save it (e.g. USER
	// table in an underlying database or maybe encrypted in the etcd). After that,
	// the credentials must be printed to console so that the dogu client saves the
	// credentials for the dogu which requested the service account. It is also
	// important that these outputs are the only ones from the script, otherwise the
	// dogu client will use them as credentials.
	//
	// Example output:
	//  echo "database: ${DATABASE}"
	//  echo "username: ${USER}"
	//  echo "password: ${PASSWORD}"
	//
	// service-account-delete:
	// If a dogu will be deleted the used service account
	// should be deleted too. Because of this, the service account producing dogu
	// (like postgresql) must implement a service-account-delete command. This
	// script will be executed in the deletion process of a dogu. The dogu
	// processing client passes the service account name as the only parameter. In
	// contrast to the service-account-create command the script can output logs
	// because the dogu client won't use any of this data.
	//
	// Example commands for service accounts:
	//  "ExposedCommands": [
	//     {
	//       "Name": "service-account-create",
	//       "Description": "Creates a new service account",
	//       "Command": "/create-sa.sh"
	//     },
	//     {
	//       "Name": "service-account-remove",
	//       "Description": "Removes a service account",
	//       "Command": "/remove-sa.sh"
	//     }
	//   ],
	//
	// Custom commands:
	// Moreover, a dogu can specify commands which are not executed
	// automatically in the dogu lifecycle (e.g. upgrade). For example, a dogu like
	// Redmine can specify a command to delete or install a plugin. At runtime an
	// administrator can call this command with the cesapp like:
	//  cesapp command redmine plugin-install scrum-plugin
	//
	// Example:
	//  "ExposedCommands": [
	//     {
	//       "Name": "plugin-install",
	//       "Description": "Installs a plugin",
	//       "Command": "/installPlugin.sh"
	//     }
	//   ],
	//
	ExposedCommands []ExposedCommand
	// Volumes contains a list of [Volume] which defines [OCI container volumes]. This field is optional.
	//
	// Volumes are created during the dogu creation or upgrade.
	//
	// All dynamic data of a Dogu should be stored inside volumes. This holds multiple advantages:
	//   - Data inside volumes is not lost when the Dogu is upgraded.
	//   - Data inside volumes can be backed up and restored. See [needsbackup] flag.
	//   - The access to data stored inside volumes is much faster than to data stored inside the Dogu container.
	//
	// Examples:
	//   [ { "Name": "temp", "Path":"/tmp", "Owner":"1000", "Group":"1000", "NeedsBackup": false} ]
	//
	// [OCI container volumes]: https://opencontainers.org/
	//
	Volumes []Volume
	// HealthCheck defines a single way to check the dogu health for observability. This field is optional.
	//
	// Deprecated: use HealthChecks instead
	HealthCheck HealthCheck
	// HealthChecks defines multiple ways to check the dogu health for observability. This field is optional.
	//
	// HealthChecks are used in various use cases:
	//  - to show the `dogu is starting` page until the dogu is healthy at startup
	//  - for monitoring via `cesapp healthy <dogu-name>`
	//  - for monitoring via the admin dogu
	//  - to avoid backing up unhealthy dogus
	//
	// There are different types of health checks:
	//  - state, via the `/state/<dogu>` etcd key set by ces-setup
	//  - tcp, to check for an open port
	//  - http, to check a status code of a http response
	//
	// They must be executable without authentication required.
	//
	// Example:
	//	"HealthChecks": [
	//	  {"Type": "state"}
	//	  {"Type": "tcp",  "Port": 8080},
	//	  {"Type": "http", "Port": 8080, "Path": "/my/health/path"},
	// 	]
	HealthChecks []HealthCheck
	// ServiceAccounts contains a list of [ServiceAccount]. This field is optional.
	//
	// A ServiceAccount protects sensitive data used for another dogu. Service account consumers are distinguished from
	// service account producers. To produce a service account a dogu has to implement service-account-create and service-account-delete
	// scripts in core.ExposedCommands. To consume a service account a dogu should define a request in ServiceAccounts.
	//
	// In the installation process a dogu client recognizes a service account request and executes
	// the service-account-create script with this information from the service account producer.
	// The credentials will be then stored in the Cloudogu EcoSystem registry.
	//
	// Example paths:
	//  - config/redmine/sa-postgresql/username
	//  - config/redmine/sa-postgresql/password
	//  - config/redmine/sa-postgresql/database
	//
	// The dogu itself can then use the service account by read and decrypt the credentials with `doguctl` in the startup script.
	//
	// Examples:
	//    "ServiceAccounts": [
	//     {
	//       "Type": "ldap",
	//       "Params": [
	//         "rw"
	//       ]
	//     },
	//     {
	//      "Type": "k8s-dogu-operator",
	//      "Kind": "k8s"
	//     }
	//   ],
	//
	ServiceAccounts []ServiceAccount
	// Privileged indicates whether the Docker socket should be mounted into the container file system. This field is
	// optional. The default value is `false`.
	//
	// For security reasons, it is highly recommended to leave Privileged set to false since almost no dogu should
	// gain introspective container insights.
	//
	// Example:
	//   - false
	//
	// Deprecated: This feature will be removed in the future because no dogu should have the privilege of container
	// meta-insights for obvious security reasons. Also, this field is subject of a misnomer because mounting a
	// container socket has nothing to do with privileged execution.
	Privileged bool
	// Security defines security policies for the dogu. This field is optional.
	//
	// The Cloudogu Ecosystem may not support restricting capabilities in all environments, e.g. in docker.
	// This feature was added for the kubernetes platform via [pod security context].
	//
	// Example:
	//
	//   {
	//     "Security": {
	//       "Capabilities": {
	//         "Drop": ["All"],
	//         "Add": ["NetBindService", "Kill"]
	//       },
	//       "RunAsNonRoot": true,
	//       "ReadOnlyRootFileSystem": true
	//     }
	//   }
	//
	// [pod security context]: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	Security Security `json:"Security,omitempty"`
	// Configuration contains a list of [ConfigurationField]. This field is optional.
	//
	// It describes generic properties of the dogu in the Cloudogu EcoSystem registry.
	//
	// Examples:
	//   {
	//     "Name": "plugin_center_url",
	//     "Description": "URL of SCM-Manager Plugin Center",
	//     "Optional": true
	//   }
	//
	//   {
	//     "Name": "logging/root",
	//     "Description": "Set the root log level to one of ERROR, WARN, INFO, DEBUG or TRACE. Default is INFO",
	//     "Optional": true,
	//     "Default": "INFO",
	//     "Validation": {
	//       "Type": "ONE_OF",
	//       "Values": [
	//         "WARN",
	//         "ERROR",
	//         "INFO",
	//         "DEBUG",
	//         "TRACE"
	//       ]
	//     }
	//   }
	Configuration []ConfigurationField
	// Properties is a `map[string]string` of Properties. This field is optional.
	// It describes generic properties of the dogu which are evaluated by a client like cesapp or k8s-dogu-operator.
	//
	// Example:
	//   {
	//     "key1": "value1",
	//     "key2": "value2"
	//   }
	Properties Properties
	// EnvironmentVariables contains a list of [EnvironmentVariable] that should
	// be set in the container. This field is optional.
	//
	// Example:
	//   [
	//     {"Key": "my_key", "Value": "my_value"}
	//   ]
	EnvironmentVariables []EnvironmentVariable
	// Dependencies contains a list of [Dependency]. This field is optional.
	//
	// This field defines dependencies that must be fulfilled during the dependency check. If the dependency
	// cannot be fulfilled during the check an error will be thrown and the processing will be stopped.
	//
	// Examples:
	//  [
	//    {
	//      "type": "dogu",
	//      "name": "postgresql"
	//    },
	//    {
	//      "type": "client",
	//      "name": "cesapp",
	//      "version": ">=6.0.0"
	//    },
	//    {
	//      "type": "package",
	//      "name": "cesappd",
	//      "version": ">=3.2.0"
	//    }
	//  ]
	//
	Dependencies []Dependency
	// OptionalDependencies contains a list of [Dependency]. This field is optional.
	//
	// In contrast to core.Dependencies, OptionalDependencies allows to define
	// dependencies that may be fulfilled if they are existent. There is no negative
	// impact during the dependency check if an optional dependency does not exist.
	// But if the optional dependency is present and the dependency cannot be
	// fulfilled during the dependency check then an error will be thrown and the
	// processing will be stopped.
	//
	// Examples:
	//  [
	//    {
	//      "type": "dogu",
	//      "name": "ldap"
	//    }
	//  ]
	//
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

// GetEnvironmentVariables returns a slice of [EnvironmentVariable] formatted as array of key value objects
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

// EffectiveCapabilities returns the accumulated list of those Capabilities that are retained after applying the
// capabilities to be added and removed. DefaultCapabilities build the baseline for added or dropped capabilities.
func (d *Dogu) EffectiveCapabilities() []Capability {
	effectiveCaps := make(map[Capability]int)

	for _, defaultCap := range DefaultCapabilities {
		// note this works since go 1.22 because iteration variables now can be used as unshared variable
		effectiveCaps[defaultCap] = 0 // we only use the map to check for keys, values don't matter
	}

	for _, dropCap := range d.Security.Capabilities.Drop {
		if dropCap == All {
			effectiveCaps = make(map[Capability]int)
			break
		}
		delete(effectiveCaps, dropCap)
	}

	for _, addCap := range d.Security.Capabilities.Add {
		if addCap == All {
			// do a fast exit here because alternatives of slice-to-map conversion would be cumbersome
			return slices.Clone(AllCapabilities)
		}
		effectiveCaps[addCap] = 0 // we only use the map to check for keys, values don't matter
	}

	actualCaps := maps.Keys(effectiveCaps)
	return slices.Collect(actualCaps)
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
	dogu := DoguV1{
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

// ValidateSecurity checks the dogu's Security section for configuration errors.
func (d *Dogu) ValidateSecurity() error {
	var errs error

	for _, value := range d.Security.Capabilities.Add {
		if !slices.Contains(AllCapabilities, value) {
			err := fmt.Errorf("%s is not a valid capability to be added", value)
			errs = errors.Join(errs, err)
		}
	}

	for _, value := range d.Security.Capabilities.Drop {
		if !slices.Contains(AllCapabilities, value) {
			err := fmt.Errorf("%s is not a valid capability to be dropped", value)
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return fmt.Errorf("dogu %s:%s contains an invalid security field: %w", d.Name, d.Version, errs)
	}

	return nil
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
