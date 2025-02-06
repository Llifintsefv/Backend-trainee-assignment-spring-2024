package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/pkg/utils"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"errors"
	"log/slog"
	"strconv"

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
	UpdateBidStatus(c *fiber.Ctx) error
	EditBid(c *fiber.Ctx) error
	SubmitBidDecision(c *fiber.Ctx) error
	AddBidFeedback(c *fiber.Ctx) error
	RollbackBidVersion(c *fiber.Ctx) error
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
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error creating bid"})
	}
	return c.Status(fiber.StatusCreated).JSON(bid)
}

func (h *bidHandler) GetCurrentUserBids(c *fiber.Ctx) error {
	ctx := c.Context()

	getCurrentUserBidsRequest := new(model.GetCurrentUserBidsRequest)

	if err := c.QueryParser(getCurrentUserBidsRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(getCurrentUserBidsRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	bids, err := h.service.GetCurrentUserBids(ctx, getCurrentUserBidsRequest.Limit, getCurrentUserBidsRequest.Offset, getCurrentUserBidsRequest.Username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bids"})
	}
	return c.Status(fiber.StatusOK).JSON(bids)
}

func (h *bidHandler) GetTenderBids(c *fiber.Ctx) error {
	ctx := c.Context()

	getTenderBidsRequest := new(model.GetTenderBidsRequest)
	getTenderBidsRequest.TenderID = c.Params("tenderId")

	if err := c.QueryParser(getTenderBidsRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(getTenderBidsRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	bids, err := h.service.GetTenderBids(ctx, getTenderBidsRequest.TenderID, getTenderBidsRequest.Limit, getTenderBidsRequest.Offset, getTenderBidsRequest.Username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrTenderNotFound) || errors.Is(err, model.ErrBidNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bids"})
	}
	return c.Status(fiber.StatusOK).JSON(bids)
}

func (h *bidHandler) GetBidStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	getBidStatusRequest := new(model.GetBidStatusRequest)

	getBidStatusRequest.BidID = c.Params("bidId")

	if err := c.QueryParser(getBidStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(getBidStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	status, err := h.service.GetBidStatus(ctx, getBidStatusRequest.BidID, getBidStatusRequest.Username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting bid status", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting bid status"})
	}
	return c.Status(fiber.StatusOK).JSON(status)
}

func (h *bidHandler) UpdateBidStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	updateBidStatusRequest := new(model.UpdateBidStatusRequest)

	updateBidStatusRequest.BidID = c.Params("bidId")

	if err := c.QueryParser(updateBidStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(updateBidStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	status, err := h.service.UpdateBidStatus(ctx, updateBidStatusRequest.BidID, updateBidStatusRequest.Username, string(updateBidStatusRequest.Status))
	if err != nil {
		h.logger.ErrorContext(ctx, "Error updating bid status", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrBidNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error updating bid status"})
	}
	return c.Status(fiber.StatusOK).JSON(status)
}

func (h *bidHandler) EditBid(c *fiber.Ctx) error {
	ctx := c.Context()

	editBidRequest := new(model.EditBidRequest)
	editBidRequest.BidID = c.Params("bidId")

	if err := c.QueryParser(editBidRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := c.BodyParser(&editBidRequest.UpdateData); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid request body"})
	}

	if err := utils.ValidateStruct(editBidRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	bid, err := h.service.EditBid(ctx, editBidRequest.BidID, editBidRequest.Username, editBidRequest.UpdateData)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error editing bid", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrBidNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error editing bid"})
	}
	return c.Status(fiber.StatusOK).JSON(bid)
}

func (h *bidHandler) RollbackBidVersion(c *fiber.Ctx) error {
	ctx := c.Context()
	rollbackBidRequest := new(model.RollbackBidRequest)

	rollbackBidRequest.BidID = c.Params("bidId")
	rollbackBidRequest.Version = c.Params("version")

	if err := c.QueryParser(rollbackBidRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(rollbackBidRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	version, err := strconv.Atoi(rollbackBidRequest.Version)
	if err != nil {
		h.logger.Error("Error converting version to int", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid version parameter"})
	}

	bid, err := h.service.RollbackBidVersion(ctx, rollbackBidRequest.BidID, rollbackBidRequest.Username, version)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error rolling back bid", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(bid)
}

func (h *bidHandler) SubmitBidDecision(c *fiber.Ctx) error {
	ctx := c.Context()
	submitBidDecisionRequest := new(model.SubmitBidDecisionRequest)

	submitBidDecisionRequest.BidID = c.Params("bidId")

	if err := c.QueryParser(submitBidDecisionRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(submitBidDecisionRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	bid, err := h.service.SubmitBidDecision(ctx, submitBidDecisionRequest.BidID, submitBidDecisionRequest.Username, string(submitBidDecisionRequest.Decision))
	if err != nil {
		h.logger.ErrorContext(ctx, "Error submitting bid decision", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrBidNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrDecisionSubmit) {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error submitting bid decision"})
	}
	return c.Status(fiber.StatusOK).JSON(bid)
}


func (h *bidHandler) AddBidFeedback(c *fiber.Ctx) error{
	ctx := c.Context()
	addBidFeedbackRequest := new(model.AddBidFeedbackRequest)
	addBidFeedbackRequest.BidID = c.Params("bidId")

	if err := c.QueryParser(addBidFeedbackRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}
	
	if err := utils.ValidateStruct(addBidFeedbackRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	bid,err := h.service.AddBidFeedback(ctx, addBidFeedbackRequest.BidID, addBidFeedbackRequest.Username, addBidFeedbackRequest.Feedback)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error adding bid feedback", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		if errors.Is(err, model.ErrBidNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{Reason: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error adding bid feedback"})
	}
	return c.Status(fiber.StatusOK).JSON(bid)

}