package archive2_test

import (
	"github.com/cloudogu/cesapp-lib/archive2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCanAddBytesToArchive(t *testing.T) {
	manager := archive2.NewManager()

	bytes := manager.GetContent()
	assert.True(t, len(bytes) == 0)

	err := manager.AddContentAsFile("test", "myfile")
	require.NoError(t, err)

	err = manager.Close()
	require.NoError(t, err)

	bytes = manager.GetContent()
	assert.True(t, len(bytes) > 0)
}
