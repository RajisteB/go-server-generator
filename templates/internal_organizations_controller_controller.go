package organizations

import (
	"errors"
	"net/http"

	"{{.Module}}/internal/organizations/models"
	organizationsService "{{.Module}}/internal/organizations/service"
	"{{.Module}}/internal/shared/constants"
	httpHelpers "{{.Module}}/internal/shared/http"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/middleware"
	"{{.Module}}/internal/shared/validation"

	"github.com/gorilla/mux"
)

const (
	pkgName = "organizations"
	layer   = "controller"
)

type OrganizationsController interface {
	GetOrganizationByID(w http.ResponseWriter, r *http.Request) error
	GetOrganizationByClerkID(w http.ResponseWriter, r *http.Request) error
	CreateOrganizationFromClerk(w http.ResponseWriter, r *http.Request) error
	UpdateOrganizationFromClerk(w http.ResponseWriter, r *http.Request) error
	DeleteOrganization(w http.ResponseWriter, r *http.Request) error
}

type ControllerImpl struct {
	log     *logger.Logger
	service organizationsService.OrganizationsService
}

func NewController(logger *logger.Logger, service organizationsService.OrganizationsService) OrganizationsController {
	ctrlLogger := logger.With("package", pkgName, "layer", layer)
	return &ControllerImpl{log: ctrlLogger, service: service}
}

func (c *ControllerImpl) GetOrganizationByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "GetOrganizationByID")

	// Extract ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		l.Debug("missing organization id in path")
		return httpHelpers.RespondWithError(w, errors.New("organization id is required"))
	}

	org, err := c.service.GetOrganizationByID(ctx, id)
	if err != nil {
		l.Error("failed to get organization by id", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, org)
}

func (c *ControllerImpl) GetOrganizationByClerkID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "GetOrganizationByClerkID")

	// Extract Clerk ID from URL path
	vars := mux.Vars(r)
	clerkID := vars["clerk_id"]

	if clerkID == "" {
		l.Debug("missing clerk organization id in path")
		return httpHelpers.RespondWithError(w, errors.New("clerk organization id is required"))
	}

	org, err := c.service.GetOrganizationByClerkOrgID(ctx, clerkID)
	if err != nil {
		l.Error("failed to get organization by clerk id", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, org)
}

func (c *ControllerImpl) CreateOrganizationFromClerk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "CreateOrganizationFromClerk")

	orgRequest := &models.ClerkOrganizationRequest{}

	if err := middleware.SafeJSONDecoder(r, orgRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode create organization request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	if err := validation.ValidateStruct(orgRequest); err != nil {
		l.Debug("failed to validate create organization from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	org, err := c.service.CreateOrganization(ctx, orgRequest)
	if err != nil {
		l.Error("failed to create organization from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusCreated, org)
}

func (c *ControllerImpl) UpdateOrganizationFromClerk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "UpdateOrganizationFromClerk")

	orgRequest := &models.ClerkOrganizationRequest{}

	if err := middleware.SafeJSONDecoder(r, orgRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode update organization request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	if err := validation.ValidateStruct(orgRequest); err != nil {
		l.Debug("failed to validate update organization from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	org, err := c.service.UpdateOrganization(ctx, orgRequest)
	if err != nil {
		l.Error("failed to update organization from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, org)
}

func (c *ControllerImpl) DeleteOrganization(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "DeleteOrganization")

	orgDeleteRequest := &models.ClerkOrganizationDeleteRequest{}

	if err := middleware.SafeJSONDecoder(r, orgDeleteRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode delete organization request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	deleted, err := c.service.DeleteOrganization(ctx, orgDeleteRequest.Data.ID)
	if err != nil {
		l.Error("failed to delete organization", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, deleted)
}
