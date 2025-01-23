package configs

import (
	register_service "go-backend/app/services"
	auth_service "go-backend/app/services/auth"
	users_service "go-backend/app/services/user"
	database_configs "go-backend/internal/configs/database"
	env_configs "go-backend/internal/configs/env"
	exception_configs "go-backend/internal/configs/exception"
	logger_configs "go-backend/internal/configs/log"
	"go-backend/internal/configs/version"
	"go-backend/internal/constants"
	"go-backend/internal/routes"
	"go-backend/internal/utils"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Server struct {
	App *fiber.App
	Db  *gorm.DB
}

func NewServer() *Server {
	env_configs.LoadEnv()
	app := fiber.New()

	db, err := database_configs.ConnectDatabase()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	s := &Server{
		App: app,
		Db:  db,
	}

	s.Config()

	return s
}

func (s *Server) Start() error {
	port := version.GetInfoVersion()

	if err := s.App.Listen(port, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		log.Fatal("\033[31m❌ Error starting app: \033[0m", err)
		return err
	}

	return nil
}

func (s *Server) Config() {
	RegisterServices(s)
	routes.RegisterRoutes(s.App)
}

func RegisterServices(s *Server) {
	db := s.Db

	authService := auth_service.NewAuthService(db)
	userService := users_service.NewUserService(db)

	s.App.Use(
		register_service.AddComponentContext(constants.UserService, userService),
		register_service.AddComponentContext(constants.AuthService, authService),
		register_service.AddComponentContext(constants.Db, db),
		utils.AddJSONResponse,
		exception_configs.SetCustomValidatorContext,
		logger_configs.Logger(),
	)

}
