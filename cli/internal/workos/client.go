package workos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

const DEFAULT_HOSTNAME = "api.workos.com"

// Errors
type CodeExchangeError struct {
	Message string
}

func (e *CodeExchangeError) Error() string {
	return e.Message
}

type RefreshError struct {
	Message string
}

func (e *RefreshError) Error() string {
	return e.Message
}

// Authentication response types
type AuthenticationResponseRaw struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	LastSignInAt      string `json:"last_sign_in_at"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ExternalID        string `json:"external_id"`
}

type AuthenticationResponse struct {
	AccessToken  string
	RefreshToken string
	User         User
}

type JWKsResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty     string   `json:"kty"`
	Kid     string   `json:"kid"`
	Use     string   `json:"use"`
	Alg     string   `json:"alg"`
	N       string   `json:"n"`
	E       string   `json:"e"`
	X5c     []string `json:"x5c"`
	X5tS256 string   `json:"x5t#S256"`
}

// Authorization URL options
type GetAuthorizationUrlOptions struct {
	ConnectionID        string
	Context             string
	DomainHint          string
	LoginHint           string
	OrganizationID      string
	Provider            string
	RedirectURI         string
	State               string
	ScreenHint          string
	PasswordResetToken  string
	InvitationToken     string
	CodeChallenge       string
	CodeChallengeMethod string
}

// HttpClient represents the WorkOS HTTP client
type HttpClient struct {
	baseURL  string
	clientID string
	client   *http.Client
}

// NewHttpClient creates a new WorkOS HTTP client
func NewHttpClient(clientID string, options ...ClientOption) *HttpClient {
	config := &clientConfig{
		hostname: DEFAULT_HOSTNAME,
		scheme:   "https",
	}

	for _, option := range options {
		option(config)
	}

	baseURL := &url.URL{
		Scheme: config.scheme,
		Host:   config.hostname,
	}

	if config.port != 0 {
		baseURL.Host = fmt.Sprintf("%s:%d", config.hostname, config.port)
	}

	return &HttpClient{
		baseURL:  baseURL.String(),
		clientID: clientID,
		client:   &http.Client{},
	}
}

type clientConfig struct {
	hostname string
	port     int
	scheme   string
}

type ClientOption func(*clientConfig)

func WithHostname(hostname string) ClientOption {
	return func(c *clientConfig) {
		c.hostname = hostname
	}
}

func WithPort(port int) ClientOption {
	return func(c *clientConfig) {
		c.port = port
	}
}

func WithScheme(scheme string) ClientOption {
	return func(c *clientConfig) {
		c.scheme = scheme
	}
}

func (h *HttpClient) GetJWK(kid string) (JWK, error) {
	jwks, err := h.GetJWKs()
	if err != nil {
		return JWK{}, err
	}

	for _, jwk := range jwks {
		if jwk.Kid == kid {
			return jwk, nil
		}
	}

	return JWK{}, fmt.Errorf("JWK not found")
}

func (h *HttpClient) GetJWKs() ([]JWK, error) {
	jwkPath := path.Join("sso/jwks", h.clientID)

	response, err := h.get(jwkPath)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jwksResponse JWKsResponse
	if err := json.Unmarshal(body, &jwksResponse); err != nil {
		return nil, err
	}

	return jwksResponse.Keys, nil
}

// AuthenticateWithCode authenticates using an authorization code
func (h *HttpClient) AuthenticateWithCode(code, codeVerifier string) (*AuthenticationResponse, error) {
	body := map[string]interface{}{
		"code":          code,
		"client_id":     h.clientID,
		"grant_type":    "authorization_code",
		"code_verifier": codeVerifier,
	}

	response, err := h.post("/user_management/authenticate", body)
	if err != nil {
		return nil, err
	}

	// read the body into a string
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		var data AuthenticationResponseRaw
		if err := json.Unmarshal(resBody, &data); err != nil {
			return nil, err
		}
		return deserializeAuthenticationResponse(data), nil
	}

	return nil, &CodeExchangeError{Message: fmt.Sprintf("Error authenticating with API, status: %d, body: %s", response.StatusCode, string(resBody))}
}

// post performs a POST request to the specified path
func (h *HttpClient) post(path string, body map[string]interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Join the base URL with the path safely
	fullURL := baseURL.JoinPath(path)

	req, err := http.NewRequest("POST", fullURL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")

	return h.client.Do(req)
}

func (h *HttpClient) get(path string) (*http.Response, error) {
	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	fullURL := baseURL.JoinPath(path)

	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")

	return h.client.Do(req)
}

// AuthenticateWithRefreshToken authenticates using a refresh token
func (h *HttpClient) AuthenticateWithRefreshToken(refreshToken string, organizationId *string) (*AuthenticationResponse, error) {
	body := map[string]interface{}{
		"client_id":     h.clientID,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	if organizationId != nil {
		body["organization_id"] = *organizationId
	}

	response, err := h.post("/user_management/authenticate", body)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		var data AuthenticationResponseRaw
		if err := json.Unmarshal(resBody, &data); err != nil {
			return nil, err
		}
		return deserializeAuthenticationResponse(data), nil
	}

	return nil, &RefreshError{Message: fmt.Sprintf("Error authenticating with API, status: %d, body: %s", response.StatusCode, string(resBody))}
}

// GetAuthorizationUrl generates an authorization URL
func (h *HttpClient) GetAuthorizationUrl(options GetAuthorizationUrlOptions) (string, error) {
	if options.Provider == "" && options.ConnectionID == "" && options.OrganizationID == "" {
		return "", fmt.Errorf("incomplete arguments. need to specify either a 'connectionId', 'organizationId', or 'provider'")
	}

	if options.Provider != "" && options.Provider != "authkit" && options.ScreenHint != "" {
		return "", fmt.Errorf("'screenHint' is only supported for 'authkit' provider")
	}

	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Join the base URL with the authorize path
	authURL := baseURL.JoinPath("user_management", "authorize")

	// Build query parameters
	params := url.Values{}

	if options.ConnectionID != "" {
		params.Set("connection_id", options.ConnectionID)
	}
	if options.Context != "" {
		params.Set("context", options.Context)
	}
	if options.OrganizationID != "" {
		params.Set("organization_id", options.OrganizationID)
	}
	if options.DomainHint != "" {
		params.Set("domain_hint", options.DomainHint)
	}
	if options.LoginHint != "" {
		params.Set("login_hint", options.LoginHint)
	}
	if options.Provider != "" {
		params.Set("provider", options.Provider)
	}
	params.Set("client_id", h.clientID)
	if options.RedirectURI != "" {
		params.Set("redirect_uri", options.RedirectURI)
	}
	params.Set("response_type", "code")
	if options.State != "" {
		params.Set("state", options.State)
	}
	if options.ScreenHint != "" {
		params.Set("screen_hint", options.ScreenHint)
	}
	if options.InvitationToken != "" {
		params.Set("invitation_token", options.InvitationToken)
	}
	if options.PasswordResetToken != "" {
		params.Set("password_reset_token", options.PasswordResetToken)
	}
	if options.CodeChallenge != "" {
		params.Set("code_challenge", options.CodeChallenge)
	}
	if options.CodeChallengeMethod != "" {
		params.Set("code_challenge_method", options.CodeChallengeMethod)
	}

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

// GetLogoutUrl generates a logout URL
func (h *HttpClient) GetLogoutUrl(sessionID, returnTo string) string {
	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		// If base URL is invalid, return a basic string (this shouldn't happen with proper initialization)
		return ""
	}

	// Join the base URL with the logout path
	logoutURL := baseURL.JoinPath("user_management", "sessions", "logout")

	// Build query parameters
	params := url.Values{}
	params.Set("session_id", sessionID)
	if returnTo != "" {
		params.Set("return_to", returnTo)
	}

	logoutURL.RawQuery = params.Encode()
	return logoutURL.String()
}

// deserializeAuthenticationResponse converts the raw response to the structured response
func deserializeAuthenticationResponse(raw AuthenticationResponseRaw) *AuthenticationResponse {
	return &AuthenticationResponse{
		AccessToken:  raw.AccessToken,
		RefreshToken: raw.RefreshToken,
		User:         raw.User,
	}
}
