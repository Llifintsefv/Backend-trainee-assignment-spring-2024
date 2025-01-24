package handler

import "Backend-trainee-assignment-autumn-2024/internal/service"

type TenderHandler struct {
	tenderService service.TenderService
}

func NewTenderHandler(tenderService service.TenderService) *TenderHandler {
	return &TenderHandler{tenderService: tenderService}
}

