package model

import "time"


type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
	BidStatusRejected  BidStatus = "Rejected"
)


type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)


type CreateBidRequest struct {
	Name          string      `json:"name" validate:"required,max=255"` 
	Description   string      `json:"description" validate:"required,max=1000"` 
	Status        BidStatus   `json:"status" validate:"required,bidstatus"`    
	TenderID      string      `json:"tenderId" validate:"required"`         
	OrganizationID string      `json:"organizationId,omitempty" validate:"omitempty"` 
	CreatorUsername string      `json:"creatorUsername" validate:"required"`     
}


type Bid struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        BidStatus     `json:"status"`
	TenderID      string        `json:"tenderId"`
	AuthorType    BidAuthorType `json:"authorType"`
	AuthorID      string        `json:"authorId"`
	CreatorUsername string        `json:"creatorUsername"`
	Version       int           `json:"version"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}