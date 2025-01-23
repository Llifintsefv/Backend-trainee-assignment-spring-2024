package model

import "time"
type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "CREATED"   
	TenderStatusPublished TenderStatus = "PUBLISHED" 
	TenderStatusClosed    TenderStatus = "CLOSED"    
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