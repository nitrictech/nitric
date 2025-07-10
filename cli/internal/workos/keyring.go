package workos

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

func hashTokenKey(tokenKey string) string {
	// Hash the token key for a consistent length.
	hash := sha256.Sum256([]byte(tokenKey))
	return fmt.Sprintf("%x", hash)
}

type KeyringTokenStore struct {
	service  string
	tokenKey string
}

func NewKeyringTokenStore(serviceName, tokenKey string) (*KeyringTokenStore, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	if tokenKey == "" {
		return nil, fmt.Errorf("token key is required")
	}

	hashedTokenKey := hashTokenKey(tokenKey)

	return &KeyringTokenStore{service: serviceName, tokenKey: hashedTokenKey}, nil
}

func (s *KeyringTokenStore) GetTokens() (*Tokens, error) {
	token, err := keyring.Get(s.service, s.tokenKey)
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

	return keyring.Set(s.service, s.tokenKey, string(json))
}

func (s *KeyringTokenStore) Clear() error {
	err := keyring.Delete(s.service, s.tokenKey)
	if err != nil {
		if err == keyring.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}
