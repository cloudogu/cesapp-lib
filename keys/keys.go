package keys

import (
	"github.com/cloudogu/cesapp-lib/core"
)

// KeyProviders contains the list of implemented key providers and is dynamically filled by the providers in keys package
var KeyProviders []string

var log = core.GetLogger()

var (
	providers = make(map[string]*KeyProvider)
)

func register(name string, provider *KeyProvider) {
	providers[name] = provider
	KeyProviders = append(KeyProviders, name)
}
