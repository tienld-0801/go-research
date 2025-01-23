package user_router

import (
	users_handler "go-backend/app/handlers/users"
	"go-backend/internal/middlewares"

	"github.com/gofiber/fiber/v3"
)

func UserRouter(app *fiber.App) {
	userGroup := app.Group("/user", middlewares.VerifyToken)
	userGroup.Get("/", users_handler.GetAllUser)
	userGroup.Get("/:id", users_handler.GetUserById)
	userGroup.Post("/", users_handler.CreateUser)
	userGroup.Delete("/:id", users_handler.DeleteUser)
}
