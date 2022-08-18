//go:build integration
// +build integration

package registry

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/retrier"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"net/url"

	"math/rand"

	"os"

	"github.com/stretchr/testify/require"
)

func newFaultyServer() *httptest.Server {
	// initialize rand
	rand.Seed(time.Now().UnixNano())

	errCounter := 0

	reverseProxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// inject randomly errors, but max 4 times
		if errCounter < 4 && rand.Float32() < 0.50 {
			// increase error counter
			errCounter++

			// write error 500
			w.WriteHeader(500)
			return
		}
		// reset error counter
		errCounter = 0

		// create etcd address, for local execution and on ci
		etcd := os.Getenv("ETCD")
		if etcd == "" {
			etcd = "localhost"
		}

		// proxy request to local etcd
		director := func(req *http.Request) {
			path := "http://" + etcd + ":4001" + r.URL.RequestURI()
			req.URL, _ = url.Parse(path)
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	return httptest.NewServer(reverseProxyHandler)
}

func newServer() *httptest.Server {
	// initialize rand
	rand.Seed(time.Now().UnixNano())

	reverseProxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create etcd address, for local execution and on ci
		etcd := os.Getenv("ETCD")
		if etcd == "" {
			etcd = "localhost"
		}

		// proxy request to local etcd
		director := func(req *http.Request) {
			path := "http://" + etcd + ":4001" + r.URL.RequestURI()
			req.URL, _ = url.Parse(path)
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	return httptest.NewServer(reverseProxyHandler)
}

func Test_getMainNode_inttest(t *testing.T) {
	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	cl, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	defer func() {
		_ = cl.DeleteRecursive("/dir_test")
	}()

	_, err = cl.Set("/dir_test/key1/subkey1", "val1", nil)
	require.Nil(t, err)

	_, err = cl.Set("/dir_test/key1/subkey2", "val2", nil)
	require.Nil(t, err)

	_, err = cl.Set("/dir_test/key2", "val3", nil)
	require.Nil(t, err)

	node, err := cl.getMainNode()
	require.NoError(t, err)

	found := false
	for _, node := range node.Nodes {
		if node.Key == "/dir_test" {
			found = true
			assert.Len(t, node.Nodes, 2)
			key1Node := node.Nodes[0]
			key2Node := node.Nodes[1]
			// slice is randomly sorted. Make sure to test the correct nodes
			if node.Nodes[1].Key == "/dir_test/key1" {
				key1Node = node.Nodes[1]
				key2Node = node.Nodes[0]
			}
			assert.Equal(t, "/dir_test/key1", key1Node.Key)
			assert.Equal(t, "/dir_test/key2", key2Node.Key)
			assert.Len(t, key1Node.Nodes, 2)
			subkey1Node := key1Node.Nodes[0]
			subkey2Node := key1Node.Nodes[1]
			// slice is randomly sorted. Make sure to test the correct nodes
			if subkey2Node.Key == "/dir_test/key1/subkey1" {
				subkey1Node = subkey1Node.Nodes[1]
				subkey2Node = subkey2Node.Nodes[0]
			}
			assert.Equal(t, "/dir_test/key1/subkey1", subkey1Node.Key)
			assert.Equal(t, "/dir_test/key1/subkey2", subkey2Node.Key)
		}
	}

	assert.True(t, found)
}

func findChildByKey(node *client.Node, key string) *client.Node {
	for _, n := range node.Nodes {
		if n.Key == key {
			return n
		}
	}
	return nil
}

func TestGetSetDeleteWithRetry_inttest(t *testing.T) {
	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	cl, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	_, err = cl.Set("/test/one", "1", nil)
	require.Nil(t, err)

	_, err = cl.Set("/mydir", "", &client.SetOptions{
		Dir:       true,
		PrevExist: client.PrevIgnore,
	})
	require.Nil(t, err)

	_, err = cl.Set("/mydir/key/one", "val", nil)
	require.Nil(t, err)

	exists, err := cl.Exists("/test/one")
	require.Nil(t, err)
	require.True(t, exists)

	value, err := cl.Get("/test/one")
	require.Nil(t, err)

	require.Equal(t, "1", value)

	err = cl.Delete("/test/one", nil)
	require.Nil(t, err)

	exists, err = cl.Exists("/test/one")
	require.Nil(t, err)
	require.False(t, exists)

	exists, err = cl.Exists("/mydir")
	require.Nil(t, err)
	require.True(t, exists)

	exists, err = cl.Exists("/mydir/key/one")
	require.Nil(t, err)
	require.True(t, exists)

	err = cl.DeleteRecursive("/mydir")

	exists, err = cl.Exists("/mydir/key/one")
	require.Nil(t, err)
	require.False(t, exists)
}

func TestWatch_inttest(t *testing.T) {

	// start http reverse proxy on random port
	server := newServer()
	defer server.Close()

	cl, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	myResponseChannel := make(chan *client.Response)
	ctx, clearFunc := context.WithTimeout(context.Background(), time.Second*10)

	// schedule a key change event in 2 seconds
	changeKeyTimer := time.NewTimer(time.Second * 2)
	go func() {
		<-changeKeyTimer.C
		_, err = cl.Set("/mywatchkey", "myvalue", &client.SetOptions{})
		require.NoError(t, err)
		changeKeyTimer.Stop()
	}()

	go cl.Watch(ctx, "/mywatchkey", false, myResponseChannel)

	for {
		select {
		case response := <-myResponseChannel:
			t.Logf("Watch result: %+v", response)
			assert.Equal(t, "set", response.Action)
			assert.Equal(t, "/mywatchkey", response.Node.Key)
			clearFunc()
			return
		case <-ctx.Done():
			t.Logf("Test failed as it should detect the change in the /mywatchkey")
			clearFunc()
			t.Fail()
			return
		}
	}
}

func TestSetWithTTL_inttest(t *testing.T) {
	ttl := 5
	ttlParsed, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	require.Nil(t, err)

	setOptions := &client.SetOptions{
		TTL: ttlParsed,
	}

	refreshOptions := &client.SetOptions{
		TTL:     ttlParsed,
		Refresh: true,
	}

	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	etcdClient, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	_, err = etcdClient.Set("/test/one", "1", setOptions)
	require.Nil(t, err)

	exists, err := etcdClient.Exists("/test/one")
	require.Nil(t, err)
	require.True(t, exists)

	value, err := etcdClient.Get("/test/one")
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Refresh to have the maximum ttl
	value, err = etcdClient.Set("/test/one", "", refreshOptions)
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait ttl-2 seconds
	refreshWaitDuration := ttl - 2
	refreshWaitDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", refreshWaitDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)

	value, err = etcdClient.Set("/test/one", "", refreshOptions)
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait again ttl-2 seconds and make sure that the value still exists
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)
	value, err = etcdClient.Get("/test/one")
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait until expiration
	expireDuration := ttl + 1
	expireDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", expireDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", expireDuration)
	time.Sleep(expireDurationParsed)

	exists, err = etcdClient.Exists("/test/one")
	require.Nil(t, err)
	require.False(t, exists)
}

