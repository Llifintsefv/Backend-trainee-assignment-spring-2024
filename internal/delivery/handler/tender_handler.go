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
	logger *slog.Logger	
}

type TenderHandler interface {
	CreateTender(c *fiber.Ctx) error
	GetTenders(c *fiber.Ctx) error
	GetCurrentUserTenders(c *fiber.Ctx) error
	GetTenderStatus(c *fiber.Ctx) error
	UpdateTenderStatus(c *fiber.Ctx) error
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
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason:err.Error()})
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



	limit := getTendersRequest.Limit
	offset := getTendersRequest.Offset
	serviceTypes := getTendersRequest.ServiceTypes


	var tenders []model.Tender
	tenders, err := h.tenderService.GetTenders(ctx,limit,offset,serviceTypes)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)
}


func (h *tenderHandler) GetCurrentUserTenders(c *fiber.Ctx) error {
	ctx := c.Context()

	limitStr := c.Queries()["limit"]
	offsetStr := c.Queries()["offset"]
	user := c.Queries()["username"]

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing limit", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid limit"})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing offset", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{Reason: "Invalid offset"})
	}

	var tenders []model.Tender
	tenders, err = h.tenderService.GetCurrentUserTenders(ctx,limit,offset,user)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)

}


func (h *tenderHandler) GetTenderStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("tenderId")
	status, err := h.tenderService.GetTenderStatus(ctx,id)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tender status", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error getting tender status"})
	}
	return c.Status(fiber.StatusOK).JSON(status)
}

func (h *tenderHandler) UpdateTenderStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("tenderId")
	user := c.Queries()["username"]
	TargetStatus := c.Queries()["status"]

	tender,err := h.tenderService.UpdateTenderStatus(ctx,id,user,TargetStatus)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error updating tender status", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{Reason: "Error updating tender status"})
	}
	return c.Status(fiber.StatusOK).JSON(tender)

}

