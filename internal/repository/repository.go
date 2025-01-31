package repository

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"context"
)


type TenderRepository interface {
	CreateTender(context.Context, *model.Tender) (*model.Tender,error)
	GetTenders(context.Context, int,int, []model.TenderServiceType) ([]model.Tender, error)
	GetTenderById(context.Context, string) (*model.Tender, error)
	GetCurrentUserTenders(context.Context, int,int, string) ([]model.Tender, error)
	UpdateTender(context.Context, *model.Tender) (*model.Tender, error)
	IsUserResponsibleForTender(context.Context, string, string) (bool, error)
}

type OrganizationRepository interface {
	GetOrganizationById(context.Context, string) (*model.Organization, error)
	IsUserResponsibleForOrganization(context.Context, string, string) (bool, error)
}

type UserRepository interface {
	GetUserById(context.Context, string) (*model.User, error)

}



type BidRepository interface {
	CreateBid(context.Context, *model.Bid) (*model.Bid, error)
}