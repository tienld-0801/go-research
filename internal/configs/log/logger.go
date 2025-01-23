package logger_configs

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:       "\u001b[36m${time}\u001b[0m \u001b[32m${ip}\u001b[0m \u001b[33m${status}\u001b[0m -\u001b[33m${latency}\u001b[0m \u001b[34m${method}\u001b[0m \u001b[31m${path}\u001b[0m ${error}\n",
		TimeInterval: 500 * time.Millisecond,
		TimeZone:     "Local",
	})
}
