package main

import (
	"log"
	"os"
	"time"
	"valuai/auth"
	"valuai/common"
	core "valuai/core/state_engine"
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
	authService := auth.NewAuthenticationService(mailSender, authConfig)
	auth.InitAuthenticationController(app, authService)

	core.InitAnalysisFlowStateEngine("resources/analysis_flow.yml")

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
