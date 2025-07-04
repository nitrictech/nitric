package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/pkg/browser"
)

//go:embed login_success.html
var loginSuccessPage []byte

// The port that the local auth callback server will listen on.
// This is used to handle the callback from the WorkOS auth provider.
var LOCAL_AUTH_CALLBACK_PORT = 48321

func WithAuthHeader(req *http.Request) {
	token, err := GetOrRefreshWorkosToken()
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
}

// ErrPortNotAvailable is returned when the local auth callback port is not available.
var ErrPortNotAvailable = fmt.Errorf("port %d is not available, unable to start local auth callback server", LOCAL_AUTH_CALLBACK_PORT)

func newCallbackHandler(callbackResult chan error, client *workos.HttpClient, pkceChallenge *workos.CodeVerifier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			w.WriteHeader(http.StatusBadRequest)

			callbackResult <- fmt.Errorf("login code was not provided with login callback")
			return
		}

		res, err := client.AuthenticateWithCode(code, pkceChallenge.Verifier)
		if err != nil {
			// TODO: make this pretty
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			callbackResult <- err
			return
		}

		err = StoreWorkosToken(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			callbackResult <- err
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(loginSuccessPage)

		fmt.Printf("\n%s Login successful, welcome %s\n", style.Green(icons.Check), style.Teal(res.User.FirstName))

		callbackResult <- nil
	}
}

func PerformPKCEFlow() error {
	// Check if the callback port is available and create a listener
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", LOCAL_AUTH_CALLBACK_PORT))
	if err != nil {
		return ErrPortNotAvailable
	}
	defer listener.Close()

	client, err := getWorkOSClient()
	if err != nil {
		return err
	}

	pkceChallenge, err := workos.CreatePkceChallenge()
	if err != nil {
		return err
	}

	callbackResult := make(chan error)

	callbackServer := &http.Server{
		// We only bind to loopback for security
		// The users own browser is the only client that should connect to this server, during a redirect
		Handler: http.HandlerFunc(newCallbackHandler(callbackResult, client, pkceChallenge)),
	}

	go callbackServer.Serve(listener)
	defer callbackServer.Shutdown(context.Background())

	authUrl, err := client.GetAuthorizationUrl(workos.GetAuthorizationUrlOptions{
		Provider:            "authkit",
		RedirectURI:         fmt.Sprintf("http://127.0.0.1:%d/callback", LOCAL_AUTH_CALLBACK_PORT),
		CodeChallenge:       pkceChallenge.Challenge,
		CodeChallengeMethod: "S256",
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nOpening browser to %s\n", style.Gray(authUrl))

	// Open the browser
	err = browser.OpenURL(authUrl)
	if err != nil {
		return err
	}

	// Wait for the callback to be received or the server to shutdown
	err = <-callbackResult
	if err != nil {
		fmt.Printf("\n%s Login failed due to an error: %s\n", style.Red(icons.Cross), err)
	}

	return nil
}

var workosClient *workos.HttpClient

func getWorkOSClient() (*workos.HttpClient, error) {
	if workosClient != nil {
		return workosClient, nil
	}

	nitricApiClient := api.NewNitricApiClient(config.GetNitricServerUrl())
	workosDetails, err := nitricApiClient.GetWorkOSPublicDetails()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset by peer") {
			return nil, fmt.Errorf("unable to connect to the Nitric API. Please check your connection and try again. If the problem persists, please contact support.")
		}

		return nil, err
	}

	workosClient = workos.NewHttpClient(workosDetails.ClientID, workos.WithHostname(workosDetails.ApiHostname))
	return workosClient, nil
}
