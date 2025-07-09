package api

type TokenProvider interface {
	GetTokens() (*Tokens, error)
	SaveTokens(*Tokens) error
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
