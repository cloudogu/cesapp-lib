package registry

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/coreos/etcd/client"
	"github.com/prometheus/common/log"
	"sync"
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
		kapi:        client.NewKeysAPI(conn),
		retrier:     r,
		recentIndex: 0,
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
	kapi        client.KeysAPI
	retrier     *retrier.Retrier
	indexMutex  sync.Mutex
	recentIndex uint64
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

		etcd.updateIndexIfNecessary(response.Index)

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

func (etcd *resilentEtcdClient) getMainNode() (*client.Node, error) {
	response, err := etcd.kapi.Get(context.Background(), "/", &client.GetOptions{Recursive: true})
	if err != nil {
		return nil, fmt.Errorf("cannot get main node from etcd: %w", err)
	}

	// `config/_global` is a hidden directory and has to be queried explicit. There is no other way to list hidden dirs.
	for _, node := range response.Node.Nodes {
		if node.Key == "/config" {
			globalConfig, err := etcd.kapi.Get(context.Background(), "/config/_global", &client.GetOptions{Recursive: true})
			if err == nil {
				node.Nodes = append(node.Nodes, globalConfig.Node)
			} else if !strings.Contains(err.Error(), "Key not found (/config/_global)") {
				return nil, fmt.Errorf("cannot get global node from etcd: %w", err)
			} else {
				log.Warn("Key '/config/_global' not found.")
			}
		}
	}

	return response.Node, err
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

// Watch watches for changes of the provided key and sends the event through the channel
func (etcd *resilentEtcdClient) Watch(ctx context.Context, key string, recursive bool, eventChannel chan *client.Response) {
	options := client.WatcherOptions{AfterIndex: etcd.recentIndex, Recursive: recursive}
	watcher := etcd.kapi.Watcher(key, &options)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			etcd.doWatch(ctx, watcher, eventChannel)
		}
	}
}

func (etcd *resilentEtcdClient) doWatch(ctx context.Context, watcher client.Watcher, eventChannel chan *client.Response) {
	resp, err := watcher.Next(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "etcd cluster is unavailable or misconfigured") {
			core.GetLogger().Infof("Cannot reach etcd cluster. Try again in 300 seconds. Error: %v", err)
			etcd.indexMutex.Lock()
			defer etcd.indexMutex.Unlock()
			etcd.recentIndex = 0
			time.Sleep(time.Minute * 5)
			return
		} else {
			core.GetLogger().Infof("Could not get event. Try again in 30 seconds. Error: %v", err)
			etcd.indexMutex.Lock()
			defer etcd.indexMutex.Unlock()
			etcd.recentIndex = 0
			time.Sleep(time.Second * 30)
			return
		}
	}
	eventChannel <- resp
}

// We only update the recent index iff it is 0; which happens only in 2 cases:
// 1. At startup
// 2. In case of an error during watch
// We do this, to not miss any changes made to etcd between
// 1. Startup and starting the watcher
// 2. An error and the restart of the watcher
func (etcd *resilentEtcdClient) updateIndexIfNecessary(index uint64) {
	if etcd.recentIndex == 0 {
		etcd.indexMutex.Lock()
		defer etcd.indexMutex.Unlock()
		if etcd.recentIndex == 0 {
			etcd.recentIndex = index
		}
	}
}
