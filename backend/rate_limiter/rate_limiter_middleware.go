package ratelimiter

import (
	"time"
	"valuai/common"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// PerIPRateLimiter returns a Fiber handler that limits requests per IP.
func PerIPRateLimiter(maxRequest int, replenishRate time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: replenishRate,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return common.SendError(c, fiber.ErrTooManyRequests.Code, "rate limit exceeded, try again later")
		},
	})
}

// GlobalRateLimiter returns a Fiber handler that limits requests globally across all clients.
// It uses a constant key so the limiter counts all requests together.
func GlobalRateLimiter(maxRequest int, replenishRate time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: replenishRate,
		KeyGenerator: func(c *fiber.Ctx) string {
			// single shared key -> global limit
			return "__global__"
		},
		LimitReached: func(c *fiber.Ctx) error {
			return common.SendError(c, fiber.ErrTooManyRequests.Code, "global rate limit exceeded, try again later")
		},
	})
}
