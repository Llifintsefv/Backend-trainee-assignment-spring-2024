package router

import (
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(tenderHandler *handler.TenderHandler) *fiber.App {
	app := fiber.New()


	return app
}