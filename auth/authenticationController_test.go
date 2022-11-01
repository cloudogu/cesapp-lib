package auth

import (
	"github.com/cloudogu/cesapp-lib/credentials"
	"github.com/cloudogu/cesapp-lib/credentials/mocks"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
)

func Test_realAuthenticationController_authenticationHandler(t *testing.T) {
	sut := realAuthenticationController{
		configuration: AuthenticationConfig{AuthenticationEndpoint: "testHost"},
		store:         nil,
		server:        nil,
	}

	responseWriter := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://my.url.org/index.html", nil)

	sut.authenticationHandler(responseWriter, request)

	actualResponse := responseWriter.Result()

	assert.Equal(t, http.StatusSeeOther, actualResponse.StatusCode)
	assert.Contains(t, actualResponse.Header.Get("Location"), "testHost/instance-registration?ces_redirect_uri=http://my.url.org")
}

func Test_realAuthenticationController_createSystemToken(t *testing.T) {
	store := getStoreMock()
	var previousID string
	abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(`{ "id": "id", "secret": "secret"}`))
		require.NoError(t, err)
		previousID = r.URL.Query().Get("previousId")
	}))

	sut := realAuthenticationController{
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

	sut.authenticationHandler(responseWriter, request)

	actualResponse := responseWriter.Result()
	body, _ := ioutil.ReadAll(actualResponse.Body)
	assert.Equal(t, `<a href="/shutdown">Found</a>.`, strings.TrimSpace(string(body)))
	assert.Equal(t, http.StatusFound, actualResponse.StatusCode)
	assert.Equal(t, "", previousID)
}

func Test_realAuthenticationController_createSystemTokenWithPreviousInstance(t *testing.T) {
	store := getStoreMock()
	var previousID string
	abcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(`{ "id": "id", "secret": "secret"}`))
		require.NoError(t, err)
		previousID = r.URL.Query().Get("previousId")
	}))

	sut := realAuthenticationController{
		configuration: AuthenticationConfig{AuthenticationEndpoint: abcServer.URL, PreviousInstanceID: "oldID", CredentialsStore: "store"},
		store:         store,
		server:        nil,
		client:        &http.Client{},
	}

	responseWriter := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://my.url.org/index.html?id=myaccount&secret=myaccountsecret", nil)

	sut.authenticationHandler(responseWriter, request)

	actualResponse := responseWriter.Result()
	body, _ := ioutil.ReadAll(actualResponse.Body)
	assert.Equal(t, `<a href="/shutdown">Found</a>.`, strings.TrimSpace(string(body)))
	assert.Equal(t, http.StatusFound, actualResponse.StatusCode)
	assert.Equal(t, "oldID", previousID)
}

func Test_realAuthenticationController_isAuthenticated_doesNotAuthenticate(t *testing.T) {
	store := getStoreMock()
	store.On("Get", "noStore").Return(nil)

	sut := realAuthenticationController{
		configuration: AuthenticationConfig{CredentialsStore: "noStore"},
		store:         store,
		server:        nil,
	}

	actual := sut.IsAuthenticated()

	assert.False(t, actual)
}

func Test_realAuthenticationController_isAuthenticated_authenticates(t *testing.T) {
	store := getStoreMock()
	sut := realAuthenticationController{
		configuration: AuthenticationConfig{CredentialsStore: "store"},
		store:         store,
		server:        nil,
	}

	actual := sut.IsAuthenticated()

	assert.True(t, actual)
}

func Test_realAuthenticationController_ShutdownHandler(t *testing.T) {
	mockResponse := httptest.NewRecorder()
	mockResponse.WriteHeader(http.StatusOK)
	_, _ = mockResponse.Write([]byte(`{ "id": "id", "secret": "secret"}`))

	sut := realAuthenticationController{
		configuration: AuthenticationConfig{},
		store:         nil,
		server:        nil,
	}

	responseWriter := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://my.url.org/", nil)

	cancelCalled := false
	sut.shutdownHandler(responseWriter, request, func() {
		cancelCalled = true
	})

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

func TestNewAuthenticationController(t *testing.T) {
	// given
	authConfig := AuthenticationConfig{}
	httpServer := NewHttpServer("")
	store, _ := credentials.NewStore("/tmp")

	// when
	result := NewAuthenticationController(authConfig, httpServer, nil, store)

	// then
	require.NotNil(t, result)
}
