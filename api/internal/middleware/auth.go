package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AdminJWT(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		tokenValue, err := extractBearerToken(auth)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		return c.Next()
	}
}

func extractBearerToken(value string) (string, error) {
	if !strings.HasPrefix(value, "Bearer ") {
		return "", errors.New("missing bearer token")
	}
	trimmed := strings.TrimSpace(strings.TrimPrefix(value, "Bearer "))
	if trimmed == "" {
		return "", errors.New("empty bearer token")
	}
	return trimmed, nil
}
