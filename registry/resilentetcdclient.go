package registry

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/coreos/etcd/client"
	"time"

	"context"

	"strings"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/pkg/errors"
)

// newResilentEtcdClient is build up on the the kapi of etcd and adds constant retries for every failed request.
func newResilentEtcdClient(endpoints []string) (*resilentEtcdClient, error) {
	core.GetLogger().Debug("create etcd client for endpoints", endpoints)

	r := retrier.New(
		retrier.ExponentialBackoff(5, 100*time.Millisecond),
		&etcdClassifier{},
	)

	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
	}

	conn, err := client.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create etcd client")
	}

	return &resilentEtcdClient{
		kapi:    client.NewKeysAPI(conn),
		retrier: r,
	}, nil
}

type etcdClassifier struct{}

// Classify returns succeeds if the error is nil or an etcd not found error, in all other cases the classifier will return retry.
func (classifier *etcdClassifier) Classify(err error) retrier.Action {
	if err == nil || IsKeyNotFoundError(err) {
		return retrier.Succeed
	}
	return retrier.Retry
}

type resilentEtcdClient struct {
	kapi    client.KeysAPI
	retrier *retrier.Retrier
}

// Exists returns true if the key exists
func (etcd *resilentEtcdClient) Exists(key string) (bool, error) {
	var exists bool
	err := etcd.retrier.Run(func() error {
		core.GetLogger().Debugf("check if key %s exists", key)
		_, err := etcd.kapi.Get(context.Background(), key, nil)
		if err == nil {
			exists = true
			return nil
		}
		if IsKeyNotFoundError(err) {
			exists = false
			return nil
		}

		return err
	})

	if err != nil {
		return false, errors.Wrapf(err, "failed to read key %s", key)
	}

	return exists, nil
}

// Get returns the value of the given node, otherwise it returns an error. If the given key cannot be found a
// KeyNotFoundError is returned.
func (etcd *resilentEtcdClient) Get(key string) (string, error) {
	var result string
	err := etcd.retrier.Run(func() error {
		core.GetLogger().Debugf("read key %s", key)
		response, err := etcd.kapi.Get(context.Background(), key, nil)
		if err != nil {
			return err
		}

		result = response.Node.Value
		return nil
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to read key %s", key)
	}

	return result, nil
}

// GetRecursive returns a map of key value pairs below the given key
func (etcd *resilentEtcdClient) GetRecursive(key string) (map[string]string, error) {
	var result map[string]string
	err := etcd.retrier.Run(func() error {
		core.GetLogger().Debugf("read key %s recursive", key)
		response, err := etcd.kapi.Get(context.Background(), key, &client.GetOptions{Recursive: true})
		if err != nil {
			return err
		}

		result = map[string]string{}
		etcd.addNodeDirValuesToMap(result, response.Node, "")
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to read key %s recursive", key)
	}

	return result, nil
}

func (etcd *resilentEtcdClient) addNodeDirValuesToMap(keyValuePairs map[string]string, node *client.Node, parent string) {
	for _, child := range node.Nodes {
		childKey := etcd.createKey(parent, child.Key)
		if child.Dir {
			etcd.addNodeDirValuesToMap(keyValuePairs, child, childKey)
		} else {
			keyValuePairs[childKey] = child.Value
		}
	}
}

func (etcd *resilentEtcdClient) createKey(parent string, nodeKey string) string {
	key := parent
	if parent != "" {
		key += "/"
	}

	indexOfLastSlash := strings.LastIndex(nodeKey, "/")
	if indexOfLastSlash > 0 {
		key += nodeKey[indexOfLastSlash+1:]
	} else {
		key += nodeKey
	}

	return key
}

// GetChildrenPaths returns an array of all children keys of the given key
func (etcd *resilentEtcdClient) GetChildrenPaths(key string) ([]string, error) {
	children := []string{}
	err := etcd.retrier.Run(func() error {
		core.GetLogger().Debugf("read children paths from %s", key)

		resp, err := etcd.kapi.Get(context.Background(), key, nil)
		if err != nil && IsKeyNotFoundError(err) {
			return nil
		} else if err != nil {
			return err
		}

		for _, child := range resp.Node.Nodes {
			children = append(children, child.Key)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to read children from %s", key)
	}

	return children, nil
}

// Set sets the key to the given value
func (etcd *resilentEtcdClient) Set(key string, value string, options *client.SetOptions) (string, error) {
	var result string
	err := etcd.retrier.Run(func() error {
		response, err := etcd.kapi.Set(context.Background(), key, value, options)
		if err != nil {
			return err
		}

		result = response.Node.Value
		return nil
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to set key %s", key)
	}

	return result, nil
}

// Delete deletes the given key or directory
func (etcd *resilentEtcdClient) Delete(key string, options *client.DeleteOptions) error {
	err := etcd.retrier.Run(func() error {
		core.GetLogger().Debugf("delete key %s", key)
		_, err := etcd.kapi.Delete(context.Background(), key, options)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete key %s", key)
	}

	return nil
}

// DeleteRecursive deletes the given key and all its children
func (etcd *resilentEtcdClient) DeleteRecursive(key string) error {
	core.GetLogger().Debugf("delete key %s recursive", key)

	return etcd.Delete(key, &client.DeleteOptions{
		Recursive: true,
	})
}
