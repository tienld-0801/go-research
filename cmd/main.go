package main

import (
	configs "go-backend/internal/configs"

	_ "go-backend/docs"
)

// @title Go Backend API
// @version 1.0
// @description API documentation for Go Backend application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4000
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	server := configs.NewServer()
	server.Start()
}
