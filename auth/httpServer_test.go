package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newHttpServer(t *testing.T) {
	actual := NewHttpServer("test")

	require.NotNil(t, actual)
}
