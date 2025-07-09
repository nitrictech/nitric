package workos

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nitrictech/nitric/cli/internal/workos/http"
)

var (
	ErrNotFound        = errors.New("no token found")
	ErrUnauthenticated = errors.New("unauthenticated")
)

type TokenStore interface {
	GetTokens() (*Tokens, error)
	SaveTokens(*Tokens) error
	Clear() error
}

type Tokens struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         *http.User `json:"user"`
}

type WorkOSAuth struct {
	tokenStore TokenStore
	tokens     *Tokens
	httpClient *http.HttpClient
}

func NewWorkOSAuth(tokenStore TokenStore, clientID string, endpoint string) *WorkOSAuth {
	httpClient := http.NewHttpClient(clientID, http.WithHostname(endpoint))

	return &WorkOSAuth{tokenStore: tokenStore, httpClient: httpClient}
}

func (a *WorkOSAuth) Login() (*http.User, error) {
	err := a.performPKCE()
	if err != nil {
		return nil, err
	}

	return a.tokens.User, nil
}

func (a *WorkOSAuth) GetAccessToken() (string, error) {

	if a.tokens == nil {
		tokens, err := a.tokenStore.GetTokens()
		if err != nil {
			return "", fmt.Errorf("no stored tokens found, please login: %w", err)
		}
		a.tokens = tokens
	}

	// Decode the JWT to check if it's expired
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(a.tokens.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		return a.httpClient.GetRSAPublicKey(kid)
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	// Add a 1 second buffer to the expiry time to account for a slight delay in the token being sent to the server
	// i.e. the token must remain valid for at least another second, or we'll refresh it early for good measure
	if err != nil || !parsedToken.Valid || claims.ExpiresAt.Before(time.Now().Add(1+time.Second)) {
		if err := a.refreshToken(); err != nil {
			return "", fmt.Errorf("token refresh failed: %w", err)
		}
	}

	return a.tokens.AccessToken, nil
}

func (a *WorkOSAuth) refreshToken() error {
	if a.tokens.RefreshToken == "" {
		return fmt.Errorf("no refresh token", ErrUnauthenticated)
	}

	workosToken, err := a.httpClient.AuthenticateWithRefreshToken(a.tokens.RefreshToken, nil)
	if err != nil {
		return err
	}

	a.tokens = &Tokens{
		AccessToken:  workosToken.AccessToken,
		RefreshToken: workosToken.RefreshToken,
		User:         &workosToken.User,
	}

	err = a.tokenStore.SaveTokens(a.tokens)
	if err != nil {
		return err
	}

	return nil
}

func (a *WorkOSAuth) Logout() error {
	a.tokens = nil
	return a.tokenStore.Clear()
}
