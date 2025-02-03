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
	logger  *slog.Logger
}

type BidHandler interface {
	CreateBid(c *fiber.Ctx) error
	GetCurrentUserBids(c *fiber.Ctx) error
	GetTenderBids(c *fiber.Ctx) error
	GetBidStatus(c *fiber.Ctx) error
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

func (h *bidHandler) GetCurrentUserBids(c *fiber.Ctx) error {
	ctx := c.Context()

	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)
	user := c.Queries()["username"]

	if user == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "username is required"})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid offset"})
	}

	var bids []model.Bid
	bids, err := h.service.GetCurrentUserBids(ctx, limit, offset, user)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bids"})
	}
	return c.Status(fiber.StatusOK).JSON(bids)
}

func (h *bidHandler) GetTenderBids(c *fiber.Ctx) error {
	ctx := c.Context()

	limit := c.QueryInt("limit", 5)
	offset := c.QueryInt("offset", 0)
	username := c.Query("username")
	tenderID := c.Params("tenderId")

	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "username is required"})
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid offset"})
	}

	var bids []model.Bid
	bids, err := h.service.GetTenderBids(ctx, tenderID, limit, offset, username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bids"})
	}
	return c.Status(fiber.StatusOK).JSON(bids)
}

func (h *bidHandler) GetBidStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("bidId")
	username := c.Query("username")

	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "username is required"})
	}

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "bidId is required"})
	}
	status, err := h.service.GetBidStatus(ctx, id, username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bid status", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bid status"})
	}
	return c.Status(fiber.StatusOK).JSON(status)
}
