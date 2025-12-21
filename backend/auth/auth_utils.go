package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetClaim returns a specific claim as a pointer to type T, or nil if missing or wrong type.
func GetClaim[T any](c *fiber.Ctx, key string) *T {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return nil
	}

	value, exists := claims[key]
	if !exists {
		return nil
	}

	typedValue, ok := value.(T)
	if !ok {
		return nil
	}

	return &typedValue
}
