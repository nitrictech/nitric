package auth

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/zalando/go-keyring"
)

var ErrNotFound = errors.New("no token found")

const KEYRING_SERVICE = "nitric.v2.cli"

// Store 1 token per API
var WORKOS_TOKEN_KEY = getWorkosTokenKey(config.GetApiUrl())

func getWorkosTokenKey(apiUrl string) string {
	// Hash the API URL for a consistent length.
	hash := sha256.Sum256([]byte(apiUrl + ".workos"))
	return fmt.Sprintf("%x", hash)
}

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
	client, err := getWorkOSClient()
	if err != nil {
		return nil, err
	}

	workosToken, err = client.AuthenticateWithRefreshToken(workosToken.RefreshToken, nil)
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

// GetOrRefreshWorkosToken retrieves the authentication token from the keyring, and refreshes it if it's expired
func GetOrRefreshWorkosToken() (*workos.AuthenticationResponse, error) {
	workosToken, err := GetWorkosToken()
	if err != nil {
		return nil, err
	}

	// Decode the JWT to check if it's expired
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(workosToken.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		client, err := getWorkOSClient()
		if err != nil {
			return nil, err
		}

		jwk, err := client.GetJWK(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get token validation key: %v", err)
		}

		return jwkToRSAPublicKey(jwk)
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	// Add a 1 second buffer to the expiry time to account for a slight delay in the token being sent to the server
	// i.e. the token must remain valid for at least another second, or we'll refresh it early for good measure
	if err != nil || !parsedToken.Valid || claims.ExpiresAt.Before(time.Now().Add(1+time.Second)) {
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
