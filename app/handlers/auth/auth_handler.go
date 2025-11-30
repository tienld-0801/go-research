package auth_handler

import (
	auth_service "go-backend/app/services/auth"
	"go-backend/internal/constants"

	"github.com/gofiber/fiber/v3"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" example:"admin@gmail.com"`
	Password string `json:"password" example:"admin123"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Successfully authenticated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /auth/login [post]
func Login(c fiber.Ctx) error {
	as := c.Locals(constants.AuthService).(auth_service.IAuthService)

	return as.Login(c)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "Successfully refreshed token"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /auth/refreshToken [post]
func RefreshToken(c fiber.Ctx) error {
	as := c.Locals((constants.AuthService)).(auth_service.IAuthService)

	return as.RefreshToken(c)
}
