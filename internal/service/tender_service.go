package service

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type TenderService interface {
	CreateTender(context.Context, *model.CreateTenderRequest) (*model.Tender,error)
	GetTenders(context.Context, int,int, []model.TenderServiceType) ([]model.Tender, error)
}

type tenderService struct {
	TenderRepository repository.TenderRepository
	logger *slog.Logger
}

func NewTenderService(tenderRepository repository.TenderRepository, logger *slog.Logger) TenderService {
	return &tenderService{tenderRepository, logger}
}

func (s *tenderService) CreateTender(ctx context.Context, createTenderRequest *model.CreateTenderRequest) (*model.Tender, error) {
	tender := &model.Tender{}

	if createTenderRequest.ServiceType != "Construction" && createTenderRequest.ServiceType != "Delivery" && createTenderRequest.ServiceType != "Manufacture" {
		s.logger.ErrorContext(ctx, "Error creating tender", slog.Any("error", "Invalid service type"))
		return nil, fmt.Errorf("Invalid service type")
	}


	tender.ID = uuid.NewString()
	tender.Name = createTenderRequest.Name
	tender.Description = createTenderRequest.Description
	tender.ServiceType = createTenderRequest.ServiceType
	tender.OrganizationID = createTenderRequest.OrganizationID
	tender.CreatorUsername = createTenderRequest.CreatorUsername
	tender.Version = 1
	tender.Status = "CREATED"

	tender, err := s.TenderRepository.CreateTender(ctx, tender)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error creating tender", slog.Any("error", err))
		return nil, err
	}

	return tender, nil
}

func (s *tenderService) GetTenders(ctx context.Context, limit int, offset int, serviceTypes []model.TenderServiceType) ([]model.Tender, error) {

	tenders, err := s.TenderRepository.GetTenders(ctx, limit, offset, serviceTypes)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return nil, err
	}

	return tenders, nil
}