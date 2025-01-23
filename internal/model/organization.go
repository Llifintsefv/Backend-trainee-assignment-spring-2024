package model

import "time"

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

type Organization struct {
	id          string `json:"id"`
	name        string `json:"name"`
	description string `json:"description"`
	Type        OrganizationType `json:"type" `
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"updated_at"`
}

type OrganizationResponsible struct {
	ID string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	UserID string`json:"user_id"`
}

