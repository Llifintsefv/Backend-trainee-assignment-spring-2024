package model

import (
	"time"
)
type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"   
	TenderStatusPublished TenderStatus = "Published" 
	TenderStatusClosed    TenderStatus = "Closed"    
)

type TenderServiceType string

const (
	TenderServiceTypeConstruction TenderServiceType = "Construction"
	TenderServiceTypeDelivery     TenderServiceType = "Delivery"
	TenderServiceTypeManufacture  TenderServiceType = "Manufacture"
)



type Tender struct {
	ID             string            `json:"id"`
	Name           string            `json:"name" `
	Description    string            `json:"description" `
	ServiceType    TenderServiceType `json:"serviceType" `
	Status         TenderStatus      `json:"status" `
	OrganizationID string            `json:"organizationId" `
	CreatorUsername string            `json:"creatorUsername" `
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type CreateTenderRequest struct {
	Name         string            `json:"name" validate:"required,max=100"` 
	Description  string            `json:"description" validate:"required,max=1000"` 
	ServiceType  TenderServiceType `json:"serviceType" validate:"required,servicetype"` 
	OrganizationID string            `json:"organizationId" " validate:"required"`
	CreatorUsername string            `json:"creatorUsername" validate:"required"`
}

type GetTendersRequest struct {
	Limit        int               `query:"limit" validate:"min=1,max=100"`
	Offset       int               `query:"offset" validate:"min=0"`
	ServiceTypes []TenderServiceType `query:"service_type" validate:"dive,servicetype"`
}

type UpdateData struct {
    Name        *string `json:"name"`
    Description *string `json:"description"`
    ServiceType *string `json:"serviceType"`
}
   