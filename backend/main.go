package main

import (
	"log"
	"os"
	"time"
	"valuai/auth"
	"valuai/common"
	"valuai/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		common.LogInfo("No .env file found, using environment variables.")
	}

}

func main() {
	app := fiber.New()

	mailSender := mail.InitMailSender()

	authConfig := auth.InitAuthConfig(os.Getenv("JWT_SECRET"), time.Hour*24)
	// authMiddleware := auth.InitAuthMiddleware(authConfig)
	authService := auth.NewAuthenticationService(mailSender, authConfig)

	auth.InitAuthenticationController(app, authService)

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
