package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

type etcdState struct {
	path   string
	client *resilentEtcdClient
}

// Get returns the current state value
func (es *etcdState) Get() (string, error) {
	core.GetLogger().Debug("try to get state key at", es.path)
	keyExists, err := es.client.Exists(es.path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check if key %s exists", es.path)
	}
	if !keyExists {
		core.GetLogger().Debugf("key %s not found. returning empty state", es.path)
		return "", nil
	}
	state, err := es.client.Get(es.path)
	if err != nil {
		return "", errors.Wrapf(err, "could not get state value at %s", es.path)
	}
	return state, nil
}

// Set sets the state of the dogu
func (es *etcdState) Set(value string) error {
	core.GetLogger().Debug("try to set state key", es.path)
	_, err := es.client.Set(es.path, value, nil)
	if err != nil {
		return errors.Wrapf(err, "could not set state value at %s", es.path)
	}
	return nil
}

// Set sets the state of the dogu
func (es *etcdState) Remove() error {
	core.GetLogger().Debug("try to remove state key", es.path)

	exists, err := es.client.Exists(es.path)
	if err != nil {
		return errors.Wrapf(err, "failed to check if state key exists")
	}

	if !exists {
		return nil
	}

	err = es.client.Delete(es.path, nil)
	if err != nil {
		return errors.Wrap(err, "could not remove dogu state")
	}

	return nil
}
