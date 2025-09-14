//go:generate mockgen -destination=../../mocks/mock_users_service.go -package=mocks {{.Module}}/internal/users/service UsersService

package users

import (
	"context"

	"{{.Module}}/internal/shared/assertions"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/uuid"
	"{{.Module}}/internal/shared/validation"
	usersDatasource "{{.Module}}/internal/users/datasource"
	"{{.Module}}/internal/users/models"
)

const (
	pkgName = "users"
	layer   = "service"
)

type UsersService interface {
	CreateUser(ctx context.Context, user *models.ClerkUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.ClerkUserRequest) (bool, error)
	GetUserByClerkUserID(ctx context.Context, clerkUserID string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	UpdateUserOrganization(ctx context.Context, clerkUserID string, orgID string) (bool, error)
}

type ServiceImpl struct {
	log  *logger.Logger
	data usersDatasource.UsersDatasource
}

func NewService(logger *logger.Logger, datasource usersDatasource.UsersDatasource) UsersService {
	serviceLogger := logger.With("package", pkgName, "layer", layer)
	return &ServiceImpl{log: serviceLogger, data: datasource}
}

func (s *ServiceImpl) CreateUser(ctx context.Context, request *models.ClerkUserRequest) (*models.User, error) {
	l := s.log.WithContext(ctx).With("operation", "CreateUser")
	userID := uuid.GenerateNamespaceUUID("usr")
	user := request.ToUser()
	user.ID = userID

	if err := validation.ValidateStruct(user); err != nil {
		l.Error("failed to parse clerk user request to user", "error", err)
		return &models.User{}, err
	}

	created, err := s.data.CreateUser(ctx, &user)
	if err != nil {
		return &models.User{}, err
	}

	return created, nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, request *models.ClerkUserRequest) (bool, error) {
	l := s.log.WithContext(ctx).With("operation", "UpdateUser")
	user := request.ToUser()

	if err := assertions.AssertNonEmptyString(user.ClerkUserID); err != nil {
		l.Error("failed to validate clerk user id", "error", err)
		return false, err
	}

	updated, err := s.data.UpdateUser(ctx, &user)
	if err != nil {
		return false, err
	}

	return updated, nil
}

func (s *ServiceImpl) GetUserByClerkUserID(ctx context.Context, clerkUserID string) (*models.User, error) {
	l := s.log.WithContext(ctx).With("operation", "GetUserByClerkUserID")

	if err := assertions.AssertNonEmptyString(clerkUserID); err != nil {
		l.Debug("failed to validate clerk user id", "error", err)
		return &models.User{}, err
	}

	user, err := s.data.GetUserByClerkUserID(ctx, clerkUserID)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (s *ServiceImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	l := s.log.WithContext(ctx).With("operation", "GetUserByID")

	if err := assertions.AssertNonEmptyString(id); err != nil {
		l.Debug("failed to validate user id", "error", err)
		return &models.User{}, err
	}

	user, err := s.data.GetUserByID(ctx, id)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (s *ServiceImpl) DeleteUser(ctx context.Context, clerkID string) (bool, error) {
	l := s.log.WithContext(ctx).With("operation", "DeleteUser")

	if err := assertions.AssertNonEmptyString(clerkID); err != nil {
		l.Debug("failed to validate clerk user id", "error", err)
		return false, err
	}

	deleted, err := s.data.DeleteUserByClerkID(ctx, clerkID)
	if err != nil {
		return false, err
	}

	return deleted, nil
}

func (s *ServiceImpl) UpdateUserOrganization(ctx context.Context, clerkUserID string, orgID string) (bool, error) {
	l := s.log.WithContext(ctx).With("operation", "UpdateUserOrganization")

	if err := assertions.AssertNonEmptyString(clerkUserID); err != nil {
		l.Debug("failed to validate clerk user id", "error", err)
		return false, err
	}

	if err := assertions.AssertNonEmptyString(orgID); err != nil {
		l.Debug("failed to validate organization id", "error", err)
		return false, err
	}

	updated, err := s.data.UpdateUserOrganization(ctx, clerkUserID, orgID)
	if err != nil {
		return false, err
	}

	return updated, nil
}
