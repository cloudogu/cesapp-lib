package tasks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/cesapp-lib/registry/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestCanStartAndStopCriticalProcess(t *testing.T) {
	// given
	emptyState := `{}`
	state := "{\"SystemProcess\":\"testprocess\"}"
	reg := mocks.NewRegistry(t)
	globalConfig := mocks.NewConfigurationContext(t)

	reg.On("GlobalConfig").Return(globalConfig)
	globalConfig.On("Exists", "critical_process_running").Return(true, nil)
	globalConfig.On("Get", "critical_process_running").Once().Return(emptyState, nil)
	globalConfig.On("Get", "critical_process_running").Once().Return(state, nil)
	globalConfig.On("Get", "critical_process_running").Once().Return(emptyState, nil)
	globalConfig.On("SetWithLifetime", "critical_process_running", state, 11).Return(nil)
	globalConfig.On("Refresh", "critical_process_running", 11).Return(nil)
	globalConfig.On("Delete", "critical_process_running").Return(nil)

	ctx := context.Background()

	process := NewCriticalSystemState(reg, "testprocess")
	process.criticalProcessTimeoutDuration = 11 // The interval has an offset of 10

	// when
	err := process.Stop()
	require.Error(t, err, "is not running")

	err = process.Start(ctx)
	require.Nil(t, err)

	current, err := process.getCurrentCriticalSystemState()
	require.Nil(t, err)

	// then
	require.Equal(t, "testprocess", current.SystemProcess)

	// when
	err = process.Pause()
	require.Nil(t, err)

	err = process.Unpause()
	require.Nil(t, err)

	err = process.Stop()
	require.NoError(t, err)
	current, err = process.getCurrentCriticalSystemState()

	// then
	require.Nil(t, err)
	require.Equal(t, "", current.SystemProcess)
}

func TestCanOnlySttOneCriticalProcess(t *testing.T) {
	// given
	ctx := context.Background()
	emptyState := `{}`
	state := "{\"SystemProcess\":\"testprocess\"}"
	reg := mocks.NewRegistry(t)
	globalConfig := mocks.NewConfigurationContext(t)

	reg.On("GlobalConfig").Return(globalConfig)
	globalConfig.On("Exists", "critical_process_running").Return(true, nil)
	globalConfig.On("Get", "critical_process_running").Once().Return(emptyState, nil)
	globalConfig.On("Get", "critical_process_running").Once().Return(state, nil)
	globalConfig.On("SetWithLifetime", "critical_process_running", state, 60).Return(nil)

	process1 := NewCriticalSystemState(reg, "testprocess")
	process2 := NewCriticalSystemState(reg, "testprocess")

	// when
	err := process1.Start(ctx)
	// then
	require.Nil(t, err)

	// when
	err = process2.Start(ctx)
	// then
	require.NotNil(t, err)
}

func TestCannotPauseOrStopPressWithoutStart(t *testing.T) {
	reg := mocks.NewRegistry(t)

	process := NewCriticalSystemState(reg, "testprocess")
	require.NotNil(t, process)

	err := process.Stop()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "the critical system state 'testprocess' is not running")

	err = process.Pause()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "the critical system state 'testprocess' is not running")
}

func TestCriticalProcessFailsOnExistsError(t *testing.T) {
	r := mocks.CreateMockRegistry(nil)
	reg := r.Registry
	regs := r.SubRegistries
	mocks.OnExists(regs["_global"], CriticalProcessIndicatorName, false, errors.New("test"))

	css := NewCriticalSystemState(reg, "test")
	err := css.Start(context.Background())

	assert.NotNil(t, err)
	require.Contains(t, err.Error(), "test")
	regs["_global"].AssertExpectations(t)
}

func TestCriticalProcessFailsOnGetError(t *testing.T) {
	r := mocks.CreateMockRegistry(nil)
	reg := r.Registry
	regs := r.SubRegistries
	mocks.OnExists(regs["_global"], CriticalProcessIndicatorName, true, nil)
	mocks.OnGet(regs["_global"], mock.Anything, mock.Anything, errors.New("test"))

	css := NewCriticalSystemState(reg, "test")
	err := css.Start(context.Background())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "test")
	regs["_global"].AssertExpectations(t)
}

func TestCriticalProcessFailsOnInvalidJson(t *testing.T) {
	r := mocks.CreateMockRegistry(nil)
	reg := r.Registry
	regs := r.SubRegistries
	mocks.OnExists(regs["_global"], CriticalProcessIndicatorName, true, nil)
	mocks.OnGet(regs["_global"], CriticalProcessIndicatorName, "", nil)

	css := NewCriticalSystemState(reg, "test")
	err := css.Start(context.Background())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unexpected end of JSON input")
}

func TestCriticalProcessFailsOnSet(t *testing.T) {
	r := mocks.CreateMockRegistry(nil)
	reg := r.Registry
	regs := r.SubRegistries
	mocks.OnExists(regs["_global"], CriticalProcessIndicatorName, false, nil)
	mocks.OnSetWithLifetime(regs["_global"], CriticalProcessIndicatorName, mocks.Anything, mocks.AnyLifetime, errors.New("test"))

	css := NewCriticalSystemState(reg, "test")
	err := css.Start(context.Background())

	assert.NotNil(t, err)
	require.Contains(t, err.Error(), "test")
	regs["_global"].AssertExpectations(t)
}

func TestUnpauseFailsWhenNotStarted(t *testing.T) {
	reg := mocks.NewRegistry(t)
	css := NewCriticalSystemState(reg, "test")
	err := css.Pause()

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not running")
}

func TestStopFailsOnExistsError(t *testing.T) {
	first := true
	r := mocks.CreateMockRegistry(nil)
	reg := r.Registry
	regs := r.SubRegistries
	regs["_global"].On("Exists", mock.MatchedBy(func(key string) bool {
		if first {
			first = false
			return true
		}
		return false
	})).Return(false, nil)

	mocks.OnSetWithLifetime(regs["_global"], CriticalProcessIndicatorName, mocks.Anything, mocks.AnyLifetime, nil)

	css := NewCriticalSystemState(reg, "test")
	err := css.Start(context.Background())
	require.Nil(t, err)
	mocks.OnExists(regs["_global"], CriticalProcessIndicatorName, false, errors.New("testerror"))

	err = css.Stop()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "test")
	regs["_global"].AssertExpectations(t)
}
