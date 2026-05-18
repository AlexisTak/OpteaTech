package middleware

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/optea-tech/api/internal/models"
	"github.com/optea-tech/api/internal/repository"
	"github.com/optea-tech/api/internal/services"
)

const ClientRequestKey = "client_request"

func ClientAuth(repo *repository.RequestsRepo, accessLogRepo *repository.AccessLogRepo) fiber.Handler {
	return func(c fiber.Ctx) error {
		if repo == nil || !repo.Ready() {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "client portal unavailable"})
		}

		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")
		rawToken := strings.TrimSpace(c.Query("token"))
		if rawToken == "" {
			if token, err := extractBearerToken(c.Get("Authorization")); err == nil {
				rawToken = token
			}
		}

		if rawToken == "" {
			logClientAccess(accessLogRepo, nil, "client_auth", ipAddress, userAgent, false, "token_missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token d'acces requis.", "code": "TOKEN_MISSING"})
		}

		if !services.ValidateClientTokenFormat(rawToken) {
			logClientAccess(accessLogRepo, nil, "client_auth", ipAddress, userAgent, false, "token_invalid_format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token invalide.", "code": "TOKEN_INVALID"})
		}

		request, err := repo.FindByTokenHash(c.Context(), services.HashClientToken(rawToken))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logClientAccess(accessLogRepo, nil, "client_auth", ipAddress, userAgent, false, "token_not_found")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Lien invalide ou expire. Contactez-nous pour obtenir un nouveau lien.", "code": "TOKEN_INVALID"})
			}

			slog.Error("client token lookup failed", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur interne. Reessayez."})
		}

		if time.Now().After(request.TokenExpiresAt) {
			logClientAccess(accessLogRepo, &request.ID, "client_auth", ipAddress, userAgent, false, "token_expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Votre lien d'acces a expire. Contactez-nous pour en obtenir un nouveau.", "code": "TOKEN_EXPIRED"})
		}

		if request.TokenUseCount > 200 {
			slog.Warn("client token use count high", "request_id", request.ID, "count", request.TokenUseCount, "ip", ipAddress)
		}

		requestID := request.ID
		go func() {
			if err := repo.UpdateTokenUsage(context.Background(), requestID); err != nil {
				slog.Error("update client token usage failed", "error", err, "request_id", requestID)
			}
		}()
		logClientAccess(accessLogRepo, &requestID, "client_auth", ipAddress, userAgent, true, "")

		c.Locals(ClientRequestKey, request)
		return c.Next()
	}
}

func GetClientRequest(c fiber.Ctx) *models.ClientRequest {
	request, _ := c.Locals(ClientRequestKey).(*models.ClientRequest)
	return request
}

func logClientAccess(repo *repository.AccessLogRepo, requestID *uuid.UUID, action string, ipAddress string, userAgent string, success bool, failureReason string) {
	if repo == nil {
		return
	}

	go func() {
		if err := repo.Log(context.Background(), repository.AccessLogEntry{
			RequestID:     requestID,
			Action:        action,
			IPAddress:     ipAddress,
			UserAgent:     userAgent,
			Success:       success,
			FailureReason: failureReason,
		}); err != nil {
			slog.Error("client access log failed", "error", err, "action", action)
		}
	}()
}
