package constants

import "time"

// HTTP Methods
type AllowedMethod string

const (
	AllowedMethodGET     AllowedMethod = "GET"
	AllowedMethodPOST    AllowedMethod = "POST"
	AllowedMethodPUT     AllowedMethod = "PUT"
	AllowedMethodPATCH   AllowedMethod = "PATCH"
	AllowedMethodDELETE  AllowedMethod = "DELETE"
	AllowedMethodOPTIONS AllowedMethod = "OPTIONS"
)

// Server Configuration
const (
	SERVICE_API_PREFIX = "api/v1"

	// CORS Origins
	ServerAllowedOriginLocal       = "http://localhost:3000"
	ServerAllowedOriginVite        = "http://localhost:5173"
	ServerAllowedOriginReact       = "http://localhost:3001"
	ServerAllowedOriginReactNative = "http://localhost:8081"
	ServerAllowedOriginPostman     = "https://www.postman.com"
)

// Timeouts
const (
	WriteTimeout        = 15 * time.Second
	ReadTimeout         = 15 * time.Second
	IdleTimeout         = 60 * time.Second
	ShutdownGracePeriod = 30 * time.Second
)

// Request Limits
const (
	JSONMaxSize = 10 * 1024 * 1024 // 10MB
)

// Clerk Webhook Event Types
type WebhookEventType string

const (
	WebhookEventUserCreated         WebhookEventType = "user.created"
	WebhookEventUserUpdated         WebhookEventType = "user.updated"
	WebhookEventUserDeleted         WebhookEventType = "user.deleted"
	WebhookEventOrganizationCreated WebhookEventType = "organization.created"
	WebhookEventOrganizationUpdated WebhookEventType = "organization.updated"
	WebhookEventOrganizationDeleted WebhookEventType = "organization.deleted"
)
