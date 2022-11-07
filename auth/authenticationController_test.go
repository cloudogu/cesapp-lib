package auth

import (
	"github.com/cloudogu/cesapp-lib/credentials"
	"github.com/cloudogu/cesapp-lib/credentials/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticationController(t *testing.T) {
	// given
	authConfig := AuthenticationConfig{}
	httpServer := NewHttpServer("")
	store, _ := credentials.NewStore(os.TempDir())

	// when
	result := NewHttpAuthenticationController(authConfig, httpServer, nil, store)

	// then
	require.NotNil(t, result)
}

func Test_realAuthenticationController_authenticationHandler(t *testing.T) {
	t.Run("redirect to auth endpoint with no credentials", func(t *testing.T) {
		// given
		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{AuthenticationEndpoint: "testHost"},
			store:         nil,
			server:        nil,
		}

		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "http://my.url.org/index.html", nil)

		sut.authenticationHandler(responseWriter, request)

		// when
		actualResponse := responseWriter.Result()

		// then
		assert.Equal(t, http.StatusSeeOther, actualResponse.StatusCode)
		assert.Contains(t, actualResponse.Header.Get("Location"), "testHost/instance-registration?ces_redirect_uri=http://my.url.org")
	})

	t.Run("authenticate with system token", func(t *testing.T) {
		// given
		store := getStoreMock()
		var previousID string
		abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, err := w.Write([]byte(`{ "id": "id", "secret": "secret"}`))
			require.NoError(t, err)
			previousID = r.URL.Query().Get("previousId")
		}))

		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{
				AuthenticationEndpoint: abcServer.URL,
				CredentialsStore:       "store",
			},
			store:  store,
			server: nil,
			client: &http.Client{},
		}

		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "http://my.url.org/index.html?id=myaccount&secret=myaccountsecret", nil)

		// when
		sut.authenticationHandler(responseWriter, request)

		// then
		actualResponse := responseWriter.Result()
		body, _ := ioutil.ReadAll(actualResponse.Body)
		assert.Equal(t, `<a href="/shutdown">Found</a>.`, strings.TrimSpace(string(body)))
		assert.Equal(t, http.StatusFound, actualResponse.StatusCode)
		assert.Equal(t, "", previousID)
	})

	t.Run("success with previous instance id", func(t *testing.T) {
		// given
		store := getStoreMock()
		var previousID string
		abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, err := w.Write([]byte(`{ "id": "id", "secret": "secret"}`))
			require.NoError(t, err)
			previousID = r.URL.Query().Get("previousId")
		}))

		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{AuthenticationEndpoint: abcServer.URL, PreviousInstanceID: "oldID", CredentialsStore: "store"},
			store:         store,
			server:        nil,
			client:        &http.Client{},
		}

		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "http://my.url.org/index.html?id=myaccount&secret=myaccountsecret", nil)

		// when
		sut.authenticationHandler(responseWriter, request)

		// then
		actualResponse := responseWriter.Result()
		body, _ := ioutil.ReadAll(actualResponse.Body)
		assert.Equal(t, `<a href="/shutdown">Found</a>.`, strings.TrimSpace(string(body)))
		assert.Equal(t, http.StatusFound, actualResponse.StatusCode)
		assert.Equal(t, "oldID", previousID)
	})

	t.Run("should return internal server error on http >= 300 on auth request", func(t *testing.T) {
		// given
		store := getStoreMock()
		var previousID string
		abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(300)
			previousID = r.URL.Query().Get("previousId")
		}))

		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{
				AuthenticationEndpoint: abcServer.URL,
				CredentialsStore:       "store",
			},
			store:  store,
			server: nil,
			client: &http.Client{},
		}

		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "http://my.url.org/index.html?id=myaccount&secret=myaccountsecret", nil)

		// when
		sut.authenticationHandler(responseWriter, request)

		// then
		actualResponse := responseWriter.Result()
		body, _ := ioutil.ReadAll(actualResponse.Body)
		assert.Equal(t, "backend returned failure status code", string(body))
		assert.Equal(t, http.StatusInternalServerError, actualResponse.StatusCode)
		assert.Equal(t, "", previousID)
	})

	t.Run("should return error if credentials coould not added to store", func(t *testing.T) {
		// given
		store := getStoreMock()
		store.On("Add", "errorStore", mock.Anything).Return(assert.AnError)
		var previousID string
		abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, err := w.Write([]byte(`{ "id": "id", "secret": "secret"}`))
			require.NoError(t, err)
			previousID = r.URL.Query().Get("previousId")
		}))

		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{
				AuthenticationEndpoint: abcServer.URL,
				CredentialsStore:       "errorStore",
			},
			store:  store,
			server: nil,
			client: &http.Client{},
		}

		responseWriter := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "http://my.url.org/index.html?id=myaccount&secret=myaccountsecret", nil)

		// when
		sut.authenticationHandler(responseWriter, request)

		// then
		actualResponse := responseWriter.Result()
		body, _ := ioutil.ReadAll(actualResponse.Body)
		assert.Equal(t, "assert.AnError general error for testing", string(body))
		assert.Equal(t, http.StatusInternalServerError, actualResponse.StatusCode)
		assert.Equal(t, "", previousID)
	})
}

func Test_realAuthenticationController_IsAuthenticated(t *testing.T) {
	t.Run("don't authenticate on empty store", func(t *testing.T) {
		// given
		store := getStoreMock()
		store.On("Get", "noStore").Return(nil)

		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{CredentialsStore: "noStore"},
			store:         store,
			server:        nil,
		}

		// when
		actual := sut.IsAuthenticated()

		// then
		assert.False(t, actual)
	})

	t.Run("authenticate on valid store", func(t *testing.T) {
		// given
		store := getStoreMock()
		sut := httpAuthenticationController{
			configuration: AuthenticationConfig{CredentialsStore: "store"},
			store:         store,
			server:        nil,
		}

		// when
		actual := sut.IsAuthenticated()

		// then
		assert.True(t, actual)
	})
}

func Test_realAuthenticationController_ShutdownHandler(t *testing.T) {
	// given
	mockResponse := httptest.NewRecorder()
	mockResponse.WriteHeader(http.StatusOK)
	_, _ = mockResponse.Write([]byte(`{ "id": "id", "secret": "secret"}`))

	sut := httpAuthenticationController{
		configuration: AuthenticationConfig{},
		store:         nil,
		server:        nil,
	}

	responseWriter := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://my.url.org/", nil)

	cancelCalled := false

	// when
	sut.shutdownHandler(responseWriter, request, func() {
		cancelCalled = true
	})

	// then
	assert.True(t, cancelCalled)
}

func getStoreMock() *mocks.Store {
	store := &mocks.Store{}
	creds := &core.Credentials{
		Username: "id",
		Password: "secret",
	}
	store.On("Add", "store", creds).Return(nil)
	store.On("Get", "store").Return(creds)
	return store
}
