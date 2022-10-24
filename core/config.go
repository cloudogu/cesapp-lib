package core

import "strconv"

type RetryPolicy struct {
	Interval      int64 `json:"interval"`
	MaxRetryCount int   `json:"maxRetryCount"`
}

// Registry contains Cloudogu EcoSystem registration details.
type Registry struct {
	Type        string      `validate:"eq=etcd"`
	Endpoints   []string    `validate:"required,min=1"`
	RetryPolicy RetryPolicy `json:"retryPolicy,omitempty"`
}

// Remote contains dogu registry configuration details.
type Remote struct {
	Endpoint               string `validate:"url"`
	AuthenticationEndpoint string `validate:"omitempty,url"`
	URLSchema              string `json:"urlSchema,omitempty" validate:"omitempty,oneof=default index"`
	CacheDir               string `validate:"required"`
	ProxySettings          ProxySettings
	AnonymousAccess        bool        `json:",omitempty"`
	Insecure               bool        `json:",omitempty"`
	RetryPolicy            RetryPolicy `json:"retryPolicy,omitempty"`
}

// ProxySettings contains the settings for http proxy
type ProxySettings struct {
	Enabled  bool
	Server   string `json:",omitempty"`
	Port     int    `json:",omitempty"`
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

// CreateURL creates a proxy http url
func (proxy ProxySettings) CreateURL() string {
	return "http://" + proxy.Server + ":" + strconv.Itoa(proxy.Port)
}

// Credentials for a remote system
type Credentials struct {
	Username string
	Password string
}
