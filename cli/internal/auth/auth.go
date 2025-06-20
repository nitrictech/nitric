package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/pkg/browser"
)

//go:embed login_success.html
var loginSuccessPage []byte

type Auth interface {
}

// TODO: These values are not secret, but we may want to pull them remotely incase of a change.
var AUTH_SERVER_PORT = 54321
var WORKOS_CLIENT_ID = "client_01JXRX76QGQMDJ621V6B16RVED"
var workosClient = workos.NewHttpClient(WORKOS_CLIENT_ID)

type WorkOsPKCE struct {
	pkceChallenge  *workos.CodeVerifier
	err            error
	callbackServer *http.Server
	done           chan error
}

func (p *WorkOsPKCE) getCallbackHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			w.WriteHeader(http.StatusBadRequest)

			p.done <- fmt.Errorf("login code was not provided with login callback")
			return
		}

		res, err := workosClient.AuthenticateWithCode(code, p.pkceChallenge.Verifier)
		if err != nil {
			// TODO: make this pretty
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			p.done <- err
			return
		}

		err = StoreWorkosToken(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			p.done <- err
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(loginSuccessPage)

		fmt.Printf("\n%s Login successful, welcome %s\n", style.Green(icons.Check), style.Teal(res.User.FirstName))

		p.done <- nil
	}
}

var nitric = style.Purple(icons.Lightning + " Nitric")

func (p *WorkOsPKCE) PerformPKCEFlow() error {
	fmt.Printf("\n%s Logging in...\n", nitric)

	router := http.NewServeMux()

	router.HandleFunc("/callback", p.getCallbackHandler())

	p.callbackServer = &http.Server{
		// We only bind to loopback for security
		// The users own browser is the only client that should connect to this server, during a redirect
		Addr:    fmt.Sprintf("127.0.0.1:%d", AUTH_SERVER_PORT),
		Handler: router,
	}

	p.done = make(chan error)
	go func() {
		p.callbackServer.ListenAndServe()
	}()

	// Start the Flow
	var err error
	p.pkceChallenge, err = workos.CreatePkceChallenge()
	if err != nil {
		return err
	}

	authUrl, err := workosClient.GetAuthorizationUrl(workos.GetAuthorizationUrlOptions{
		Provider:            "authkit",
		RedirectURI:         "http://127.0.0.1:54321/callback",
		CodeChallenge:       p.pkceChallenge.Challenge,
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

	err = <-p.done
	if err != nil {
		fmt.Printf("\n%s Login failed due to an error: %s\n", style.Red(icons.Cross), err)
	}

	p.callbackServer.Shutdown(context.Background())

	return nil
}
