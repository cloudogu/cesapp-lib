package remote_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"net/url"
	"strconv"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearCache() {
	_ = os.RemoveAll("/tmp/ces/cache/remote_test")
}

func TestAnonymousOnAnonymousServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			data, _ := core.WriteDoguToString(&core.Dogu{Name: "Test", Version: "3.0"})
			w.Write([]byte(data))
			w.WriteHeader(200)
		}
	}))
	defer ts.Close()
	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{
		AnonymousAccess: true,
	})

	version, err := testRemote.GetVersion("Test", "3.0")
	assert.NotNil(t, version)
	assert.Nil(t, err)
}

func TestAuthorizedOnAnonymousServer(t *testing.T) {
	ts := createAnonymousHttpServer()
	defer ts.Close()
	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{
		CacheDir: "Test",
	})
	testRemote.SetUseCache(false)

	version, err := testRemote.GetVersion("Test", "3.0")
	assert.Nil(t, version)
	assert.NotNil(t, err)
}

func TestAnonymousOnAuthorizedServer(t *testing.T) {
	ts := createAuthorrizedHttpServer()
	defer ts.Close()
	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{
		CacheDir:        "Test",
		AnonymousAccess: true,
	})
	testRemote.SetUseCache(false)

	version, err := testRemote.GetVersion("Test", "3.0")
	assert.NotNil(t, version)
	assert.Nil(t, err)
}

func TestAuthorizedOnAuthorizedServer(t *testing.T) {
	ts := createAuthorrizedHttpServer()
	defer ts.Close()
	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{
		CacheDir: "Test",
	})
	testRemote.SetUseCache(false)

	version, err := testRemote.GetVersion("Test", "3.0")
	assert.NotNil(t, version)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.Equal(t, "trillian", username)
		assert.Equal(t, "secret", password)
		assert.True(t, ok)

		w.WriteHeader(200)

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)

		dogu, _, err := core.ReadDoguFromString(string(body))
		assert.Nil(t, err)

		assert.Equal(t, "Test", dogu.Name)
		assert.Equal(t, "1.0.0", dogu.Version)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	err := testRemote.Create(&core.Dogu{
		Name:    "Test",
		Version: "1.0.0",
	})
	assert.Nil(t, err)

	clearCache()
}

func TestCreateWithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	err := testRemote.Create(&core.Dogu{
		Name:    "Test",
		Version: "1.0.0",
	})
	assert.NotNil(t, err)

	clearCache()
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/Test"))

		data, err := core.WriteDoguToString(&core.Dogu{Name: "Test", Version: "3.0"})
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	dogu, err := testRemote.Get("Test")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)

	assert.Equal(t, "Test", dogu.Name)
	assert.Equal(t, "3.0", dogu.Version)

	clearCache()
}

func TestGetWithRetry(t *testing.T) {
	counter := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter > 0 {
			data, err := core.WriteDoguToString(&core.Dogu{Name: "official/Hansolo", Version: "3.0"})
			assert.Nil(t, err)
			w.Write([]byte(data))
			w.WriteHeader(200)
		} else {
			counter++
			w.WriteHeader(500)
		}
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	dogu, err := testRemote.Get("official/Hansolo")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)

	assert.Equal(t, "official/Hansolo", dogu.Name)
	assert.Equal(t, "3.0", dogu.Version)

	assert.Equal(t, 1, counter)

	clearCache()
}

func TestGetVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/Test/3.0"))

		data, err := core.WriteDoguToString(&core.Dogu{Name: "Test", Version: "3.0"})
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	dogu, err := testRemote.GetVersion("Test", "3.0")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)

	assert.Equal(t, "Test", dogu.Name)
	assert.Equal(t, "3.0", dogu.Version)

	clearCache()
}

