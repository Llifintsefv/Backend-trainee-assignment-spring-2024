package handler

import (
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"log/slog"
)

type TenderHandler struct {
	tenderService service.TenderService
	logger *slog.Logger	
}

func NewTenderHandler(tenderService service.TenderService, logger *slog.Logger) *TenderHandler {
	return &TenderHandler{tenderService: tenderService, logger: logger}
}

