package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)


type pingHandler struct {
	logger *slog.Logger
}

type PingHandler interface {
	Ping(c *fiber.Ctx) error
}

func NewPingHandler(logger *slog.Logger) PingHandler {
	return &pingHandler{logger: logger}
}
func (p *pingHandler)Ping(c *fiber.Ctx) error {
	return c.SendString("ok")
}