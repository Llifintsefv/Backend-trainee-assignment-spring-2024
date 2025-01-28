package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/pkg/utils"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type TenderHandler struct {
	tenderService service.TenderService
	logger *slog.Logger	
}

func NewTenderHandler(tenderService service.TenderService, logger *slog.Logger) *TenderHandler {
	return &TenderHandler{tenderService: tenderService, logger: logger}
}


func (h *TenderHandler) CreateTender(c *fiber.Ctx) error {
	createTenderRequest := new(model.CreateTenderRequest)
	ctx := c.Context()
	err := c.BodyParser(createTenderRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid request body"})
	}
	
	if createTenderRequest.Name == "" || createTenderRequest.Description == "" || createTenderRequest.ServiceType == "" || createTenderRequest.OrganizationID == "" || createTenderRequest.CreatorUsername == "" {
		h.logger.ErrorContext(ctx, "Error parsing request body", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid request body"})
	}
	

	tender, err := h.tenderService.CreateTender(ctx, createTenderRequest)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error creating tender", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error creating tender"})
	}
	return c.Status(fiber.StatusCreated).JSON(tender)
	
}

func (h *TenderHandler) GetTenders(c *fiber.Ctx) error {
	ctx := c.Context()

	getTendersRequest := new(model.GetTendersRequest)

	if err := c.QueryParser(getTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Error parsing query parameters", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid query parameters"})
	}


	if err := utils.ValidateStruct(getTendersRequest); err != nil {
		h.logger.ErrorContext(ctx, "Validation error", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": err.Error()})
	}



	limit := getTendersRequest.Limit
	offset := getTendersRequest.Offset
	serviceTypes := getTendersRequest.ServiceTypes


	var tenders []model.Tender
	tenders, err := h.tenderService.GetTenders(ctx,limit,offset,serviceTypes)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)
}