package register_service

import (
	"github.com/gofiber/fiber/v3"
)

func AddComponentContext(key string, value interface{}) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Locals(key, value)
		return c.Next()
	}
}

func GetServiceContext[T any](c fiber.Ctx, key string) T {
	service, ok := c.Locals(key).(T)
	if !ok {
		panic("Invalid service type or service not found in context")
	}
	return service
}
