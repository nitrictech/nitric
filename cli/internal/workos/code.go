package workos

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

type CodeVerifier struct {
	Verifier  string
	Challenge string
}

// CreatePkceChallenge generates both a code verifier and code challenge for PKCE
func CreatePkceChallenge() (*CodeVerifier, error) {
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