func TestIsDoingAnonymousAccessOnceWithRetry(t *testing.T) {
	var accessCounter = 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessCounter++
		if (accessCounter % 2) != 0 {
			//First access and third access will be here (because we do one retry)
			assert.Empty(t, r.Header.Get("Authorization"))
		} else {
			assert.NotEmpty(t, r.Header.Get("Authorization"))
		}
	}))
	defer ts.Close()

	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{
		AnonymousAccess: true,
	})
	_, _ = testRemote.GetVersion("Test", "3.0")
	assert.Equal(t, 4, accessCounter)

	// Try again with anonymous docker access also activated. Should not affect the result
	testRemote = createRemoteWithConfiguration(t, ts, &core.Remote{
		AnonymousAccess: true,
	})
	_, _ = testRemote.GetVersion("Test", "3.0")
	assert.Equal(t, 8, accessCounter)
}

func TestIsDoingNoAnynymousAcces(t *testing.T) {
	var accessCounter = 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEmpty(t, r.Header.Get("Authorization"))
		accessCounter++
	}))
	defer ts.Close()

	testRemote := createRemoteWithConfiguration(t, ts, &core.Remote{})
	_, _ = testRemote.GetVersion("Test", "3.0")
	assert.Equal(t, 2, accessCounter)

	// Try again with anonymous docker access also activated. Should not affect the result
	testRemote = createRemoteWithConfiguration(t, ts, &core.Remote{})
	_, _ = testRemote.GetVersion("Test", "3.0")
	assert.Equal(t, 4, accessCounter)
}

func TestGetDoNotRetryWithFailedAuthentication(t *testing.T) {
	counter := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(403)
		}
		counter++
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	_, err := testRemote.GetVersion("Test", "3.0")
	assert.NotNil(t, err)
	assert.Equal(t, 1, counter)

	_, err = testRemote.GetVersion("Test", "3.0")
	assert.NotNil(t, err)
	assert.Equal(t, 2, counter)

	clearCache()
}

func TestGetAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dogus := []*core.Dogu{
			{Name: "a", Version: "1.0"},
			{Name: "b", Version: "2.0"},
		}
		data, err := core.WriteDogusToString(dogus)
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	dogus, err := testRemote.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, dogus)
	assert.Equal(t, 2, len(dogus))
	assert.Equal(t, "a", dogus[0].Name)
	assert.Equal(t, "b", dogus[1].Name)

	clearCache()
}

func TestGetVersions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []string{
			"2.89.2-1",
			"2.73.1-2",
			"2.73.1-1",
		}
		data, err := json.Marshal(versions)
		assert.Nil(t, err)
		w.Write(data)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	testRemote := createRemote(t, ts)
	versions, err := testRemote.GetVersionsOf("official/jenkins")
	assert.Nil(t, err)
	assert.NotNil(t, versions)

	assert.Equal(t, "2.89.2-1", versions[0].Raw)
	assert.Equal(t, "2.73.1-2", versions[1].Raw)
	assert.Equal(t, "2.73.1-1", versions[2].Raw)

	clearCache()
}

func TestGetVersionsCached(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []string{
			"2.89.2-1",
			"2.73.1-2",
			"2.73.1-1",
		}
		data, err := json.Marshal(versions)
		assert.Nil(t, err)
		w.Write(data)
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)
	versions, err := testRemote.GetVersionsOf("official/jenkins")
	assert.Nil(t, err)
	assert.NotNil(t, versions)

	assert.Equal(t, "2.89.2-1", versions[0].Raw)
	assert.Equal(t, "2.73.1-2", versions[1].Raw)
	assert.Equal(t, "2.73.1-1", versions[2].Raw)

	ts.Close()

	versions, err = testRemote.GetVersionsOf("official/jenkins")
	assert.Nil(t, err)
	assert.NotNil(t, versions)

	assert.Equal(t, "2.89.2-1", versions[0].Raw)
	assert.Equal(t, "2.73.1-2", versions[1].Raw)
	assert.Equal(t, "2.73.1-1", versions[2].Raw)

	clearCache()
}

