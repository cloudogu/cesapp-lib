package registry

// State handles state related functions of a dogu
type State interface {
	// Get returns the current state Value
	Get() (string, error)
	// Set sets the state of the dogu
	Set(value string) error
	// Remove removes the state of the dogu
	Remove() error
}
