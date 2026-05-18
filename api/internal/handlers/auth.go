package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/optea-tech/api/internal/models"
)

type AuthHandler struct {
	store     *Store
	db        *pgxpool.Pool
	validator *validator.Validate
	jwtSecret string
	ttl       int64
}

func NewAuthHandler(store *Store, db *pgxpool.Pool, jwtSecret string, ttl int64) *AuthHandler {
	return &AuthHandler{
		store:     store,
		db:        db,
		validator: validator.New(),
		jwtSecret: jwtSecret,
		ttl:       ttl,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input models.LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validator.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}

	adminEmail := input.Email
	if h.db != nil {
		passwordHash, err := h.getAdminPasswordHash(c.Context(), input.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)) != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		_ = h.touchLastLogin(c.Context(), input.Email)
	} else {
		fallbackEmail := getEnv("ADMIN_EMAIL", "admin@optea.tech")
		fallbackPassword := getEnv("ADMIN_PASSWORD", "admin123")
		if input.Email != fallbackEmail || input.Password != fallbackPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		adminEmail = fallbackEmail
	}

	accessToken, err := h.signAccessToken(adminEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to sign token"})
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to generate refresh token"})
	}

	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	if h.db != nil {
		if err := h.saveRefreshToken(c.Context(), refreshToken, adminEmail, refreshExpiry, sessionFingerprint(c), c.Get("User-Agent"), c.IP()); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to persist refresh token"})
		}
	} else {
		h.store.mu.Lock()
		h.store.refresh[refreshToken] = refreshExpiry
		h.store.mu.Unlock()
	}

	return c.JSON(models.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: h.ttl})
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().Body(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	adminEmail := getEnv("ADMIN_EMAIL", "admin@optea.tech")
	if h.db != nil {
		newRefreshToken, err := createRefreshToken()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to generate refresh token"})
		}
		refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
		email, valid, err := h.rotateRefreshToken(
			c.Context(),
			payload.RefreshToken,
			newRefreshToken,
			refreshExpiry,
			sessionFingerprint(c),
			c.Get("User-Agent"),
			c.IP(),
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to validate refresh token"})
		}
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
		}
		adminEmail = email

		accessToken, signErr := h.signAccessToken(adminEmail)
		if signErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to sign token"})
		}
		return c.JSON(models.AuthResponse{AccessToken: accessToken, RefreshToken: newRefreshToken, ExpiresIn: h.ttl})
	} else {
		h.store.mu.Lock()
		expiresAt, ok := h.store.refresh[payload.RefreshToken]
		if !ok || time.Now().After(expiresAt) {
			h.store.mu.Unlock()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
		}
		delete(h.store.refresh, payload.RefreshToken)
		h.store.mu.Unlock()
	}

	accessToken, err := h.signAccessToken(adminEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to sign token"})
	}

	newRefreshToken, err := createRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to generate refresh token"})
	}

	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	h.store.mu.Lock()
	h.store.refresh[newRefreshToken] = refreshExpiry
	h.store.mu.Unlock()

	return c.JSON(models.AuthResponse{AccessToken: accessToken, RefreshToken: newRefreshToken, ExpiresIn: h.ttl})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = c.Bind().Body(&payload)

	if h.db != nil {
		_ = h.revokeRefreshToken(c.Context(), payload.RefreshToken)
	} else {
		h.store.mu.Lock()
		delete(h.store.refresh, payload.RefreshToken)
		h.store.mu.Unlock()
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) signAccessToken(email string) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(h.ttl) * time.Second).Unix(),
		"rol": "admin",
	})
	return token.SignedString([]byte(h.jwtSecret))
}

func createRefreshToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (h *AuthHandler) saveRefreshToken(ctx context.Context, token string, userEmail string, expiresAt time.Time, fingerprintHash string, userAgent string, ipAddress string) error {
	_, err := h.db.Exec(ctx, `
		INSERT INTO admin_refresh_tokens (token, user_email, expires_at, fingerprint_hash, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, token, userEmail, expiresAt, fingerprintHash, userAgent, ipAddress)
	return err
}

func (h *AuthHandler) rotateRefreshToken(ctx context.Context, token string, newToken string, newExpiresAt time.Time, fingerprintHash string, userAgent string, ipAddress string) (string, bool, error) {
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return "", false, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userEmail string
	var expiresAt time.Time
	var revokedAt *time.Time
	var storedFingerprint *string
	err = tx.QueryRow(ctx, `
		SELECT user_email, expires_at, revoked_at, fingerprint_hash
		FROM admin_refresh_tokens
		WHERE token = $1
		FOR UPDATE
	`, token).Scan(&userEmail, &expiresAt, &revokedAt, &storedFingerprint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}

	if revokedAt != nil || time.Now().After(expiresAt) {
		return "", false, nil
	}
	if storedFingerprint != nil && *storedFingerprint != "" && *storedFingerprint != fingerprintHash {
		return "", false, nil
	}

	if _, err = tx.Exec(ctx, `
		UPDATE admin_refresh_tokens
		SET revoked_at = NOW(), replaced_by_token = $2
		WHERE token = $1
	`, token, newToken); err != nil {
		return "", false, err
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO admin_refresh_tokens (token, user_email, expires_at, fingerprint_hash, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, newToken, userEmail, newExpiresAt, fingerprintHash, userAgent, ipAddress); err != nil {
		return "", false, err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", false, err
	}
	_, _ = h.db.Exec(ctx, `DELETE FROM admin_refresh_tokens WHERE expires_at <= NOW()`)
	return userEmail, true, nil
}

func (h *AuthHandler) revokeRefreshToken(ctx context.Context, token string) error {
	_, err := h.db.Exec(ctx, `
		UPDATE admin_refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`, token)
	return err
}

func sessionFingerprint(c fiber.Ctx) string {
	raw := strings.TrimSpace(c.Get("X-Client-Fingerprint"))
	if raw == "" {
		raw = c.Get("User-Agent") + "|" + c.IP()
	}
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func (h *AuthHandler) getAdminPasswordHash(ctx context.Context, email string) (string, error) {
	var hash string
	err := h.db.QueryRow(ctx, `SELECT password_hash FROM admin_users WHERE email = $1 LIMIT 1`, email).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (h *AuthHandler) touchLastLogin(ctx context.Context, email string) error {
	_, err := h.db.Exec(ctx, `UPDATE admin_users SET last_login_at = NOW() WHERE email = $1`, email)
	return err
}