func TestGetChildrenPathsAndRecursiveOperations_inttest(t *testing.T) {

	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	etcdClient, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	_, err = etcdClient.Set("/parent/child0/cchild0", "1", nil)
	require.Nil(t, err)

	_, err = etcdClient.Set("/parent/child0/cchild1", "1", nil)
	require.Nil(t, err)

	_, err = etcdClient.Set("/parent/child1/cchild0", "1", nil)
	require.Nil(t, err)

	_, err = etcdClient.Set("/parent/child2", "1", nil)
	require.Nil(t, err)

	childrenPaths, err := etcdClient.GetChildrenPaths("/parent")
	require.Nil(t, err)

	require.Contains(t, childrenPaths, "/parent/child0")
	require.Contains(t, childrenPaths, "/parent/child1")

	children, err := etcdClient.GetRecursive("/parent")
	require.Nil(t, err)

	require.Equal(t, "1", children["child0/cchild0"])
	require.Equal(t, "1", children["child0/cchild1"])
	require.Equal(t, "1", children["child1/cchild0"])

	err = etcdClient.DeleteRecursive("/parent")
	require.Nil(t, err)

	node, err := etcdClient.Get("/parent")
	require.NotNil(t, err)
	require.Equal(t, "", node)
}

func Test_resilentEtcdClient_Get_inttest(t *testing.T) {
	t.Run("should return error which can be tested with IsKeyNotFoundError ", func(t *testing.T) {
		mockedRetrier := retrier.New(
			retrier.ConstantBackoff(1, time.Millisecond),
			&etcdClassifier{},
		)

		mockedKeysAPI := new(mockKeysAPI)
		clientErr := client.Error{Code: client.ErrorCodeKeyNotFound}
		clientResponse := &client.Response{}
		mockedKeysAPI.On("Get", mock.Anything, "/config/theKey", mock.Anything).Return(clientResponse, clientErr)
		sut := resilentEtcdClient{kapi: mockedKeysAPI, retrier: mockedRetrier}

		actual, err := sut.Get("/config/theKey")

		require.Error(t, err)
		require.True(t, IsKeyNotFoundError(err))
		require.Empty(t, actual)
		mockedKeysAPI.AssertExpectations(t)
	})
}

func Test_resilentEtcdClient_Watch_inttest(t *testing.T) {
	t.Run("successfull terminated watch with context timeout", func(t *testing.T) {
		// given
		mockedRetrier := retrier.New(
			retrier.ConstantBackoff(1, time.Millisecond),
			&etcdClassifier{},
		)
		clientResponse := &client.Response{}
		watcherMock := new(mockWatcher)
		watcherMock.On("Next", mock.Anything).Return(clientResponse, nil)
		mockedKeysAPI := new(mockKeysAPI)
		mockedKeysAPI.On("Watcher", "/key", mock.Anything).Return(watcherMock)
		underTest := resilentEtcdClient{kapi: mockedKeysAPI, retrier: mockedRetrier}
		eventChannel := make(chan *client.Response)

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)

		go func() {
			for range eventChannel {
				assert.True(t, true)
			}
		}()

		// when
		underTest.Watch(ctx, "/key", false, eventChannel)
		cancelFunc()
	})
}

type mockWatcher struct {
	mock.Mock
}

func (m *mockWatcher) Next(ctx context.Context) (*client.Response, error) {
	args := m.Called(ctx)
	return args.Get(0).(*client.Response), args.Error(1)
}

type mockKeysAPI struct {
	mock.Mock
}

func (m *mockKeysAPI) Get(ctx context.Context, key string, opts *client.GetOptions) (*client.Response, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) Set(ctx context.Context, key, value string, opts *client.SetOptions) (*client.Response, error) {
	args := m.Called(ctx, key, value, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) Delete(ctx context.Context, key string, opts *client.DeleteOptions) (*client.Response, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) Create(ctx context.Context, key, value string) (*client.Response, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) CreateInOrder(ctx context.Context, dir, value string, opts *client.CreateInOrderOptions) (*client.Response, error) {
	args := m.Called(ctx, dir, value, opts)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) Update(ctx context.Context, key, value string) (*client.Response, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(*client.Response), args.Error(1)
}

func (m *mockKeysAPI) Watcher(key string, opts *client.WatcherOptions) client.Watcher {
	args := m.Called(key, opts)
	return args.Get(0).(client.Watcher)
}
