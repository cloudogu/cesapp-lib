package core

// Registry struct for cesapp
type Registry struct {
	Type      string   `validate:"eq=etcd"`
	Endpoints []string `validate:"required,min=1"`
}
