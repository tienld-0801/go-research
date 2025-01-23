package auth_service

import (
	users_model "go-backend/app/models/users"
	auth_dto "go-backend/app/services/auth/dto"
	exception_configs "go-backend/internal/configs/exception"
	"go-backend/internal/constants"

	"log"
	"os"

	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	Login(c fiber.Ctx) error
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) IAuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(c fiber.Ctx) error {
	auth := new(auth_dto.AuthDTO)
	cv := c.Locals(constants.CustomValidator).(*exception_configs.CustomValidator)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := c.Bind().Body(auth); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", nil)
	}

	errors := cv.Validate(auth)
	if len(errors) > 0 {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not valid", errors)
	}

	var user users_model.User

	if err := s.db.Where("email = ?", auth.Email).First(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusUnauthorized, "Email not found", nil)
	}

	if !CheckPassword(user.Password, auth.Password) {
		return jsonResponse(c, fiber.StatusUnauthorized, "Password is Correct", nil)
	}

	token, clams, err := generateJWT(user)

	if err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to generate token", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Login successfully", fiber.Map{
		"info":  clams,
		"token": token,
	})
}

func CheckPassword(storedPassword, inputPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		log.Println("Password comparison failed:", err)
		return false
	}
	return true
}

func generateJWT(user users_model.User) (string, jwt.MapClaims, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", claims, err
	}

	return signedToken, claims, nil
}
