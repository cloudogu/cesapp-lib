package doguConf

import (
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func Test_createEntryValidator_oneOf(t *testing.T) {
	oneofDescriptor := core.ValidationDescriptor{
		Type:   "ONE_OF",
		Values: []string{""},
	}
	val, err := CreateEntryValidator(oneofDescriptor)
	assert.NoError(t, err)
	assert.IsType(t, &oneOfValidator{}, val)
}

func Test_createEntryValidator_unknown(t *testing.T) {
	oneofDescriptor := core.ValidationDescriptor{
		Type:   "notImplemented",
		Values: []string{""},
	}
	_, err := CreateEntryValidator(oneofDescriptor)
	assert.Error(t, err)
}
