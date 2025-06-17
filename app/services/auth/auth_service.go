package auth_service

import (
	users_model "go-backend/app/models/users"
	auth_dto "go-backend/app/services/auth/dto"
	exception_configs "go-backend/internal/configs/exception"
	"go-backend/internal/constants"
	"go-backend/internal/middlewares"
	"go-backend/internal/utils"

	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IAuthService interface {
	Login(c fiber.Ctx) error
	RefreshToken(c fiber.Ctx) error
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) IAuthService {
	return &AuthService{db: db}
}

// Handle Login
func (s *AuthService) Login(c fiber.Ctx) error {
	auth := new(auth_dto.AuthDTO)
	cv := c.Locals(constants.CustomValidator).(*exception_configs.CustomValidator)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := c.Bind().Body(auth); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", nil)
	}

	// Validate auth
	errors := cv.Validate(auth)
	if len(errors) > 0 {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not valid", errors)
	}

	// Find this user in the Database by field "email"
	var user users_model.User
	if err := s.db.Where("email = ?", auth.Email).First(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusUnauthorized, "Email not found", nil)
	}

	// Authenticate password
	if !utils.CheckPassword(user.Password, auth.Password) {
		return jsonResponse(c, fiber.StatusUnauthorized, "Password is Correct", nil)
	}

	// Start generate AccessToken and RefreshToken
	accessToken, refreshToken, accessClaims, err := generateJWT(user)
	if err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to generate token", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Login successfully", fiber.Map{
		"info":         accessClaims,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// Handle Refresh Token
func (s *AuthService) RefreshToken(c fiber.Ctx) error {
	bodyRequest := new(auth_dto.RefreshTokenDTO)
	cv := c.Locals(constants.CustomValidator).(*exception_configs.CustomValidator)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := c.Bind().Body(bodyRequest); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", nil)
	}

	// Validate auth
	errors := cv.Validate(bodyRequest)
	if len(errors) > 0 {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not valid", errors)
	}

	// Parse the token to get claims info
	claims := &middlewares.Claims{}
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	parsedToken, err := jwt.ParseWithClaims(bodyRequest.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check the token is valid or not
	if err != nil || !parsedToken.Valid {
		return jsonResponse(c, fiber.StatusBadRequest, "invalid refresh token", nil)
	}

	// Check the claims gotten from token is valid or not
	claimsMap, ok := parsedToken.Claims.(*middlewares.Claims)
	if !ok {
		return jsonResponse(c, fiber.StatusBadRequest, "invalid refresh token claims", nil)
	}

	// Find this user in the Database by field "uuid"
	var user users_model.User
	if err := s.db.First(&user, "uuid = ?", claimsMap.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Start generate AccessToken from RefreshToken
	accessToken, accessClaims, err := generateJWTFromRefreshToken(user, claimsMap.RefreshTokenExpiresAt)
	if err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to refresh new token", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Refresh token successfully", fiber.Map{
		"info":        accessClaims,
		"accessToken": accessToken,
	})
}

// Handle generate JWT
func generateJWT(user users_model.User) (string, string, *middlewares.Claims, error) {
	// Access Token Claims
	accessClaims := &middlewares.Claims{
		UserID:                user.UUID,
		Email:                 user.Email,
		Role:                  user.Role,
		TokenType:             "Bearer",
		RefreshTokenExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// Generate AccessToken from Claims
	accessTokenGenerated := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	// Refresh Token Claims
	refreshClaims := &middlewares.Claims{
		UserID:    user.UUID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "Refresh Token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// Generate RefreshToken from Claims
	refreshTokenGenerated := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign the secretKey into AccessToken
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	accessToken, err := accessTokenGenerated.SignedString(secretKey)
	if err != nil {
		return "", "", accessClaims, err
	}

	// Sign the secretKey into RefreshToken
	refreshToken, err := refreshTokenGenerated.SignedString(secretKey)
	if err != nil {
		return "", "", refreshClaims, err
	}

	return accessToken, refreshToken, accessClaims, nil
}

// Handle generate JWT from RefreshToken
func generateJWTFromRefreshToken(user users_model.User, refreshTokenExpiresAt int64) (string, *middlewares.Claims, error) {
	// Access Token Claims
	accessClaims := &middlewares.Claims{
		UserID:                user.UUID,
		Email:                 user.Email,
		Role:                  user.Role,
		TokenType:             "Bearer",
		RefreshTokenExpiresAt: refreshTokenExpiresAt, // Current RefreshToken
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Renew for 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// Generate RefreshToken from Claims
	accessTokenGenerated := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	// Sign the secretKey into AccessToken
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	accessToken, err := accessTokenGenerated.SignedString(secretKey)
	if err != nil {
		return "", accessClaims, err
	}

	return accessToken, accessClaims, nil
}
