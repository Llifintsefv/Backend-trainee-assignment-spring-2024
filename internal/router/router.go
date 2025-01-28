package router

import (
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"
	"Backend-trainee-assignment-autumn-2024/internal/pkg/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(tenderHandler *handler.TenderHandler,pingHandler *handler.PingHandler,bidHandler *handler.BidHandler) *fiber.App {
	app := fiber.New()

	app.Get("/api/ping", pingHandler.Ping)

	api := app.Group("/",middleware.AuthMiddleware) // Имитация авторизации

	api.Post("/tender/new",tenderHandler.CreateTender)

	api.Get("/tenders",tenderHandler.GetTenders)

	api.Post("/bids/new",bidHandler.CreateBid)

	return app
}