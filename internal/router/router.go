package router

import (
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(tenderHandler *handler.TenderHandler,pingHandler *handler.PingHandler,bidHandler *handler.BidHandler) *fiber.App {
	app := fiber.New()

	app.Get("/api/ping", pingHandler.Ping)

	app.Post("/tender/new",tenderHandler.CreateTender)

	app.Get("/tenders",tenderHandler.GetTenders)

	app.Post("/bids/new",bidHandler.CreateBid)

	return app
}