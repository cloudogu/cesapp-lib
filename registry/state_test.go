//go:build integration
// +build integration

package registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	t.Run("default case", func(t *testing.T) {
		state := reg.State("unit-test-1")

		err := state.Set("installing")
		assert.Nil(t, err)

		value, err := state.Get()
		assert.Nil(t, err)
		assert.Equal(t, "installing", value)

		err = state.Set("ready")
		assert.Nil(t, err)

		value, err = state.Get()
		assert.Nil(t, err)
		assert.Equal(t, "ready", value)

		err = state.Remove()
		assert.Nil(t, err)

		value, err = state.Get()
		assert.Nil(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("remove key which is not set throws no error", func(t *testing.T) {
		state := reg.State("unit-test-2")

		value, err := state.Get()
		assert.Nil(t, err)
		assert.Equal(t, "", value)

		err = state.Remove()

		value, err = state.Get()
		assert.Nil(t, err)
		assert.Equal(t, "", value)
	})
}
