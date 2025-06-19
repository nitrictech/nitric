package token

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/zalando/go-keyring"
)

var ErrNotFound = errors.New("no token found")

const KEYRING_SERVICE = "nitric.v2.cli"
const WORKOS_TOKEN_KEY = "workos"

// StoreToken saves the authentication token to the keyring
func StoreWorkosToken(token *workos.AuthenticationResponse) error {
	json, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	err = keyring.Set(KEYRING_SERVICE, WORKOS_TOKEN_KEY, string(json))
	if err != nil {
		return fmt.Errorf("failed to store token in keyring: %w", err)
	}
	return nil
}

// GetToken retrieves the authentication token from the keyring
func GetWorkosToken() (*workos.AuthenticationResponse, error) {
	token, err := keyring.Get(KEYRING_SERVICE, WORKOS_TOKEN_KEY)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve token from keyring: %w", err)
	}

	var workosToken workos.AuthenticationResponse
	err = json.Unmarshal([]byte(token), &workosToken)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}
	return &workosToken, nil
}

// DeleteToken removes the authentication token from the keyring
func DeleteWorkosToken() error {
	err := keyring.Delete(KEYRING_SERVICE, WORKOS_TOKEN_KEY)
	if err != nil {
		if err == keyring.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}
