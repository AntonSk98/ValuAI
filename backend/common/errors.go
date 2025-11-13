package common

import (
	"github.com/gofiber/fiber/v2"
)

// SendError is a helper funcion to send error responses in a consistent format.
func SendError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{"error": msg})
}
