package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type BidHandler struct {
	service service.BidService
	logger *slog.Logger
}

func NewBidHandler(bidService service.BidService, logger *slog.Logger) *BidHandler {
	return &BidHandler{service: bidService, logger: logger}
}

func (h *BidHandler) CreateBid(c *fiber.Ctx) error {
	createBidRequest := new(model.CreateBidRequest)
	ctx := c.Context()
	err := c.BodyParser(createBidRequest)
	fmt.Println(createBidRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid request body"})
	}

	if createBidRequest.TenderID == "" || createBidRequest.Name == ""  || createBidRequest.Description == "" {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid request body"})
	}

	bid, err := h.service.CreateBid(ctx, createBidRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error creating bid", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error creating bid"})
	}
	return c.Status(fiber.StatusCreated).JSON(bid)


	
}	