package middleware

import (
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS(allowedOrigins []string) cors.Config {
	return cors.Config{
		AllowOrigins:  allowedOrigins,
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"X-Total-Count"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}
}
