package users

import (
	"{{.Module}}/internal/shared/constants"
	httpHelpers "{{.Module}}/internal/shared/http"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/middleware"
	"{{.Module}}/internal/shared/validation"
	"{{.Module}}/internal/users/models"
	users "{{.Module}}/internal/users/service"

	"net/http"
)

const (
	pkgName = "users"
	layer   = "controller"
)

type UsersController interface {
	CreateUserFromClerk(w http.ResponseWriter, r *http.Request) error
	UpdateUserFromClerk(w http.ResponseWriter, r *http.Request) error
	DeleteUser(w http.ResponseWriter, r *http.Request) error
}

type ControllerImpl struct {
	log     *logger.Logger
	service users.UsersService
}

func NewController(logger *logger.Logger, service users.UsersService) UsersController {
	ctrlLogger := logger.With("package", pkgName, "layer", layer)
	return &ControllerImpl{log: ctrlLogger, service: service}
}

func (c *ControllerImpl) CreateUserFromClerk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "CreateUserFromClerk")

	userRequest := &models.ClerkUserRequest{}

	if err := middleware.SafeJSONDecoder(r, userRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode create user request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	if err := validation.ValidateStruct(userRequest); err != nil {
		l.Debug("failed to validate create user from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	user, err := c.service.CreateUser(ctx, userRequest)
	if err != nil {
		l.Error("failed to create user from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusCreated, user)
}

func (c *ControllerImpl) UpdateUserFromClerk(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "UpdateUserFromClerk")

	userRequest := &models.ClerkUserRequest{}

	if err := middleware.SafeJSONDecoder(r, userRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode update user request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	if err := validation.ValidateStruct(userRequest); err != nil {
		l.Debug("failed to validate update user from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	user, err := c.service.UpdateUser(ctx, userRequest)

	if err != nil {
		l.Error("failed to update user from clerk request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, user)
}

func (c *ControllerImpl) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "DeleteUser")

	userDeleteRequest := &models.ClerkUserDeleteRequest{}

	if err := middleware.SafeJSONDecoder(r, userDeleteRequest, constants.JSONMaxSize); err != nil {
		l.Debug("failed to decode delete user request", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	deleted, err := c.service.DeleteUser(ctx, userDeleteRequest.Data.ID)
	if err != nil {
		l.Error("failed to delete user", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, deleted)
}
