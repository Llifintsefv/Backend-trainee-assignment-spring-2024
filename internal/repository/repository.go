package repository

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"context"
)

type OrganizationRepository interface {
	
}

type TenderRepository interface {
	CreateTender(context.Context, *model.Tender) (*model.Tender,error)
}

type UserRepository interface {
}
