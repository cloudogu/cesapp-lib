package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_newResilientEtcdClient(t *testing.T) {
	t.Run("should return an error if creating the backoff fails", func(t *testing.T) {
		config := core.RetryPolicy{Interval: -2}

		_, err := newResilientEtcdClient(nil, config)

		require.Error(t, err)
	})
}
