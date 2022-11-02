package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/pkg/errors"
)

const (
	CriticalProcessIndicatorName   = "critical_process_running"
	criticalProcessTimeoutDuration = 60
)

type CriticalSystemState struct {
	SystemProcess                  string
	reg                            registry.Registry
	channel                        chan string
	errors                         chan error
	cancel                         context.CancelFunc
	criticalProcessTimeoutDuration time.Duration
}

// NewCriticalSystemState creates a new critical state object
func NewCriticalSystemState(reg registry.Registry, processName string) *CriticalSystemState {
	return &CriticalSystemState{
		reg:                            reg,
		SystemProcess:                  processName,
		channel:                        make(chan string),
		errors:                         make(chan error),
		criticalProcessTimeoutDuration: criticalProcessTimeoutDuration,
	}
}

// getCurrentCriticalSystemState returns the value of the currently set critical system state key
func (css *CriticalSystemState) getCurrentCriticalSystemState() (*CriticalSystemState, error) {
	exists, err := css.reg.GlobalConfig().Exists(CriticalProcessIndicatorName)
	if err != nil {
		return nil, errors.Wrap(err, "could not find critical process key in registry")
	}

	if !exists {
		return &CriticalSystemState{}, nil
	}

	val, err := css.reg.GlobalConfig().Get(CriticalProcessIndicatorName)
	if err != nil {
		return &CriticalSystemState{}, errors.Wrap(err, "error while getting current critical process in registry")
	}

	var current *CriticalSystemState
	err = json.Unmarshal([]byte(val), &current)
	if err != nil {
		return &CriticalSystemState{}, err
	}
	return current, nil
}

// Start sets the critical process key for this process and starts the refresh routine
func (css *CriticalSystemState) Start(ctx context.Context) error {
	existing, err := css.isAnotherProcessRunning()
	if err != nil {
		return errors.Wrapf(err, "could not find out if there is a critical process running")
	}
	if existing {
		return errors.Errorf("there is already a critical process running")
	}

	err = css.setKey()
	if err != nil {
		return errors.Wrap(err, "could not set critical process key")
	}

	interval := (css.criticalProcessTimeoutDuration - 10) * time.Second

	subctx, cancel := context.WithCancel(ctx)
	css.cancel = cancel

	go func(ctx context.Context) {
		for {
			time.Sleep(interval)
			select {
			case <-ctx.Done():
				return
			case s := <-css.channel:
				if s == "wait" {
					css.waitForStart(ctx)
				}
			default:
			}

			err := css.refreshKey()
			if err != nil {
				go func() { css.errors <- err }()
			}
		}
	}(subctx)

	return nil
}

// Pause sends the "wait" signal to the refresh routine
func (css *CriticalSystemState) Pause() error {
	err := css.validateProcessIsRunning()
	if err != nil {
		return err
	}

	css.channel <- "wait"
	return nil
}

// Unpause sends the "start" signal to the refresh routine
func (css *CriticalSystemState) Unpause() error {
	err := css.validateProcessIsRunning()
	if err != nil {
		return err
	}

	css.channel <- "start"
	return nil
}

// Stop removes the critical process key of this process from the registry and stops the refresh routine
func (css *CriticalSystemState) Stop() error {
	err := css.validateProcessIsRunning()
	if err != nil {
		return err
	}

	err = css.removeKey()
	if err != nil {
		return errors.Wrap(err, "error removing critical process key")
	}
	css.cancel()
	return nil
}

// GetErrors returns a slice with any error that was sent to the errors channel
func (css *CriticalSystemState) GetErrors() []error {
	errs := make([]error, 0)
	for {
		select {
		case err := <-css.errors:
			errs = append(errs, err)
		default:
			return errs
		}
	}
}

// waitForStart holds until the string "start" is sent to the channel
func (css *CriticalSystemState) waitForStart(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case w := <-css.channel:
			if w == "start" {
				return
			}
		}
	}
}

// setKey sets the critical process key of this process in registry
func (css *CriticalSystemState) setKey() error {
	data, err := json.Marshal(css)
	if err != nil {
		return err
	}
	err = css.reg.GlobalConfig().SetWithLifetime(CriticalProcessIndicatorName, string(data), int(css.criticalProcessTimeoutDuration))
	if err != nil {
		return err
	}
	return nil
}

// refreshKey refreshes the lifetime of this process' key in the registry. If unset, the key will be set.
func (css *CriticalSystemState) refreshKey() error {
	exists, err := css.reg.GlobalConfig().Exists(CriticalProcessIndicatorName)
	if err != nil {
		return err
	}

	if !exists {
		err = css.setKey()
		if err != nil {
			return err
		}
	}

	err = css.reg.GlobalConfig().Refresh(CriticalProcessIndicatorName, int(css.criticalProcessTimeoutDuration))
	if err != nil {
		return err
	}
	return nil
}

// removeKey removes the critical process key of this process from registry
func (css *CriticalSystemState) removeKey() error {
	keyExists, err := css.reg.GlobalConfig().Exists(CriticalProcessIndicatorName)
	if err != nil {
		return err
	}
	if keyExists {
		err := css.reg.GlobalConfig().Delete(CriticalProcessIndicatorName)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateProcessIsRunning checks if the process is currently running
func (css *CriticalSystemState) validateProcessIsRunning() error {
	if css.channel == nil || css.cancel == nil {
		return errors.Errorf("the critical system state '%s' is not running", css.SystemProcess)
	}

	return nil
}

// isAnotherProcessRunning checks whether there is a process running which is not this process
func (css *CriticalSystemState) isAnotherProcessRunning() (bool, error) {
	current, err := css.getCurrentCriticalSystemState()
	if err != nil {
		return false, errors.Wrapf(err, "could not get the current running critical process")
	}
	return current.SystemProcess == css.SystemProcess, nil
}
