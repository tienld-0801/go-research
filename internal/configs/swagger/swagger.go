package swagger_configs

import (
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	BasePath string
	FilePath string
	Path     string
	Title    string
	CacheAge int
}

func DefaultConfig() Config {
	return Config{
		BasePath: "/",
		FilePath: "docs/swagger.json",
		Path:     "swagger",
		Title:    "Go Backend API Documentation",
		CacheAge: 0,
	}
}

func RegisterSwagger(app *fiber.App) {
	cfg := DefaultConfig()

	app.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})

	app.Use(swaggerui.New(swaggerui.Config{
		BasePath: cfg.BasePath,
		FilePath: cfg.FilePath,
		Path:     cfg.Path,
		Title:    cfg.Title,
		CacheAge: cfg.CacheAge,
	}))
}

func RegisterSwaggerWithConfig(app *fiber.App, cfg Config) {

	app.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})

	app.Use(swaggerui.New(swaggerui.Config{
		BasePath: cfg.BasePath,
		FilePath: cfg.FilePath,
		Path:     cfg.Path,
		Title:    cfg.Title,
		CacheAge: cfg.CacheAge,
	}))
}
