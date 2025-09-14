package models

import (
	"time"

	"{{.Module}}/internal/shared/validation"
)

type Organization struct {
	ID         string    `json:"id" gorm:"unique" validate:"required"`
	ClerkOrgID string    `json:"clerk_org_id" gorm:"unique" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Slug       string    `json:"slug" gorm:"unique"`
	ImageURL   *string   `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type ClerkOrganizationRequest struct {
	Data            OrganizationData `json:"data" validate:"required"`
	EventAttributes EventAttributes  `json:"event_attributes" validate:"required"`
	Object          string           `json:"object" validate:"required"`
	Timestamp       int64            `json:"timestamp" validate:"required"`
	Type            string           `json:"type" validate:"required"`
}

type ClerkOrganizationDeleteRequest struct {
	Data            DeleteData      `json:"data" validate:"required"`
	EventAttributes EventAttributes `json:"event_attributes" validate:"required"`
	Object          string          `json:"object" validate:"required"`
	Timestamp       int64           `json:"timestamp" validate:"required"`
	Type            string          `json:"type" validate:"required"`
}

type DeleteData struct {
	ID      string `json:"id" validate:"required"`
	Deleted bool   `json:"deleted"`
}

type OrganizationData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ImageURL  string `json:"image_url"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

func (cor *ClerkOrganizationRequest) ToOrganization() Organization {
	return Organization{
		ClerkOrgID: cor.Data.ID,
		Name:       cor.Data.Name,
		Slug:       cor.Data.Slug,
		ImageURL:   &cor.Data.ImageURL,
		CreatedAt:  time.UnixMilli(cor.Data.CreatedAt),
		UpdatedAt:  time.UnixMilli(cor.Data.UpdatedAt),
	}
}

// Sanitize cleans all string fields in the Organization struct
func (o *Organization) Sanitize() {
	o.ID = validation.SanitizeString(o.ID)
	o.ClerkOrgID = validation.SanitizeString(o.ClerkOrgID)
	o.Name = validation.SanitizeString(o.Name)
	o.Slug = validation.SanitizeString(o.Slug)

	if o.ImageURL != nil {
		sanitized := validation.SanitizeString(*o.ImageURL)
		o.ImageURL = &sanitized
	}
}

// Sanitize cleans all string fields in the ClerkOrganizationRequest struct
func (cor *ClerkOrganizationRequest) Sanitize() {
	cor.Object = validation.SanitizeString(cor.Object)
	cor.Type = validation.SanitizeString(cor.Type)
}

// Sanitize cleans all string fields in the ClerkOrganizationDeleteRequest struct
func (codr *ClerkOrganizationDeleteRequest) Sanitize() {
	codr.Object = validation.SanitizeString(codr.Object)
	codr.Type = validation.SanitizeString(codr.Type)
}

// Sanitize cleans all string fields in the DeleteData struct
func (dd *DeleteData) Sanitize() {
	dd.ID = validation.SanitizeString(dd.ID)
}
