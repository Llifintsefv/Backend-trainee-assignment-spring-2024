package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
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