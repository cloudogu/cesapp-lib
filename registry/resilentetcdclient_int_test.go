//go:build integration
// +build integration

package registry

import (
	"fmt"
	"github.com/coreos/etcd/client"
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

func TestGetSetDeleteWithRetry(t *testing.T) {
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

func TestSetWithTTL(t *testing.T) {
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

	client, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	_, err = client.Set("/test/one", "1", setOptions)
	require.Nil(t, err)

	exists, err := client.Exists("/test/one")
	require.Nil(t, err)
	require.True(t, exists)

	value, err := client.Get("/test/one")
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Refresh to have the maximum ttl
	value, err = client.Set("/test/one", "", refreshOptions)
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait ttl-2 seconds
	refreshWaitDuration := ttl - 2
	refreshWaitDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", refreshWaitDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)

	value, err = client.Set("/test/one", "", refreshOptions)
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait again ttl-2 seconds and make sure that the value still exists
	fmt.Printf("Waiting %d seconds...\n", refreshWaitDuration)
	time.Sleep(refreshWaitDurationParsed)
	value, err = client.Get("/test/one")
	require.Nil(t, err)
	require.Equal(t, "1", value)

	// Wait until expiration
	expireDuration := ttl + 1
	expireDurationParsed, err := time.ParseDuration(fmt.Sprintf("%ds", expireDuration))
	require.Nil(t, err)
	fmt.Printf("Waiting %d seconds...\n", expireDuration)
	time.Sleep(expireDurationParsed)

	exists, err = client.Exists("/test/one")
	require.Nil(t, err)
	require.False(t, exists)
}

func TestGetChildrenPathsAndRecursiveOperations(t *testing.T) {

	// start http reverse proxy on random port
	server := newFaultyServer()
	defer server.Close()

	client, err := newResilentEtcdClient([]string{server.URL})
	require.Nil(t, err)

	_, err = client.Set("/parent/child0/cchild0", "1", nil)
	require.Nil(t, err)

	_, err = client.Set("/parent/child0/cchild1", "1", nil)
	require.Nil(t, err)

	_, err = client.Set("/parent/child1/cchild0", "1", nil)
	require.Nil(t, err)

	_, err = client.Set("/parent/child2", "1", nil)
	require.Nil(t, err)

	childrenPaths, err := client.GetChildrenPaths("/parent")
	require.Nil(t, err)

	require.Contains(t, childrenPaths, "/parent/child0")
	require.Contains(t, childrenPaths, "/parent/child1")

	children, err := client.GetRecursive("/parent")
	require.Nil(t, err)

	require.Equal(t, "1", children["child0/cchild0"])
	require.Equal(t, "1", children["child0/cchild1"])
	require.Equal(t, "1", children["child1/cchild0"])

	err = client.DeleteRecursive("/parent")
	require.Nil(t, err)

	node, err := client.Get("/parent")
	require.NotNil(t, err)
	require.Equal(t, "", node)
}

func Test_resilentEtcdClient_Get(t *testing.T) {
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