func TestGetVersionsNotCached(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []string{
			"2.89.2-1",
			"2.73.1-2",
			"2.73.1-1",
		}
		data, err := json.Marshal(versions)
		assert.Nil(t, err)
		w.Write(data)
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)
	testRemote.SetUseCache(false)
	versions, err := testRemote.GetVersionsOf("official/jenkins")
	assert.Nil(t, err)
	assert.NotNil(t, versions)

	assert.Equal(t, "2.89.2-1", versions[0].Raw)
	assert.Equal(t, "2.73.1-2", versions[1].Raw)
	assert.Equal(t, "2.73.1-1", versions[2].Raw)

	ts.Close()

	versions, err = testRemote.GetVersionsOf("official/jenkins")
	assert.NotNil(t, err)
}

func TestGetAllCached(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dogus := []*core.Dogu{
			&core.Dogu{Name: "a", Version: "1.0"},
			&core.Dogu{Name: "b", Version: "2.0"},
		}
		data, err := core.WriteDogusToString(dogus)
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)

	dogus, err := testRemote.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, dogus)
	assert.Equal(t, 2, len(dogus))
	assert.Equal(t, "a", dogus[0].Name)
	assert.Equal(t, "b", dogus[1].Name)

	ts.Close()

	dogus, err = testRemote.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, dogus)
	assert.Equal(t, 2, len(dogus))
	assert.Equal(t, "a", dogus[0].Name)
	assert.Equal(t, "b", dogus[1].Name)

	clearCache()
}

func TestGetAllNotCached(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dogus := []*core.Dogu{
			&core.Dogu{Name: "a", Version: "1.0"},
			&core.Dogu{Name: "b", Version: "2.0"},
		}
		data, err := core.WriteDogusToString(dogus)
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)

	testRemote.SetUseCache(false)
	dogus, err := testRemote.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, dogus)
	assert.Equal(t, 2, len(dogus))
	assert.Equal(t, "a", dogus[0].Name)
	assert.Equal(t, "b", dogus[1].Name)

	ts.Close()

	dogus, err = testRemote.GetAll()
	assert.NotNil(t, err)
}

func TestGetCached(t *testing.T) {
	expectedDogu := core.Dogu{Name: "a", Version: "1.0"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data, err := core.WriteDoguToString(&expectedDogu)
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)

	dogu, err := testRemote.Get("a")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)
	assert.Equal(t, expectedDogu, *dogu)

	ts.Close()

	dogu, err = testRemote.Get("a")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)
	assert.Equal(t, expectedDogu, *dogu)

	clearCache()
}

func TestGetNotCached(t *testing.T) {
	expectedDogu := core.Dogu{Name: "a", Version: "1.0"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data, err := core.WriteDoguToString(&expectedDogu)
		assert.Nil(t, err)
		w.Write([]byte(data))
		w.WriteHeader(200)
	}))

	testRemote := createRemote(t, ts)

	dogu, err := testRemote.Get("a")
	assert.Nil(t, err)
	assert.NotNil(t, dogu)
	assert.Equal(t, expectedDogu, *dogu)

	ts.Close()
	testRemote.SetUseCache(false)
	dogu, err = testRemote.Get("a")
	assert.NotNil(t, err)

	clearCache()
}
func TestGetAllWithTlsServer(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := &core.Dogu{Name: "a", Version: "1.0"}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)

		data, err := core.WriteDoguToString(a)
		require.Nil(t, err)
		_, err = w.Write([]byte(data))
		require.Nil(t, err)
	}))
	defer ts.Close()

	t.Run("insecure is true", func(t *testing.T) {
		testRemote, err := remote.New(
			&core.Remote{
				Endpoint:        ts.URL,
				CacheDir:        "/tmp/ces/cache/remote_test",
				Insecure:        true,
				AnonymousAccess: true,
			},
			&core.Credentials{
				Username: "trillian",
				Password: "secret",
			},
		)
		assert.Nil(t, err)
		testRemote.SetUseCache(false)

		sample, err := testRemote.Get("sample")
		require.Nil(t, err)
		require.Equal(t, "a", sample.Name)
	})
	t.Run("insecure is false", func(t *testing.T) {

		testRemote, err := remote.New(
			&core.Remote{
				Endpoint:               ts.URL,
				AuthenticationEndpoint: "",
				URLSchema:              "",
				CacheDir:               "/tmp/ces/cache/remote_test",
				ProxySettings:          core.ProxySettings{},
				AnonymousAccess:        true,
				Insecure:               false,
			},
			&core.Credentials{
				Username: "trillian",
				Password: "secret",
			},
		)
		assert.Nil(t, err)
		testRemote.SetUseCache(false)

		_, err = testRemote.Get("sample")
		require.NotNil(t, err)
	})
}

func TestGetAllWithProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := &core.Dogu{Name: "a", Version: "1.0"}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)

		data, err := core.WriteDoguToString(a)
		require.Nil(t, err)
		_, err = w.Write([]byte(data))
		require.Nil(t, err)
	}))

	defer ts.Close()

	parsedURL, err := url.Parse(ts.URL)
	require.Nil(t, err)

	port, err := strconv.Atoi(parsedURL.Port())
	require.Nil(t, err)

	testRemote, err := remote.New(
		&core.Remote{
			Endpoint: "http://pangalaktischer-donnergurgler.org",
			CacheDir: "/tmp/ces/cache/remote_test",
			ProxySettings: core.ProxySettings{
				Enabled:  true,
				Server:   parsedURL.Hostname(),
				Port:     port,
				Username: "trillian",
				Password: "mcmillian",
			},
		},
		&core.Credentials{
			Username: "trillian",
			Password: "secret",
		},
	)

	sample, err := testRemote.Get("sample")
	require.Nil(t, err)
	require.Equal(t, "a", sample.Name)
}

func TestGetAllWithIndexURLSchema(t *testing.T) {
	var url string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL.RequestURI()
		w.WriteHeader(200)
		w.Write([]byte("[]"))
	}))
	defer ts.Close()

	testRemote := createRemoteWithURLScheme(t, ts, "index")
	_, err := testRemote.GetAll()
	assert.Nil(t, err)

	assert.True(t, strings.HasSuffix(url, "/index.json"))

	clearCache()
}

func createRemote(t *testing.T, ts *httptest.Server) remote.Registry {
	return createRemoteWithURLScheme(t, ts, "")
}

func createRemoteWithConfiguration(t *testing.T, ts *httptest.Server, remoteConf *core.Remote) remote.Registry {
	t.Helper()
	remoteConf.Endpoint = ts.URL
	rem, err := remote.New(
		remoteConf,
		&core.Credentials{
			Username: "trillian",
			Password: "secret",
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, rem)
	return rem
}

func createRemoteWithURLScheme(t *testing.T, ts *httptest.Server, urlScheme string) remote.Registry {
	testRemote, err := remote.New(
		&core.Remote{
			Endpoint:  ts.URL,
			CacheDir:  "/tmp/ces/cache/remote_test",
			URLSchema: urlScheme,
		},
		&core.Credentials{
			Username: "trillian",
			Password: "secret",
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, testRemote)

	return testRemote
}

func createAnonymousHttpServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			dogu := &core.Dogu{Name: "Test", Version: "3.0"}
			data, _ := core.WriteDoguToString(dogu)
			w.Write([]byte(data))
			w.WriteHeader(200)
		}
	}))
	return ts
}

func createAuthorrizedHttpServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			dogu := &core.Dogu{Name: "Test", Version: "3.0"}
			data, _ := core.WriteDoguToString(dogu)
			w.Write([]byte(data))
			w.WriteHeader(200)
		}
	}))
	return ts
}
