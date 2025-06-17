package users_handler

import (
	register_service "go-backend/app/services"
	users_service "go-backend/app/services/user"
	"go-backend/internal/constants"

	"github.com/gofiber/fiber/v3"
)

func GetAllUser(c fiber.Ctx) error {
	us := register_service.GetServiceContext[users_service.IUserService](c, constants.UserService)
	return us.GetAllUser(c)
}

func GetUserByUUID(c fiber.Ctx) error {
	us := register_service.GetServiceContext[users_service.IUserService](c, constants.UserService)
	return us.GetUserByUUID(c)
}

func CreateUser(c fiber.Ctx) error {
	us := register_service.GetServiceContext[users_service.IUserService](c, constants.UserService)
	return us.CreateUser(c)
}

func DeleteUser(c fiber.Ctx) error {
	us := register_service.GetServiceContext[users_service.IUserService](c, constants.UserService)
	return us.DeleteUser(c)
}
