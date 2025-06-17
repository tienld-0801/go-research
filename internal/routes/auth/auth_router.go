package auth_router

import (
	auth_handler "go-backend/app/handlers/auth"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(app fiber.Router) {
	authGroup := app.Group("/auth")
	authGroup.Post("/login", auth_handler.Login)
	authGroup.Post("/refreshToken", auth_handler.RefreshToken)
}
