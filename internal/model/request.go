package model

type GetCurrentUserTendersRequest struct {
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	Offset   int    `query:"offset" validate:"min=0"`
	Username string `query:"username" validate:"required"`
}

type GetTenderStatusRequest struct {
	TenderID string `params:"tenderId" validate:"required"`
}

type UpdateTenderStatusRequest struct {
	TenderID string       `params:"tenderId" validate:"required"`
	Username string       `query:"username" validate:"required"`
	Status   TenderStatus `query:"status" validate:"required,oneof=Created Published Closed"`
}

type EditTenderRequest struct {
	TenderID   string     `params:"tenderId" validate:"required"`
	Username   string     `query:"username" validate:"required"`
	UpdateData UpdateData `json:"updateData" validate:"required"`
}

type RollbackTenderRequest struct {
	TenderID string `params:"tenderId" validate:"required"`
	Version  string `params:"version" validate:"required,number,min=1"`
}

type CreateTenderRequest struct {
	Name            string            `json:"name" validate:"required,max=100"`
	Description     string            `json:"description" validate:"required,max=1000"`
	ServiceType     TenderServiceType `json:"serviceType" validate:"required,servicetype"`
	OrganizationID  string            `json:"organizationId" validate:"required"`
	CreatorUsername string            `json:"creatorUsername" validate:"required"`
}

type GetTendersRequest struct {
	Limit        int                 `query:"limit" validate:"min=1,max=100"`
	Offset       int                 `query:"offset" validate:"min=0"`
	ServiceTypes []TenderServiceType `query:"service_type" validate:"dive,servicetype"`
}

type UpdateData struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	ServiceType *string `json:"serviceType" validate:"omitempty,servicetype"`
}

type GetCurrentUserBidsRequest struct {
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	Offset   int    `query:"offset" validate:"min=0"`
	Username string `query:"username" validate:"required"`
}

type GetTenderBidsRequest struct {
	TenderID string `params:"tenderId" validate:"required"`
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	Offset   int    `query:"offset" validate:"min=0"`
	Username string `query:"username" validate:"required"`
}

type GetBidStatusRequest struct {
	BidID    string `params:"bidId" validate:"required"`
	Username string `query:"username" validate:"required"`
}

type UpdateBidStatusRequest struct {
	BidID    string    `params:"bidId" validate:"required"`
	Username string    `query:"username" validate:"required"`
	Status   BidStatus `query:"status" validate:"required,bidstatus"`
}

type EditBidRequest struct {
	BidID      string     `params:"bidId" validate:"required"`
	Username   string     `query:"username" validate:"required"`
	UpdateData UpdateData `json:"updateData" validate:"required"`
}

type RollbackBidRequest struct {
	BidID    string `params:"bidId" validate:"required"`
	Version  string `params:"version" validate:"required,number,min=1"`
	Username string `query:"username" validate:"required"`
}
