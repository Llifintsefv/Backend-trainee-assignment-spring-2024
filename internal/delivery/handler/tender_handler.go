package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"log/slog"
	"strconv"

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

	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	queryArgs := c.Context().QueryArgs()
serviceTypesQuery := queryArgs.PeekMulti("service_type")

	limit,err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error getting tenders"})
	}

	if limit < 1 || limit > 100 {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", "Invalid limit"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid limit"})
	}

	offset,err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error getting tenders"})
	}

	if offset < 0 {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", "Invalid offset"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"reason": "Invalid offset"})
	}

	serviceTypes := make([]model.TenderServiceType, 0, len(serviceTypesQuery))
for _, st := range serviceTypesQuery {
	tenderServiceType := model.TenderServiceType(string(st)) // Конвертация []byte в string
	if !model.IsValidServiceType(tenderServiceType) {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", "Invalid service type"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"reason": "Invalid service type",
		})
	}
	
	serviceTypes = append(serviceTypes, tenderServiceType)

}
	var tenders []model.Tender
	tenders, err = h.tenderService.GetTenders(ctx,limit,offset,serviceTypes)
	if err != nil {
		h.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"reason": "Error getting tenders"})
	}
	return c.Status(fiber.StatusOK).JSON(tenders)

}