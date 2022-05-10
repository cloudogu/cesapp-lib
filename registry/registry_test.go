package registry_test

import (
	"github.com/cloudogu/cesapp-lib/core"
	"testing"

	"os"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/stretchr/testify/assert"
)

var reg registry.Registry

func init() {
	reg, _ = registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: createTestEndpoint(),
	})
}

func TestNew(t *testing.T) {
	_, err := registry.New(core.Registry{
		Type:      "consul",
		Endpoints: createTestEndpoint(),
	})
	assert.NotNil(t, err)
}

func createTestEndpoint() []string {
	etcd := os.Getenv("ETCD")
	if etcd == "" {
		etcd = "localhost"
	}
	return []string{"http://" + etcd + ":4001"}
}
