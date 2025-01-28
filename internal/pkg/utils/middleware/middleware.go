package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware для имитации авторизации по Bearer токену 
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"reason": "Заголовок Authorization отсутствует",
		})
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"reason": "Неверный формат заголовка Authorization. Ожидается 'Bearer <токен>'",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" { 
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"reason": "Токен авторизации не предоставлен после 'Bearer '",
		})
	}
	c.Locals("token", token)
	return c.Next()
}
