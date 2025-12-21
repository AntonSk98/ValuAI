package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthenticationConfig holds configuration for JWT authentication
type AuthMiddleware struct {
	authConfig *AuthenticationConfig
}

// InitAuthMiddleware initializes the AuthMiddleware with the given configuration.
func InitAuthMiddleware(cfg *AuthenticationConfig) *AuthMiddleware {
	return &AuthMiddleware{
		authConfig: cfg,
	}
}

// Protect is a Fiber middleware that validates the JWT and stores claims in the request context.
func (am *AuthMiddleware) Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		// Expecting "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return fiber.ErrUnauthorized
		}

		// Parse and validate JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(am.authConfig.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}
		if iss, ok := claims["iss"].(string); !ok || iss != issuer {
			return fiber.ErrUnauthorized
		}

		// Store claims in request context for handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}
