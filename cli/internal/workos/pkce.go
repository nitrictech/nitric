package workos

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	workos_http "github.com/nitrictech/nitric/cli/internal/workos/http"
	"github.com/pkg/browser"
)

//go:embed login_success.html
var loginSuccessPage []byte

// The port that the local auth callback server will listen on.
// This is used to handle the callback from the WorkOS auth provider.
var LOCAL_PKCE_CALLBACK_PORT = 48321

// ErrPortNotAvailable is returned when the local auth callback port is not available.
var ErrPortNotAvailable = fmt.Errorf("port %d is not available, unable to start local auth callback server", LOCAL_PKCE_CALLBACK_PORT)

func (a *WorkOSAuth) performPKCE() error {
	// Check if the callback port is available and create a listener
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", LOCAL_PKCE_CALLBACK_PORT))
	if err != nil {
		return ErrPortNotAvailable
	}
	defer listener.Close()

	pkceChallenge, err := createPkceChallenge()
	if err != nil {
		return err
	}

	callbackResult := make(chan CallbackResult)

	callbackServer := &http.Server{
		// We only bind to loopback for security
		// The users own browser is the only client that should connect to this server, during a redirect
		Handler: http.HandlerFunc(newCallbackHandler(callbackResult, a.httpClient, pkceChallenge)),
	}

	go callbackServer.Serve(listener)
	defer func() {
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancelShutdown()

		callbackServer.Shutdown(shutdownCtx)
	}()

	authUrl, err := a.httpClient.GetAuthorizationUrl(workos_http.GetAuthorizationUrlOptions{
		Provider:            "authkit",
		RedirectURI:         fmt.Sprintf("http://127.0.0.1:%d/callback", LOCAL_PKCE_CALLBACK_PORT),
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
	result := <-callbackResult
	if result.Error != nil {
		fmt.Printf("\n%s Login failed due to an error: %s\n", style.Red(icons.Cross), result.Error)
		return result.Error
	}

	close(callbackResult)

	a.tokens = result.Tokens
	err = a.tokenStore.SaveTokens(a.tokens)
	if err != nil {
		fmt.Printf("\n%s Error saving tokens: %s\n", style.Red(icons.Cross), err)
		return err
	}

	return nil
}

func newCallbackHandler(callbackResult chan CallbackResult, client *workos_http.HttpClient, pkceChallenge *CodeVerifier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			// TODO: Handle favicon.ico requests
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := client.AuthenticateWithCode(code, pkceChallenge.Verifier)
		if err != nil {
			// TODO: make this pretty
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			callbackResult <- CallbackResult{Error: err}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(loginSuccessPage)

		go func() {
			callbackResult <- CallbackResult{Tokens: &Tokens{
				AccessToken:  res.AccessToken,
				RefreshToken: res.RefreshToken,
				User:         &res.User,
			}}
		}()

		fmt.Println("Callback received")
	}
}

type CallbackResult struct {
	Tokens *Tokens
	Error  error
}

type CodeVerifier struct {
	Verifier  string
	Challenge string
}

// CreatePkceChallenge generates both a code verifier and code challenge for PKCE
func createPkceChallenge() (*CodeVerifier, error) {
	codeVerifier, err := createCodeVerifier()
	if err != nil {
		return nil, err
	}
	codeChallenge, err := createCodeChallenge(codeVerifier)
	if err != nil {
		return nil, err
	}
	return &CodeVerifier{
		Verifier:  codeVerifier,
		Challenge: codeChallenge,
	}, nil
}

// createCodeVerifier generates a random code verifier
func createCodeVerifier() (string, error) {
	// Generate 96 bytes (equivalent to 96 * 4 = 384 bits from Uint32Array(96))
	randomBytes := make([]byte, 96)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64urlEncode(randomBytes), nil
}

// createCodeChallenge creates a SHA-256 hash of the code verifier
func createCodeChallenge(codeVerifier string) (string, error) {
	hashed := sha256.Sum256([]byte(codeVerifier))
	return base64urlEncode(hashed[:]), nil
}

// base64urlEncode encodes bytes to URL-safe base64
func base64urlEncode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	// Replace characters for URL-safe base64
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	// Remove padding
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}
