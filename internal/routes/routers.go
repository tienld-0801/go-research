package routes

import (
	global_handler "go-backend/app/handlers/global"
	swagger_configs "go-backend/internal/configs/swagger"
	auth_router "go-backend/internal/routes/auth"
	user_router "go-backend/internal/routes/user"

	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", global_handler.Index)
	app.Get("/metrics", monitor.New())

	swagger_configs.RegisterSwagger(app)

	api := app.Group("/api")
	auth_router.AuthRouter(api)
	user_router.UserRouter(api)
}
