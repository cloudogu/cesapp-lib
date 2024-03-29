package remote

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newHTTPRemote(t *testing.T) {
	t.Run("should return an error if creating the backoff fails", func(t *testing.T) {
		config := &core.Remote{
			RetryPolicy: core.RetryPolicy{
				Interval: -2,
			},
		}

		_, err := newHTTPRemote(config, nil)

		require.Error(t, err)
	})
}

func Test_checkStatusCode(t *testing.T) {
	t.Run("should return nil for HTTP 200", func(t *testing.T) {
		mockResp := &http.Response{}
		mockResp.Status = "200 OK"
		mockResp.StatusCode = http.StatusOK
		mockResp.Body = ioutil.NopCloser(strings.NewReader(`{"status": "is well"}`))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error for HTTP statuses >= 300", func(t *testing.T) {
		mockResp := &http.Response{}
		mockResp.Status = "300 Whoopsie!"
		mockResp.StatusCode = 300
		mockResp.Body = ioutil.NopCloser(strings.NewReader(`{"status": "I, uh, well... phew!"}`))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.Error(t, err)
		assert.Equal(t, err.Error(), "remote registry returns invalid status: 300 Whoopsie!: I, uh, well... phew!: (no error)")
	})

	t.Run("should return error for HTTP 400", func(t *testing.T) {
		const errorBody = "Do not use v1 endpoint for v2 dogu creation. Use v2 endpoint instead."

		mockResp := &http.Response{}
		mockResp.Status = http.StatusText(http.StatusBadRequest)
		mockResp.StatusCode = http.StatusBadRequest
		mockResp.Body = ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"error": "%s"}`, errorBody)))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.Error(t, err)
		assert.Equal(t, err.Error(), "remote registry returns invalid status: Bad Request: 400: Do not use v1 endpoint for v2 dogu creation. Use v2 endpoint instead.")
	})

	t.Run("should return custom error for HTTP 401", func(t *testing.T) {
		mockResp := &http.Response{}
		mockResp.Status = http.StatusText(http.StatusUnauthorized)
		mockResp.StatusCode = http.StatusUnauthorized
		mockResp.Body = ioutil.NopCloser(strings.NewReader(`{"status": "unauthorized"}`))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.Error(t, err)
		assert.Equal(t, errUnauthorized, err)
	})

	t.Run("should return custom error for HTTP 403", func(t *testing.T) {
		mockResp := &http.Response{}
		mockResp.Status = http.StatusText(http.StatusForbidden)
		mockResp.StatusCode = http.StatusForbidden
		mockResp.Body = ioutil.NopCloser(strings.NewReader(`{"status": "forbidden"}`))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.Error(t, err)
		assert.Equal(t, errForbidden, err)
	})

	t.Run("should return custom error for HTTP 404", func(t *testing.T) {
		mockResp := &http.Response{}
		mockResp.Status = http.StatusText(http.StatusNotFound)
		mockResp.StatusCode = http.StatusNotFound
		mockResp.Body = ioutil.NopCloser(strings.NewReader(`{"status": "forbidden"}`))

		// when
		err := checkStatusCode(mockResp)

		// then
		require.Error(t, err)
		assert.Equal(t, errNotFound, err)
	})
}

func Test_extractRemoteErrorBody(t *testing.T) {
	t.Run("should return error body", func(t *testing.T) {
		responseBody := ioutil.NopCloser(strings.NewReader(`{"error": "the error text"}`))
		// when
		actual := extractRemoteBody(responseBody, 400)

		// then
		assert.Equal(t, "400: the error text", actual)
	})

	t.Run("should include generic error for truncated json", func(t *testing.T) {
		responseBody := ioutil.NopCloser(strings.NewReader(`{"error": "the erro...`))
		// when
		actual := extractRemoteBody(responseBody, 400)

		// then
		assert.Equal(t, "error", actual)
	})
}

func Test_remoteResponseBody_String(t *testing.T) {
	type fields struct {
		statusCode int
		Status     string
		Error      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"return mixed string", fields{Status: "aaa", Error: "bbb"}, "aaa: bbb"},
		{"return only status", fields{Status: "aaa", Error: ""}, "aaa: (no error)"},
		{"return only error", fields{statusCode: 123, Error: "bbb"}, "123: bbb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseBody := &remoteResponseBody{
				statusCode: tt.fields.statusCode,
				Status:     tt.fields.Status,
				Error:      tt.fields.Error,
			}
			if got := responseBody.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRemote_Delete(t *testing.T) {
	aDogu := &core.Dogu{
		Name:    "testing/nginx",
		Image:   "staging-registry.cloudogu.com/testing/nginx",
		Version: "2.3.4-5",
	}
	netRetrier := retrier.New(
		retrier.ExponentialBackoff(1, 1*time.Millisecond),
		retrier.BlacklistClassifier{errUnauthorized, errForbidden},
	)

	t.Run("should success when remote returns HTTP 204", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		}))
		defer ts.Close()
		urlSchema := NewURLSchemaByName("default", ts.URL+"/api/v2/dogus")
		httpClient := &http.Client{Timeout: 1 * time.Second}

		sut := httpRemote{
			retrier:   netRetrier,
			urlSchema: urlSchema,
			client:    httpClient,
		}

		// when
		err := sut.Delete(aDogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail when remote returns HTTP 401", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
		}))
		defer ts.Close()
		urlSchema := NewURLSchemaByName("default", ts.URL+"/api/v2/dogus")
		httpClient := &http.Client{Timeout: 1 * time.Second}

		sut := httpRemote{
			retrier:   netRetrier,
			urlSchema: urlSchema,
			client:    httpClient,
		}

		// when
		err := sut.Delete(aDogu)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), errUnauthorized.Error())
	})

	t.Run("should properly use credentials", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, _ := r.BasicAuth()
			// then
			assert.Equal(t, username, "username")
			assert.Equal(t, password, "password")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()
		urlSchema := NewURLSchemaByName("default", ts.URL+"/api/v2/dogus")
		httpClient := &http.Client{Timeout: 1 * time.Second}

		sut := httpRemote{
			retrier:     netRetrier,
			urlSchema:   urlSchema,
			client:      httpClient,
			credentials: &core.Credentials{Username: "username", Password: "password"},
		}

		// when
		err := sut.Delete(aDogu)

		require.NoError(t, err)
	})
}
