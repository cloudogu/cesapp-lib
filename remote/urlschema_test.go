package remote_test

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultURLSchema(t *testing.T) {
	baseURL := "https://ces.com"
	schema := remote.NewDefaultURLSchema(baseURL)
	assert.Equal(t, "https://ces.com/dogus/one/two", schema.Create("one/two"))
	assert.Equal(t, "https://ces.com/dogus/one/two", schema.Get("one/two"))
	assert.Equal(t, "https://ces.com/dogus/one/two/1", schema.GetVersion("one/two", "1"))
	assert.Equal(t, "https://ces.com/dogus/", schema.GetAll())
	assert.Equal(t, "https://ces.com/dogus/one/two/_versions", schema.GetVersionsOf("one/two"))
}

func TestDefaultURLSchemaWithEndingSlash(t *testing.T) {
	baseURL := "https://ces.com/"
	schema := remote.NewDefaultURLSchema(baseURL)
	assert.Equal(t, "https://ces.com/dogus/one/two", schema.Create("one/two"))
	assert.Equal(t, "https://ces.com/dogus/one/two", schema.Get("one/two"))
	assert.Equal(t, "https://ces.com/dogus/one/two/1", schema.GetVersion("one/two", "1"))
	assert.Equal(t, "https://ces.com/dogus/", schema.GetAll())
	assert.Equal(t, "https://ces.com/dogus/one/two/_versions", schema.GetVersionsOf("one/two"))
}

func TestNewIndexURLSchema(t *testing.T) {
	baseURL := "https://ces.com"
	schema := remote.NewIndexURLSchema(baseURL)
	assert.Equal(t, "https://ces.com/one/two/index.json", schema.Create("one/two"))
	assert.Equal(t, "https://ces.com/one/two/index.json", schema.Get("one/two"))
	assert.Equal(t, "https://ces.com/one/two/1/index.json", schema.GetVersion("one/two", "1"))
	assert.Equal(t, "https://ces.com/index.json", schema.GetAll())
	assert.Equal(t, "https://ces.com/one/two/_versions.json", schema.GetVersionsOf("one/two"))
}

func TestIndextURLSchemaWithEndingSlash(t *testing.T) {
	baseURL := "https://ces.com/"
	schema := remote.NewIndexURLSchema(baseURL)
	assert.Equal(t, "https://ces.com/one/two/index.json", schema.Create("one/two"))
	assert.Equal(t, "https://ces.com/one/two/index.json", schema.Get("one/two"))
	assert.Equal(t, "https://ces.com/one/two/1/index.json", schema.GetVersion("one/two", "1"))
	assert.Equal(t, "https://ces.com/index.json", schema.GetAll())
	assert.Equal(t, "https://ces.com/one/two/_versions.json", schema.GetVersionsOf("one/two"))
}

func TestNewURLSchemaByName(t *testing.T) {
	assert.NotNil(t, remote.NewURLSchemaByName("", "https://ces.com"))
	assert.NotNil(t, remote.NewURLSchemaByName("default", "https://ces.com"))
	assert.NotNil(t, remote.NewURLSchemaByName("index", "https://ces.com"))
	assert.Nil(t, remote.NewURLSchemaByName("sorbot", "https://ces.com"))
}
