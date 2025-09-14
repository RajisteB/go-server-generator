package models

import (
	"time"

	"{{.Module}}/internal/shared/validation"
)

type User struct {
	ID               string    `json:"id" gorm:"unique" validate:"required"`            // Custom ID for your system
	ClerkUserID      string    `json:"clerk_user_id" gorm:"unique" validate:"required"` // Clerk's User ID
	Email            string    `json:"email" gorm:"unique"`                             // User's email
	FirstName        string    `json:"first_name" validate:"required"`
	LastName         string    `json:"last_name" validate:"required"`
	OrganizationID   string    `json:"organization_id"` // Organization ID
	IsBusinessAcount bool      `json:"is_business_account"`
	BusinessName     *string   `json:"business_name"`
	ProfileImageURL  *string   `json:"profile_image_url"`
	MFAEnabled       bool      `json:"mfa_enabled"`
	TwoFactorEnabled bool      `json:"two_factor_enabled"`
	IsBanned         bool      `json:"is_banned"`
	LastActiveAt     time.Time `json:"last_active_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type ClerkUserRequest struct {
	Data            UserData        `json:"data" validate:"required"`
	EventAttributes EventAttributes `json:"event_attributes" validate:"required"`
	Object          string          `json:"object" validate:"required"`
	Timestamp       int64           `json:"timestamp" validate:"required"`
	Type            string          `json:"type" validate:"required"`
}

type ClerkUserDeleteRequest struct {
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

func (cur *ClerkUserRequest) ToUser() User {
	email := ""
	if len(cur.Data.EmailAddresses) > 0 {
		email = cur.Data.EmailAddresses[0].EmailAddress
	}

	return User{
		ID:               cur.Data.ExternalID, // Custom ID for your system
		ClerkUserID:      cur.Data.ID,         // Clerk's User ID
		Email:            email,               // Use the first email address
		OrganizationID:   cur.Data.OrganizationID,
		FirstName:        cur.Data.FirstName,
		LastName:         cur.Data.LastName,
		ProfileImageURL:  &cur.Data.ProfileImageURL,
		MFAEnabled:       cur.Data.PasswordEnabled,
		TwoFactorEnabled: cur.Data.TwoFactorEnabled,
		IsBanned:         cur.Data.IsBanned,
		CreatedAt:        time.UnixMilli(cur.Data.CreatedAt), // Convert Unix milliseconds to time.Time
		UpdatedAt:        time.UnixMilli(cur.Data.UpdatedAt), // Convert Unix milliseconds to time.Time
		LastActiveAt:     time.UnixMilli(cur.Data.LastSignInAt),
	}
}

type UserData struct {
	Birthday              string                 `json:"birthday"`
	CreatedAt             int64                  `json:"created_at"`
	EmailAddresses        []EmailAddress         `json:"email_addresses"`
	ExternalAccounts      []interface{}          `json:"external_accounts"`
	ExternalID            string                 `json:"external_id"`
	FirstName             string                 `json:"first_name"`
	Gender                string                 `json:"gender"`
	ID                    string                 `json:"id"`
	OrganizationID        string                 `json:"organization_id"`
	IsBanned              bool                   `json:"banned"`
	ImageURL              string                 `json:"image_url"`
	LastName              string                 `json:"last_name"`
	LastSignInAt          int64                  `json:"last_sign_in_at"`
	Object                string                 `json:"object"`
	PasswordEnabled       bool                   `json:"password_enabled"`
	PhoneNumbers          []interface{}          `json:"phone_numbers"`
	PrimaryEmailAddressID string                 `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  *string                `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID   *string                `json:"primary_web3_wallet_id"`
	PrivateMetadata       map[string]interface{} `json:"private_metadata"`
	ProfileImageURL       string                 `json:"profile_image_url"`
	PublicMetadata        map[string]interface{} `json:"public_metadata"`
	TwoFactorEnabled      bool                   `json:"two_factor_enabled"`
	UnsafeMetadata        map[string]interface{} `json:"unsafe_metadata"`
	UpdatedAt             int64                  `json:"updated_at"`
	Username              *string                `json:"username"`
	Web3Wallets           []interface{}          `json:"web3_wallets"`
}

type EmailAddress struct {
	EmailAddress string        `json:"email_address"`
	ID           string        `json:"id"`
	LinkedTo     []interface{} `json:"linked_to"`
	Object       string        `json:"object"`
	Verification Verification  `json:"verification"`
}

type Verification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

// Sanitize cleans all string fields in the User struct
func (u *User) Sanitize() {
	u.ID = validation.SanitizeString(u.ID)
	u.ClerkUserID = validation.SanitizeString(u.ClerkUserID)
	u.Email = validation.SanitizeEmail(u.Email)
	u.FirstName = validation.SanitizeString(u.FirstName)
	u.LastName = validation.SanitizeString(u.LastName)
	u.OrganizationID = validation.SanitizeString(u.OrganizationID)

	if u.BusinessName != nil {
		sanitized := validation.SanitizeString(*u.BusinessName)
		u.BusinessName = &sanitized
	}
	if u.ProfileImageURL != nil {
		sanitized := validation.SanitizeString(*u.ProfileImageURL)
		u.ProfileImageURL = &sanitized
	}
}

// Sanitize cleans all string fields in the ClerkUserRequest struct
func (cur *ClerkUserRequest) Sanitize() {
	cur.Object = validation.SanitizeString(cur.Object)
	cur.Type = validation.SanitizeString(cur.Type)
	// Note: Data and EventAttributes are nested structs - sanitize them separately if needed
}

// Sanitize cleans all string fields in the ClerkUserDeleteRequest struct
func (cudr *ClerkUserDeleteRequest) Sanitize() {
	cudr.Object = validation.SanitizeString(cudr.Object)
	cudr.Type = validation.SanitizeString(cudr.Type)
}

// Sanitize cleans all string fields in the DeleteData struct
func (dd *DeleteData) Sanitize() {
	dd.ID = validation.SanitizeString(dd.ID)
}
