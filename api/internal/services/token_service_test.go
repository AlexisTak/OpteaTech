package services

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGenerateClientAccessToken(t *testing.T) {
	rawToken, tokenHash, expiresAt, err := GenerateClientAccessToken()
	if err != nil {
		t.Fatalf("GenerateClientAccessToken returned error: %v", err)
	}

	if len(rawToken) != ClientTokenLength*2 {
		t.Fatalf("expected raw token length %d, got %d", ClientTokenLength*2, len(rawToken))
	}

	if !ValidateClientTokenFormat(rawToken) {
		t.Fatal("generated token should have a valid format")
	}

	if tokenHash != HashClientToken(rawToken) {
		t.Fatal("generated hash does not match the raw token")
	}

	if time.Until(expiresAt) < 89*24*time.Hour {
		t.Fatal("token expiration should be roughly 90 days in the future")
	}
}

func TestValidateClientTokenFormat(t *testing.T) {
	if ValidateClientTokenFormat("invalid-token") {
		t.Fatal("expected invalid token format to be rejected")
	}

	if ValidateClientTokenFormat(strings.Repeat("a", ClientTokenLength*2-1)) {
		t.Fatal("expected short token to be rejected")
	}
}

func TestBuildClientMagicLink(t *testing.T) {
	magicLink := BuildClientMagicLink("https://optea.tech/", "abc123", "req-1")
	parsed, err := url.Parse(magicLink)
	if err != nil {
		t.Fatalf("magic link is not a valid URL: %v", err)
	}

	if parsed.Path != "/suivi" {
		t.Fatalf("expected /suivi path, got %s", parsed.Path)
	}

	if parsed.Query().Get("token") != "abc123" {
		t.Fatal("magic link token query parameter mismatch")
	}

	if parsed.Query().Get("ref") != "req-1" {
		t.Fatal("magic link ref query parameter mismatch")
	}
}
