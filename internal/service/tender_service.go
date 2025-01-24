package service

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type TenderService interface {
	CreateTender(context.Context, *model.CreateTenderRequest) (*model.Tender,error)
}

type tenderService struct {
	TenderRepository repository.TenderRepository
	UserRepository   repository.UserRepository
	OrganizationRepository   repository.OrganizationRepository
	logger *slog.Logger
}

func NewTenderService(tenderRepository repository.TenderRepository, userRepository repository.UserRepository, organizationRepository repository.OrganizationRepository, logger *slog.Logger) TenderService {
	return &tenderService{tenderRepository, userRepository, organizationRepository, logger}
}

func (s *tenderService) CreateTender(ctx context.Context, createTenderRequest *model.CreateTenderRequest) (*model.Tender, error) {
	tender := &model.Tender{}

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