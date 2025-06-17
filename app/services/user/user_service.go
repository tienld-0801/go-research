package users_service

import (
	users_model "go-backend/app/models/users"
	user_dto "go-backend/app/services/user/dto"
	exception_configs "go-backend/internal/configs/exception"
	"go-backend/internal/constants"
	"go-backend/internal/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(c fiber.Ctx) error
	DeleteUser(c fiber.Ctx) error
	GetAllUser(c fiber.Ctx) error
	GetUserByUUID(c fiber.Ctx) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) IUserService {
	return &UserService{db: db}
}

// Handle get list of users
func (s *UserService) GetAllUser(c fiber.Ctx) error {
	users := new([]users_model.User)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := s.db.Find(&users).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to get users", nil)
	}

	// Prepared the user info response for each of user
	var usersWithoutPassword []map[string]interface{}
	for _, user := range *users {
		userMap := map[string]interface{}{
			"uuid":  user.UUID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}
		usersWithoutPassword = append(usersWithoutPassword, userMap)
	}

	return jsonResponse(c, fiber.StatusOK, "Get all users successfully", usersWithoutPassword)
}

// Handle get the user by uuid
func (s *UserService) GetUserByUUID(c fiber.Ctx) error {
	uuid := c.Params("uuid")
	user := new(users_model.User)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	err := s.db.First(&user, "uuid = ?", uuid).Error
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "User not found", nil)
	}

	// Prepared the user info response
	userResponse := users_model.UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return jsonResponse(c, fiber.StatusOK, "Get user by uuid successfully", &userResponse)
}

// Handle create a new user
func (s *UserService) CreateUser(c fiber.Ctx) error {
	user := new(user_dto.UserDTO)
	cv := c.Locals(constants.CustomValidator).(*exception_configs.CustomValidator)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := c.Bind().Body(user); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", nil)
	}

	// Validate the user info in request body
	errors := cv.Validate(user)
	if len(errors) > 0 {
		return jsonResponse(c, fiber.StatusBadRequest, "Params is not empty", errors)
	}

	// Check this user info is exited or not
	var existingUser users_model.User
	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Email already exists", nil)
	} else if err != gorm.ErrRecordNotFound {
		return jsonResponse(c, fiber.StatusInternalServerError, "Database error", nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to hash password", nil)
	}

	// Create user information to save into Database
	newUser := users_model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		Role:     user.Role,
	}
	if err := s.db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return jsonResponse(c, fiber.StatusCreated, "User created successfully", &fiber.Map{
		"uuid": newUser.UUID,
		"name": newUser.Name,
		"role": newUser.Role,
	})
}

// Handle delete user
func (s *UserService) DeleteUser(c fiber.Ctx) error {
	uuid := c.Params("uuid")
	user := new(users_model.User)
	jsonResponse := c.Locals(constants.JSONResponse).(func(c fiber.Ctx, status int, message string, data interface{}) error)

	if err := s.db.First(&user, "uuid = ?", uuid).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "User not found", nil)
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Failed to delete user", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}
