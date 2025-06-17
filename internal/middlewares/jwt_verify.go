package middlewares

import (
	"fmt"
	"go-backend/internal/constants"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID                string `json:"uuid"`
	Email                 string `json:"email"`
	Role                  int    `json:"role"`
	TokenType             string `json:"token_type"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
	jwt.RegisteredClaims
}

func VerifyToken(c fiber.Ctx) error {
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)
	tokenString := c.Get("Authorization")

	// Check empty access token
	if tokenString == "" {
		return jsonResponse(c, fiber.StatusUnauthorized, "Token is required", nil)
	}

	// Check token type
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &Claims{}
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	// Check error parser token
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return jsonResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Check valid token
	if !parsedToken.Valid {
		fmt.Println("Invalid token:", parsedToken)
		return jsonResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Set the user ID in the Fiber context
	c.Locals(constants.VerifyToken, claims.UserID)

	return c.Next()
}
