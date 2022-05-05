package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

// combinedEtcdDoguRegistry was added to provide backward compatibility in the dogu registry.
// All writing actions use v1 and v2 registry
// All reading actions use v2 registry, v1 registry if key not found in v2.
type combinedEtcdDoguRegistry struct {
	v1DoguRegistry DoguRegistry
	v2DoguRegistry DoguRegistry
}

func newCombinedEtcdDoguRegistry(client *resilentEtcdClient, pathV1 string, pathV2 string) *combinedEtcdDoguRegistry {
	formatProviderV1 := &core.DoguJsonV1FormatProvider{}
	formatProviderV2 := &core.DoguJsonV2FormatProvider{}

	return &combinedEtcdDoguRegistry{
		v1DoguRegistry: &etcdDoguRegistry{
			pathV1,
			client,
			formatProviderV1,
		},
		v2DoguRegistry: &etcdDoguRegistry{
			pathV2,
			client,
			formatProviderV2,
		},
	}
}

// Enable enables a specific dogu in the backend. The method will create the
// current key to the dogu path
// Enables the dogu in v1 as well as in v2 registry.
func (reg *combinedEtcdDoguRegistry) Enable(dogu *core.Dogu) error {
	err := reg.v1DoguRegistry.Enable(dogu)
	if err != nil {
		return errors.Wrap(err, "could not write to v1 registry")
	}

	err = reg.v2DoguRegistry.Enable(dogu)
	if err != nil {
		return errors.Wrap(err, "could not write to v2 registry")
	}

	return nil
}

// Register registers the dogu at the registry backend
// Registeres the dogu in v1 as well as in v2 registry.
func (reg *combinedEtcdDoguRegistry) Register(dogu *core.Dogu) error {
	err := reg.v1DoguRegistry.Register(dogu)
	if err != nil {
		return errors.Wrap(err, "could not write to v1 registry")
	}

	err = reg.v2DoguRegistry.Register(dogu)
	if err != nil {
		return errors.Wrap(err, "could not write to v2 registry")
	}

	return nil
}

// Get returns a dogu from the registry
// Gets the dogu from v2 registry. If dogu is not installed in v2 registry, v1 registry is used.
func (reg *combinedEtcdDoguRegistry) Get(name string) (*core.Dogu, error) {
	v2dogu, err := reg.v2DoguRegistry.Get(name)
	if err != nil {
		if !IsKeyNotFoundError(err) {
			return nil, errors.Wrap(err, "could not get dogu from v2 registry")
		}

		return reg.v1DoguRegistry.Get(name)
	}

	return v2dogu, nil
}

// GetAll returns all registered dogus
// Collects all dogus in v1 as well as in v2 registry.
func (reg *combinedEtcdDoguRegistry) GetAll() ([]*core.Dogu, error) {
	var allDogus []*core.Dogu

	v2dogus, err := reg.v2DoguRegistry.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get all v2 dogus")
	}

	v1dogus, err := reg.v1DoguRegistry.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not get all v1 dogus")
	}

	allDogus = append(allDogus, v2dogus...)

	for _, v1dogu := range v1dogus {
		if !core.ContainsDoguWithName(v2dogus, v1dogu.Name) {
			allDogus = append(allDogus, v1dogu)
		}
	}

	return allDogus, nil
}

// Unregister removes a dogu from the registry
// Removes in v1 as well as in v2 registry.
func (reg *combinedEtcdDoguRegistry) Unregister(name string) error {
	err := reg.v1DoguRegistry.Unregister(name)
	if err != nil {
		return errors.Wrap(err, "could not unregister v1 dogu")
	}

	err = reg.v2DoguRegistry.Unregister(name)
	if err != nil && !IsKeyNotFoundError(err) {
		return errors.Wrap(err, "could not unregister v2 dogu")
	}

	return nil
}

// IsEnabled returns true if the dogu is installed and enabled.
// Use v1 registry as it contains v1 as well as converted v2 dogus.
func (reg *combinedEtcdDoguRegistry) IsEnabled(name string) (bool, error) {
	return reg.v1DoguRegistry.IsEnabled(name)
}

type etcdDoguRegistry struct {
	path           string
	client         *resilentEtcdClient
	formatProvider core.DoguFormatProvider
}

// Enable enables a specific dogu in the backend. The method will create the
// current key to the dogu path
func (reg *etcdDoguRegistry) Enable(dogu *core.Dogu) error {
	core.GetLogger().Infof("enable dogu %s:%s", dogu.GetSimpleName(), dogu.Version)

	path := reg.path + "/" + dogu.GetSimpleName() + "/current"
	core.GetLogger().Debug("set etcd value at", path)
	_, err := reg.client.Set(path, dogu.Version, nil)
	return err
}

// Register registers the dogu at the registry backend
func (reg *etcdDoguRegistry) Register(dogu *core.Dogu) error {
	// convert dogu to json
	data, err := reg.formatProvider.WriteDoguToString(dogu)
	if err != nil {
		return errors.Wrap(err, "failed to marshal dogu")
	}

	core.GetLogger().Infof("register dogu %s:%s", dogu.GetSimpleName(), dogu.Version)

	// register dogu as json on etcd
	path := reg.path + "/" + dogu.GetSimpleName() + "/" + dogu.Version
	core.GetLogger().Debug("set etcd value on", path)
	_, err = reg.client.Set(path, string(data), nil)
	return err
}

// Get returns a dogu from the registry
func (reg *etcdDoguRegistry) Get(name string) (*core.Dogu, error) {
	value, err := getCurrentValue(reg.client, reg.path+"/"+name)
	if err != nil && IsKeyNotFoundError(err) {
		return nil, err
	}

	dogu, err := reg.formatProvider.ReadDoguFromString(value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dogu json")
	}
	return dogu, nil
}

func getCurrentValue(client *resilentEtcdClient, parent string) (string, error) {
	path := parent + "/current"
	core.GetLogger().Debug("get etcd value from", path)

	version, err := client.Get(path)
	if err != nil {
		core.GetLogger().Debug("could not get current version of", parent)
		return "", err
	}

	path = parent + "/" + version
	core.GetLogger().Debug("get etcd value from", path)

	doguJSON, err := client.Get(path)
	if err != nil {
		core.GetLogger().Warningf("could not read version %s of %s", version, parent)
		return "", err
	}

	return doguJSON, err
}

// GetAll returns all registered dogus
func (reg *etcdDoguRegistry) GetAll() ([]*core.Dogu, error) {
	path := reg.path
	core.GetLogger().Debug("get etcd values from", path)

	children, err := reg.client.GetChildrenPaths(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get dogus from path %s", path)
	}

	var dogus []*core.Dogu
	for _, child := range children {
		value, err := getCurrentValue(reg.client, child)
		if err == nil {
			dogu, err := reg.formatProvider.ReadDoguFromString(value)
			if err != nil {
				core.GetLogger().Warningf("could not unmarshal dogu %s: %v", child, err)
			} else {
				dogus = append(dogus, dogu)
			}
		}
	}

	return dogus, nil
}

// Unregister removes a dogu from the registry
func (reg *etcdDoguRegistry) Unregister(name string) error {
	core.GetLogger().Info("unregister dogu", name)
	path := reg.path + "/" + name + "/current"

	core.GetLogger().Debug("delete etcd key", path)
	return reg.client.Delete(path, nil)
}

// IsEnabled returns true if the dogu is installed and enabled
func (reg *etcdDoguRegistry) IsEnabled(name string) (bool, error) {
	core.GetLogger().Debugf("check if dogu %s is installed", name)
	path := reg.path + "/" + name + "/current"

	return reg.client.Exists(path)
}
