package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v4"
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

func RefreshWorkosToken(workosToken *workos.AuthenticationResponse) (*workos.AuthenticationResponse, error) {
	workosToken, err := workosClient.AuthenticateWithRefreshToken(workosToken.RefreshToken, nil)
	if err != nil {
		return nil, err
	}

	err = StoreWorkosToken(workosToken)
	if err != nil {
		return nil, err
	}

	return workosToken, nil
}

// GetToken retrieves the authentication token from the keyring, use GetOrRefreshWorkosToken instead
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

func GetOrRefreshWorkosToken() (*workos.AuthenticationResponse, error) {
	workosToken, err := GetWorkosToken()
	if err != nil {
		return nil, err
	}

	// Decode the JWT to check if it's expired
	parsedToken, err := jwt.Parse(workosToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		jwk, err := workosClient.GetJWK(kid)
		if err != nil {
			return nil, err
		}

		return jwkToRSAPublicKey(jwk)
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	if err != nil || !parsedToken.Valid {
		return RefreshWorkosToken(workosToken)
	}

	return workosToken, nil
}

func jwkToRSAPublicKey(jwk workos.JWK) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	// Create the RSA public key
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}, nil
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
