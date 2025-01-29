package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/pkg/utils"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type bidHandler struct {
	service service.BidService
	logger *slog.Logger
}

type BidHandler interface{
	CreateBid(c *fiber.Ctx) error
}

func NewBidHandler(bidService service.BidService, logger *slog.Logger) BidHandler {
	return &bidHandler{service: bidService, logger: logger}
}

func (h *bidHandler) CreateBid(c *fiber.Ctx) error {
	createBidRequest := new(model.CreateBidRequest)
	ctx := c.Context()
	err := c.BodyParser(createBidRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid request body"})
	}

	err = utils.ValidateStruct(createBidRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}
	
	bid, err := h.service.CreateBid(ctx, createBidRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error creating bid", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error creating bid"})
	}
	return c.Status(fiber.StatusCreated).JSON(bid)


	
}	