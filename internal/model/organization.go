package model

import "time"

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

type Organization struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        OrganizationType `json:"type" `
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type OrganizationResponsible struct {
	ID string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	UserID string`json:"user_id"`
}

