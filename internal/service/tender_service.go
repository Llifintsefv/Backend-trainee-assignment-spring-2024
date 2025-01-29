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
	GetTenders(context.Context, int,int, []model.TenderServiceType) ([]model.Tender, error)
	GetTenderById(context.Context, string) (*model.Tender, error)
	GetCurrentUserTenders(context.Context, int,int, string) ([]model.Tender, error)
	GetTenderStatus(context.Context, string) (string, error)
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



	tender.ID = uuid.NewString()
	tender.Name = createTenderRequest.Name
	tender.Description = createTenderRequest.Description
	tender.ServiceType = createTenderRequest.ServiceType
	tender.OrganizationID = createTenderRequest.OrganizationID
	tender.CreatorUsername = createTenderRequest.CreatorUsername
	tender.Version = 1
	tender.Status = model.TenderStatusCreated

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


func (s *tenderService) GetTenderById(ctx context.Context, id string) (*model.Tender, error) {
	tender, err := s.TenderRepository.GetTenderById(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender by id", slog.Any("error", err))
		return nil, err
	}
	return tender, nil
}


func (s *tenderService) GetCurrentUserTenders(ctx context.Context, limit int, offset int, username string) ([]model.Tender, error) {

	tenders, err := s.TenderRepository.GetCurrentUserTenders(ctx, limit, offset, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return nil, err
	}

	return tenders, nil
}


func (s *tenderService) GetTenderStatus(ctx context.Context, id string) (string, error) {

	tender, err := s.TenderRepository.GetTenderById(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender status", slog.Any("error", err))
		return "", err
	}
	return string(tender.Status), nil
	}