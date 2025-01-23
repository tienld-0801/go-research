package users_service

import (
	users_model "go-backend/app/models/users"
	user_dto "go-backend/app/services/user/dto"
	exception_configs "go-backend/internal/configs/exception"
	"go-backend/internal/constants"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(c fiber.Ctx) error
	DeleteUser(c fiber.Ctx) error
	GetAllUser(c fiber.Ctx) error
	GetUserById(c fiber.Ctx) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) IUserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUser(c fiber.Ctx) error {
	users := new([]users_model.User)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := s.db.Find(&users).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to get users", nil)
	}

	var usersWithoutPassword []map[string]interface{}
	for _, user := range *users {
		userMap := map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
		usersWithoutPassword = append(usersWithoutPassword, userMap)
	}

	return jsonResponse(c, fiber.StatusOK, "Get all users successfully", usersWithoutPassword)
}

func (s *UserService) GetUserById(c fiber.Ctx) error {
	id := c.Params("id")
	user := new(users_model.User)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := s.db.First(&user, id).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "User not found", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Get user by id sucessfully", &user)
}

func (s *UserService) CreateUser(c fiber.Ctx) error {
	user := new(user_dto.UserDTO)
	cv := c.Locals(constants.CustomValidator).(*exception_configs.CustomValidator)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := c.Bind().Body(user); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", nil)
	}

	errors := cv.Validate(user)

	if len(errors) > 0 {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", errors)
	}

	var existingUser users_model.User

	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Email already exists", nil)
	} else if err != gorm.ErrRecordNotFound {
		return jsonResponse(c, fiber.StatusInternalServerError, "Database error", nil)
	}

	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to hash password", nil)
	}

	newUser := users_model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return jsonResponse(c, fiber.StatusCreated, "User created sucessfully", &fiber.Map{
		"id":   newUser.ID,
		"name": newUser.Name,
	})
}

func (s *UserService) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	user := new(user_dto.UserDTO)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := s.db.First(&user, id).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "User not found", nil)
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to delete user", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "User deleted sucessfully", nil)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
