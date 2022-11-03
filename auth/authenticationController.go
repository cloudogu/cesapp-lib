package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudogu/cesapp-lib/credentials"

	"github.com/cloudogu/cesapp-lib/core"

	"io/ioutil"
	"net/http"
)

const (
	failedToWriteFormat = "failed to write error message: %s"
)

var log = core.GetLogger()

// AuthenticationConfig represents the settings for an authentication.
type AuthenticationConfig struct {
	AuthenticationEndpoint string
	CredentialsStore       string
	PreviousInstanceID     string
}

// AuthenticationController provide actions for authenticating a user with the authentication endpoint
type AuthenticationController interface {
	// IsAuthenticated checks whether a user is authenticated or not.
	IsAuthenticated() bool
	// Serve starts the server to listen for incoming connections.
	Serve()
}

type realAuthenticationController struct {
	configuration AuthenticationConfig
	store         credentials.Store
	server        HttpServer
	client        *http.Client
}

// NewAuthenticationController creates a new instance of 'AuthenticationController'.
func NewAuthenticationController(authConfig AuthenticationConfig, httpServer HttpServer, httpClient *http.Client,
	credentialsStore credentials.Store) AuthenticationController {

	controller := &realAuthenticationController{
		configuration: authConfig,
		store:         credentialsStore,
		server:        httpServer,
		client:        httpClient,
	}
	return controller
}

// IsAuthenticated checks whether a user is authenticated or not.
func (controller *realAuthenticationController) IsAuthenticated() bool {
	return controller.store.Get(controller.configuration.CredentialsStore) != nil
}

// Serve starts the server to listen for incoming connections.
func (controller *realAuthenticationController) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/", controller.authenticationHandler)
	http.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		controller.shutdownHandler(writer, request, cancel)
	})
	go func() {
		err := controller.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("could not start web server: %s", err.Error())
		}
	}()
	<-ctx.Done()
}

// shutdownHandler handles the shutdown action of the server.
//noinspection GoUnusedParameter
func (controller *realAuthenticationController) shutdownHandler(w http.ResponseWriter, r *http.Request, cancel func()) {
	log.Info("Shutting web server down.")
	http.Redirect(w, r, controller.configuration.AuthenticationEndpoint, http.StatusSeeOther)
	cancel()
}

// authenticationHandler handles the authentication action of the server.
func (controller *realAuthenticationController) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	secret := r.URL.Query().Get("secret")

	if id != "" && secret != "" {
		systemToken, err := controller.createSystemToken(id, secret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errorf(failedToWriteFormat, err.Error())
			}
			return
		}

		creds := core.Credentials{Username: systemToken.ID, Password: systemToken.Secret}
		err = controller.store.Add(controller.configuration.CredentialsStore, &creds)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errorf(failedToWriteFormat, err.Error())
			}
			return
		}

		log.Info("Successfully authenticated.")
		http.Redirect(w, r, "/shutdown", http.StatusFound)
	} else {
		controller.redirectToAuthEndpoint(w, r)
		return
	}
}

func (controller *realAuthenticationController) redirectToAuthEndpoint(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("%s/instance-registration?ces_redirect_uri=http://%s", controller.configuration.AuthenticationEndpoint, r.Host)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// createSystemToken creates a system token from setup tokens.
func (controller *realAuthenticationController) createSystemToken(id string, secret string) (SystemToken, error) {
	systemToken := SystemToken{}

	instanceRegistrationUrl := controller.configuration.AuthenticationEndpoint + "/api/v1/instance-registrations/" + id

	if controller.configuration.PreviousInstanceID != "" {
		instanceRegistrationUrl += "?previousId=" + controller.configuration.PreviousInstanceID
	}

	client := controller.client

	// post temporal token to get login credentials for dogu registry
	var tokenSecret = []byte(`{"secret":"` + secret + `"}`)
	request, err := http.NewRequest(http.MethodPost, instanceRegistrationUrl, bytes.NewBuffer(tokenSecret))
	if err != nil {
		return systemToken, fmt.Errorf("could not create backend request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return systemToken, fmt.Errorf("backend request failed: %w", err)
	}

	// resource moved or an error occurred
	if resp.StatusCode >= 300 {
		return systemToken, fmt.Errorf("backend returned failure status code")
	}
	body := resp.Body
	defer func() {
		closeErr := body.Close()
		if closeErr != nil {
			log.Error(closeErr.Error())
		}
	}()
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return systemToken, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(bodyData, &systemToken)
	if err != nil {
		return systemToken, fmt.Errorf("failed to parse response body: %w", err)
	}
	return systemToken, nil
}
