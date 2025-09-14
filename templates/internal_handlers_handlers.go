package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"{{.Module}}/internal/conf"
	"{{.Module}}/internal/shared/assertions"
	"{{.Module}}/internal/shared/constants"
	httpHelpers "{{.Module}}/internal/shared/http"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/middleware"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type Handler struct {
	Logger       *logger.Logger
	Dependencies *conf.Dependencies
}

func NewHandler(logger *logger.Logger, dep *conf.Dependencies) *Handler {
	return &Handler{
		Logger:       logger,
		Dependencies: dep,
	}
}

func (handler *Handler) GetCSRFToken(w http.ResponseWriter, r *http.Request) error {
	token := csrf.Token(r)
	return httpHelpers.RespondWithJSON(w, http.StatusOK, map[string]string{
		"csrf_token": token,
	})
}

func (h *Handler) HandleClerkWebhook(w http.ResponseWriter, r *http.Request) error {
	l := h.Logger.WithContext(r.Context()).With("operation", "handleClerkWebhook")

	body, err := httpHelpers.SafeBodyReader(r, 52428800)
	if err != nil {
		l.Error("failed to read webhook body", "error", err)
		return err
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		l.Error("failed to unmarshal webhook body", "error", err)
		return err
	}

	eventType := data["type"].(string)
	if err := assertions.Assert(eventType != "", "eventType is required"); err != nil {
		l.Error("eventType is required", "error", err)
		return err
	}

	l.Debug("webhook event received", "data", data)
	r.Body = io.NopCloser(bytes.NewReader(body))

	switch eventType {
	case string(constants.WebhookEventUserCreated):
		return h.Dependencies.Controllers.Users.CreateUserFromClerk(w, r)
	case string(constants.WebhookEventUserUpdated):
		return h.Dependencies.Controllers.Users.UpdateUserFromClerk(w, r)
	case string(constants.WebhookEventUserDeleted):
		return h.Dependencies.Controllers.Users.DeleteUser(w, r)
	case string(constants.WebhookEventOrganizationCreated):
		return h.Dependencies.Controllers.Organizations.CreateOrganizationFromClerk(w, r)
	case string(constants.WebhookEventOrganizationUpdated):
		return h.Dependencies.Controllers.Organizations.UpdateOrganizationFromClerk(w, r)
	case string(constants.WebhookEventOrganizationDeleted):
		return h.Dependencies.Controllers.Organizations.DeleteOrganization(w, r)
	default:
		l.Error("unsupported webhook event", "eventType", eventType)
		return errors.New("unsupported webhook event")
	}
}

func (handler *Handler) Register() *mux.Router {
	router := mux.NewRouter()
	prefix := fmt.Sprintf("/%s", constants.SERVICE_API_PREFIX)
	mw := middleware.NewMiddleware(handler.Dependencies.ExternalDependencies.Clerk, handler.Dependencies.Config.Clerk.Secret)

	// Generate CSRF auth key if not provided
	csrfAuthKey := []byte(handler.Dependencies.Config.CSRF.AuthKey)
	if len(csrfAuthKey) == 0 {
		csrfAuthKey = make([]byte, 32)
		if _, err := rand.Read(csrfAuthKey); err != nil {
			handler.Logger.Error("failed to generate CSRF auth key", "error", err)
		}
		handler.Logger.Warn("using generated CSRF auth key - set CSRF_AUTH_KEY environment variable for production")
	}

	csrfMiddleware := mw.CSRFMiddleware(csrfAuthKey, handler.Dependencies.Config.CSRF.Secure)
	securityConfig := middleware.SecurityConfig{
		CSPPolicy:          handler.Dependencies.Config.Security.CSPPolicy,
		HSTSMaxAge:         handler.Dependencies.Config.Security.HSTSMaxAge,
		FrameOptions:       handler.Dependencies.Config.Security.FrameOptions,
		ContentTypeOptions: handler.Dependencies.Config.Security.ContentTypeOptions,
		ReferrerPolicy:     handler.Dependencies.Config.Security.ReferrerPolicy,
		PermissionsPolicy:  handler.Dependencies.Config.Security.PermissionsPolicy,
	}
	securityMiddleware := mw.SecurityHeadersMiddleware(securityConfig)

	requestLimitsConfig := middleware.RequestLimitsConfig{
		MaxRequestSize:    handler.Dependencies.Config.RequestLimits.MaxRequestSize,
		MaxHeaderSize:     handler.Dependencies.Config.RequestLimits.MaxHeaderSize,
		MaxFileUploadSize: handler.Dependencies.Config.RequestLimits.MaxFileUploadSize,
		ReadTimeout:       handler.Dependencies.Config.RequestLimits.ReadTimeout,
		WriteTimeout:      handler.Dependencies.Config.RequestLimits.WriteTimeout,
		DebugHeaders:      handler.Dependencies.Config.Server.Environment == "development",
	}
	requestSizeLimitMiddleware := mw.RequestSizeLimitMiddleware(requestLimitsConfig)
	requestTimeoutMiddleware := mw.RequestTimeoutMiddleware(requestLimitsConfig)

	// Apply security headers and request limits to all routes
	router.Use(securityMiddleware)
	router.Use(requestSizeLimitMiddleware)
	router.Use(requestTimeoutMiddleware)

	api := router.PathPrefix(prefix).Subrouter()
	api.Use(mw.LoggerMiddleware)
	api.Use(mw.RateLimiterMiddleware)

	private := router.PathPrefix(prefix).Subrouter()
	private.Use(mw.LoggerMiddleware)
	private.Use(mw.RateLimiterMiddleware)
	private.Use(mw.ClerkAuthMiddleware)
	private.Use(csrfMiddleware)

	webhook := router.PathPrefix(prefix).Subrouter()
	webhook.Use(mw.LoggerMiddleware)
	webhook.Use(mw.RateLimiterMiddleware)
	webhook.Use(mw.ClerkWebhookMiddleware)

	// Health
	api.Handle("/health", httpHelpers.HandlerFunc(handler.Dependencies.Controllers.Health.GetHealth)).Methods(http.MethodGet)

	// CSRF Token
	api.Handle("/csrf-token", httpHelpers.HandlerFunc(handler.GetCSRFToken)).Methods(http.MethodGet)

	// Identity Webhook
	webhook.Handle("/identity/clerk", httpHelpers.HandlerFunc(handler.HandleClerkWebhook)).Methods(http.MethodPost)

	// Organizations
	private.Handle("/organizations/{id}", httpHelpers.HandlerFunc(handler.Dependencies.Controllers.Organizations.GetOrganizationByID)).Methods(http.MethodGet)
	private.Handle("/organizations/clerk/{clerk_id}", httpHelpers.HandlerFunc(handler.Dependencies.Controllers.Organizations.GetOrganizationByClerkID)).Methods(http.MethodGet)

	return router
}
