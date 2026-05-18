package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	ClientTokenLength    = 32
	ClientTokenExpiresIn = 90 * 24 * time.Hour
)

func GenerateClientAccessToken() (rawToken string, tokenHash string, expiresAt time.Time, err error) {
	buffer := make([]byte, ClientTokenLength)
	if _, err = rand.Read(buffer); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generate client access token: %w", err)
	}

	rawToken = hex.EncodeToString(buffer)
	tokenHash = HashClientToken(rawToken)
	expiresAt = time.Now().UTC().Add(ClientTokenExpiresIn)

	return rawToken, tokenHash, expiresAt, nil
}

func HashClientToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func ValidateClientTokenFormat(token string) bool {
	if len(token) != ClientTokenLength*2 {
		return false
	}
	_, err := hex.DecodeString(token)
	return err == nil
}

func BuildClientMagicLink(baseURL, rawToken, requestID string) string {
	trimmedBaseURL := strings.TrimRight(baseURL, "/")
	query := url.Values{}
	query.Set("token", rawToken)
	query.Set("ref", requestID)
	return trimmedBaseURL + "/suivi?" + query.Encode()
}
