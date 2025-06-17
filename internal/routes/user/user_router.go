package user_router

import (
	users_handler "go-backend/app/handlers/users"
	"go-backend/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

func UserRouter(app fiber.Router) {
	userGroup := app.Group("/user", middlewares.VerifyToken)
	userGroup.Get("/list", users_handler.GetAllUser)
	userGroup.Get("/:uuid", users_handler.GetUserByUUID)
	userGroup.Post("/create", users_handler.CreateUser)
	userGroup.Delete("/:uuid", users_handler.DeleteUser)
}
