package auth

import (
	"fmt"
	"time"
	"valuai/common"
	ratelimiter "valuai/rate_limiter"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles OTP-related HTTP endpoints.
type AuthController struct {
	authService *AuthenticationService
}

// InitAuthenticationController initializes the authentication controller with routes.
func InitAuthenticationController(app *fiber.App, authService *AuthenticationService) {
	controller := &AuthController{authService: authService}
	ipRateLimiter := ratelimiter.PerIPRateLimiter(5, 1*time.Minute)
	globalRateLimiter := ratelimiter.GlobalRateLimiter(1000, 1*time.Minute)
	authGroup := app.Group("/auth/otp")
	// Middlewares are executed in the order they are provided.
	authGroup.Post("/issue", globalRateLimiter, ipRateLimiter, controller.issueOtp)
	authGroup.Post("/verify", globalRateLimiter, ipRateLimiter, controller.verifyOtp)
}

func (ac *AuthController) issueOtp(c *fiber.Ctx) error {
	var cmd IssueOtpCommand
	if err := c.BodyParser(&cmd); err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrBadRequest.Code, "invalid request body")
	}

	if err := cmd.Validate(); err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrBadRequest.Code, fmt.Sprintf("validation failed for command: %s", err.Error()))
	}

	if err := ac.authService.GenerateOtp(cmd); err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrInternalServerError.Code, "failed to issue otp")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ac *AuthController) verifyOtp(c *fiber.Ctx) error {
	var cmd VerifyOtpCommand
	if err := c.BodyParser(&cmd); err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrBadRequest.Code, "invalid request body")

	}

	if err := cmd.Validate(); err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrBadRequest.Code, fmt.Sprintf("validation failed for command: %s", err.Error()))
	}

	token, err := ac.authService.VerifyOtp(cmd)
	if err != nil {
		common.LogError(err)
		return common.SendError(c, fiber.ErrBadRequest.Code, "otp validation failed")
	}

	return c.Status(fiber.StatusOK).JSON(token)
}
