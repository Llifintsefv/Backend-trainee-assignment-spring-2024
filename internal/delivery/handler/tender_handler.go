package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/pkg/utils"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type tenderHandler struct {
	tenderService service.TenderService
	logger        *slog.Logger
}

type TenderHandler interface {
	CreateTender(c *fiber.Ctx) error
	GetTenders(c *fiber.Ctx) error
	GetCurrentUserTenders(c *fiber.Ctx) error
	GetTenderStatus(c *fiber.Ctx) error
	UpdateTenderStatus(c *fiber.Ctx) error
	EditTender(c *fiber.Ctx) error
	RollbackTender(c *fiber.Ctx) error
}

func NewTenderHandler(tenderService service.TenderService, logger *slog.Logger) TenderHandler {
	return &tenderHandler{tenderService: tenderService, logger: logger}
}

func (h *tenderHandler) CreateTender(c *fiber.Ctx) error {
	createTenderRequest := new(model.CreateTenderRequest)
	ctx := c.Context()
	err := c.BodyParser(createTenderRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid request body"})
	}

	if err := utils.ValidateStruct(createTenderRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	tender, err := h.tenderService.CreateTender(ctx, createTenderRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error creating tender", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error creating tender"})
	}
	return c.Status(fiber.StatusCreated).JSON(tender)

}

func (h *tenderHandler) GetTenders(c *fiber.Ctx) error {
	ctx := c.Context()

	getTendersRequest := new(model.GetTendersRequest)

	if err := c.QueryParser(getTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(getTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	tenders, err := h.tenderService.GetTenders(ctx, getTendersRequest.Limit, getTendersRequest.Offset, getTendersRequest.ServiceTypes)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)
}

func (h *tenderHandler) GetCurrentUserTenders(c *fiber.Ctx) error {
	ctx := c.Context()

	getCurrentUserTendersRequest := new(model.GetCurrentUserTendersRequest)

	if err := c.QueryParser(getCurrentUserTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(getCurrentUserTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	tenders, err := h.tenderService.GetCurrentUserTenders(ctx, getCurrentUserTendersRequest.Limit, getCurrentUserTendersRequest.Offset, getCurrentUserTendersRequest.Username)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)

}

func (h *tenderHandler) GetTenderStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	getTenderStatusRequest := new(model.GetTenderStatusRequest)

	getTenderStatusRequest.TenderID = c.Params("tenderId")

	if err := utils.ValidateStruct(getTenderStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	status, err := h.tenderService.GetTenderStatus(ctx, getTenderStatusRequest.TenderID)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tender status", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tender status"})
	}
	return c.Status(fiber.StatusOK).JSON(status)
}

func (h *tenderHandler) UpdateTenderStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	updateTenderStatusRequest := new(model.UpdateTenderStatusRequest)

	updateTenderStatusRequest.TenderID = c.Params("tenderId")
	if err := c.QueryParser(updateTenderStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := utils.ValidateStruct(updateTenderStatusRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	tender, err := h.tenderService.UpdateTenderStatus(ctx, updateTenderStatusRequest.TenderID, updateTenderStatusRequest.Username, string(updateTenderStatusRequest.Status))
	if err != nil {
		h.logger.ErrorContext(ctx, "Error updating tender status", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error updating tender status"})
	}
	return c.Status(fiber.StatusOK).JSON(tender)

}

func (h *tenderHandler) EditTender(c *fiber.Ctx) error {
	ctx := c.Context()

	editTenderRequest := new(model.EditTenderRequest)
	editTenderRequest.TenderID = c.Params("tenderId")

	if err := c.QueryParser(editTenderRequest); err != nil {
		h.logger.Error("Error parsing query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid query parameters"})
	}

	if err := c.BodyParser(&editTenderRequest.UpdateData); err != nil {
		h.logger.Error("Error parsing request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid request body"})
	}

	if err := utils.ValidateStruct(editTenderRequest); err != nil {
		h.logger.Error("Validation error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	updatedTender, err := h.tenderService.EditTender(ctx, editTenderRequest.TenderID, editTenderRequest.Username, editTenderRequest.UpdateData)
	if err != nil {
		h.logger.Error("Error updating tender", "error", err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedTender)
}

func (h *tenderHandler) RollbackTender(c *fiber.Ctx) error {
	ctx := c.Context()
	rollbackTenderRequest := new(model.RollbackTenderRequest)

	rollbackTenderRequest.TenderID = c.Params("tenderId")
	rollbackTenderRequest.Version = c.Params("version")

	if err := utils.ValidateStruct(rollbackTenderRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: err.Error()})
	}

	version, err := strconv.Atoi(rollbackTenderRequest.Version)
	if err != nil {
		h.logger.Error("Error converting version to int", "error", err)
	}

	_, err = h.tenderService.RollbackTenderVersion(ctx, rollbackTenderRequest.TenderID, version)
	if err != nil {
		h.logger.Error("Error rolling back tender", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error rolling back tender"})
	}
	return c.Status(fiber.StatusOK).JSON(model.ErrorResponse{Reason: "Tender rollbacked successfully"})
}
