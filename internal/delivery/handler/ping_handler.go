package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)


type PingHandler struct {
	logger *slog.Logger
}

func NewPingHandler(logger *slog.Logger) *PingHandler {
	return &PingHandler{logger: logger}
}
func (p *PingHandler)Ping(c *fiber.Ctx) error {
	return c.SendString("ok")
}