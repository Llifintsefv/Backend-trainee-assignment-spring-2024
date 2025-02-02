package service

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type TenderService interface {
	CreateTender(context.Context, *model.CreateTenderRequest) (*model.Tender, error)
	GetTenders(context.Context, int, int, []model.TenderServiceType) ([]model.Tender, error)
	GetTenderById(context.Context, string) (*model.Tender, error)
	GetCurrentUserTenders(context.Context, int, int, string) ([]model.Tender, error)
	GetTenderStatus(context.Context, string) (string, error)
	UpdateTenderStatus(context.Context, string, string, string) (*model.Tender, error)
	EditTender(context.Context, string, string, model.UpdateData) (*model.Tender, error)
	RollbackTenderVersion(context.Context, string, int) (*model.Tender, error)
}

type tenderService struct {
	TenderRepository       repository.TenderRepository
	OrganizationRepository repository.OrganizationRepository
	logger                 *slog.Logger
}

func NewTenderService(tenderRepository repository.TenderRepository, OrganizationRepository repository.OrganizationRepository, logger *slog.Logger) TenderService {
	return &tenderService{tenderRepository, OrganizationRepository, logger}
}

func (s *tenderService) CreateTender(ctx context.Context, createTenderRequest *model.CreateTenderRequest) (*model.Tender, error) {

	isResponsible, err := s.OrganizationRepository.IsUserResponsibleForOrganization(ctx, createTenderRequest.OrganizationID, createTenderRequest.CreatorUsername)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error checking user is responsible for tender", slog.Any("error", err))
		return nil, err
	}

	if !isResponsible {
		s.logger.ErrorContext(ctx, "User is not responsible for the tender", slog.Any("error", err))
		return nil, err
	}

	tender := &model.Tender{}

	tender.ID = uuid.NewString()
	tender.Name = createTenderRequest.Name
	tender.Description = createTenderRequest.Description
	tender.ServiceType = createTenderRequest.ServiceType
	tender.OrganizationID = createTenderRequest.OrganizationID
	tender.CreatorUsername = createTenderRequest.CreatorUsername
	tender.Version = 1
	tender.Status = model.TenderStatusCreated

	tender, err = s.TenderRepository.CreateTender(ctx, tender)
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

	tenders, err := s.TenderRepository.GetTenderByUsername(ctx, limit, offset, username)
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

func (s *tenderService) UpdateTenderStatus(ctx context.Context, id string, username string, status string) (*model.Tender, error) {

	isResponsible, err := s.TenderRepository.IsUserResponsibleForTender(ctx, id, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error checking user is responsible for tender", slog.Any("error", err))
		return nil, err
	}

	if !isResponsible {
		s.logger.ErrorContext(ctx, "User is not responsible for the tender", slog.Any("error", err))
		return nil, err
	}

	tender, err := s.TenderRepository.GetTenderById(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender status", slog.Any("error", err))
		return nil, err
	}

	if tender.Status == model.TenderStatus(status) {
		s.logger.ErrorContext(ctx, "Status is the same", slog.Any("error", err))
		return nil, err
	}

	tender.Status = model.TenderStatus(status)
	tender, err = s.TenderRepository.UpdateTender(ctx, tender)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error updating tender status", slog.Any("error", err))
		return nil, err
	}

	return tender, nil
}

func (s *tenderService) EditTender(ctx context.Context, id string, username string, updateData model.UpdateData) (*model.Tender, error) {

	tender, err := s.TenderRepository.GetTenderById(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender", slog.Any("error", err))
		return nil, err
	}

	if tender.CreatorUsername != username {
		s.logger.ErrorContext(ctx, "User is not the creator of the tender", slog.Any("error", err))
		return nil, err
	}
	if *updateData.Name != "" {
		tender.Name = *updateData.Name
	}

	if *updateData.Description != "" {
		tender.Description = *updateData.Description
	}

	if *updateData.ServiceType != "" {
		tender.ServiceType = model.TenderServiceType(*updateData.ServiceType)
	}

	tender, err = s.TenderRepository.UpdateTender(ctx, tender)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error updating tender", slog.Any("error", err))
		return nil, err
	}

	return tender, nil
}

func (s *tenderService) RollbackTenderVersion(ctx context.Context, id string, version int) (*model.Tender, error) {

	tender, err := s.TenderRepository.RollbackTenderVersion(ctx, id, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error rolling back tender version", slog.Any("error", err))
		return nil, err
	}
	return tender, nil
}
