package workos

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/zalando/go-keyring"
)

// Store 1 token per API
var WORKOS_TOKEN_KEY = getWorkosTokenKey(config.GetNitricServerUrl().String())

func getWorkosTokenKey(apiUrl string) string {
	// Hash the API URL for a consistent length. We don't use the scheme or host, just the path
	hash := sha256.Sum256([]byte(apiUrl + ".workos"))
	return fmt.Sprintf("%x", hash)
}

type KeyringTokenStore struct {
	service string
}

func NewKeyringTokenStore(service string) *KeyringTokenStore {
	return &KeyringTokenStore{service: service}
}

func (s *KeyringTokenStore) GetTokens() (*Tokens, error) {
	token, err := keyring.Get(s.service, WORKOS_TOKEN_KEY)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve token from keyring: %w", err)
	}

	var tokens Tokens
	err = json.Unmarshal([]byte(token), &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &tokens, nil
}

func (s *KeyringTokenStore) SaveTokens(tokens *Tokens) error {
	json, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	return keyring.Set(s.service, WORKOS_TOKEN_KEY, string(json))
}

func (s *KeyringTokenStore) Clear() error {
	err := keyring.Delete(s.service, WORKOS_TOKEN_KEY)
	if err != nil {
		if err == keyring.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}
