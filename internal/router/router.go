package router

import (
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(tenderHandler *handler.TenderHandler,pingHandler *handler.PingHandler) *fiber.App {
	app := fiber.New()

	app.Get("/api/ping", pingHandler.Ping)

	return app
}