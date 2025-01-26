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
	Name          string      `json:"name" `
	Description   string      `json:"description" `
	Status        BidStatus   `json:"status" `
	TenderID      string      `json:"tenderId"`    
	OrganizationID string      `json:"organizationId,omitempty,uuid"`     
	CreatorUsername string      `json:"creatorUsername" `
	AuthorType    BidAuthorType `json:"authorType"` 
	AuthorID      string      `json:"authorId" `     
}


type Bid struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        BidStatus     `json:"status"`
	TenderID      string        `json:"tenderId"`
	AuthorType    BidAuthorType `json:"authorType"`
	AuthorID      string        `json:"authorId"`
	Version       int           `json:"version"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}