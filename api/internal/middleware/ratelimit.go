package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func GlobalRateLimit() limiter.Config {
	return limiter.Config{Max: 100, Expiration: time.Minute}
}

func ContactRateLimit() limiter.Config {
	return limiter.Config{Max: 5, Expiration: time.Hour}
}

func ClientRequestRateLimit() limiter.Config {
	return limiter.Config{Max: 3, Expiration: time.Hour}
}

func ClientRequestNewLinkRateLimit() limiter.Config {
	return limiter.Config{Max: 2, Expiration: time.Hour}
}

func ClientMessageRateLimit() limiter.Config {
	return limiter.Config{Max: 10, Expiration: time.Hour}
}
