package model

import (
	"time"
)
type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"   
	TenderStatusPublished TenderStatus = "Publish" 
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
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type CreateTenderRequest struct {
	Name         string            `json:"name" `
	Description  string            `json:"description" `
	ServiceType  TenderServiceType `json:"serviceType" ` 
	OrganizationID string            `json:"organizationId" "`
	CreatorUsername string            `json:"creatorUsername" `
}

type GetTendersRequest struct {
	Limit        int               `query:"limit" validate:"min=1,max=100"`
	Offset       int               `query:"offset" validate:"min=0"`
	ServiceTypes []TenderServiceType `query:"service_type" validate:"dive,servicetype"`
}