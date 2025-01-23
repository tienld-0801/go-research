package auth_handler

import (
	auth_service "go-backend/app/services/auth"
	"go-backend/internal/constants"

	"github.com/gofiber/fiber/v3"
)

func Login(c fiber.Ctx) error {
	as := c.Locals(constants.AuthService).(auth_service.IAuthService)

	return as.Login(c)
}
