package service

import "Backend-trainee-assignment-autumn-2024/internal/repository"

type TenderService interface {
}

type tenderService struct {
	TenderRepository repository.TenderRepository
	UserRepository   repository.UserRepository
	OrganizationRepository   repository.OrganizationRepository
}

func NewTenderService(tenderRepository repository.TenderRepository, userRepository repository.UserRepository, organizationRepository repository.OrganizationRepository) TenderService {
	return &tenderService{tenderRepository, userRepository, organizationRepository}
}