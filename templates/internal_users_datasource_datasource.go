package datasource

import (
	"context"

	"{{.Module}}/internal/shared/assertions"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/uuid"
	"{{.Module}}/internal/users/models"

	"gorm.io/gorm"
)

const (
	pkgName = "users"
	layer   = "datasource"
)

type UsersDatasource interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (bool, error)
	GetUserByClerkUserID(ctx context.Context, clerkUserID string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	DeleteUserByClerkID(ctx context.Context, clerkID string) (bool, error)
	UpdateUserOrganization(ctx context.Context, clerkUserID string, orgID string) (bool, error)
}

type DatasourceImpl struct {
	log *logger.Logger
	db  *gorm.DB
}

func NewDatasource(logger *logger.Logger, db *gorm.DB) UsersDatasource {
	dsLogger := logger.With("package", pkgName, "layer", layer)
	return &DatasourceImpl{log: dsLogger, db: db}
}

func (d *DatasourceImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	l := d.log.WithContext(ctx).With("operation", "CreateUser")

	// Generate UUID if not provided
	if user.ID == "" {
		user.ID = uuid.GenerateNamespaceUUID("usr")
	}

	if err := d.db.WithContext(ctx).Create(user).Error; err != nil {
		l.Error("failed to create user", "error", err)
		return nil, err
	}

	l.Debug("user created successfully", "user_id", user.ID)
	return user, nil
}

func (d *DatasourceImpl) UpdateUser(ctx context.Context, user *models.User) (bool, error) {
	l := d.log.WithContext(ctx).With("operation", "UpdateUser")

	result := d.db.WithContext(ctx).Model(user).Where("clerk_user_id = ?", user.ClerkUserID).Updates(user)
	if result.Error != nil {
		l.Error("failed to update user", "error", result.Error)
		return false, result.Error
	}

	l.Debug("user updated successfully", "user_id", user.ID, "rows_affected", result.RowsAffected)
	return result.RowsAffected > 0, nil
}

func (d *DatasourceImpl) GetUserByClerkUserID(ctx context.Context, clerkUserID string) (*models.User, error) {
	l := d.log.WithContext(ctx).With("operation", "GetUserByClerkUserID")

	if err := assertions.AssertNonEmptyString(clerkUserID); err != nil {
		l.Debug("invalid clerk user id", "error", err)
		return nil, err
	}

	var user models.User
	if err := d.db.WithContext(ctx).Where("clerk_user_id = ?", clerkUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Debug("user not found", "clerk_user_id", clerkUserID)
			return nil, err
		}
		l.Error("failed to get user by clerk user id", "error", err)
		return nil, err
	}

	l.Debug("user retrieved successfully", "user_id", user.ID)
	return &user, nil
}

func (d *DatasourceImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	l := d.log.WithContext(ctx).With("operation", "GetUserByID")

	if err := assertions.AssertNonEmptyString(id); err != nil {
		l.Debug("invalid user id", "error", err)
		return nil, err
	}

	var user models.User
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Debug("user not found", "user_id", id)
			return nil, err
		}
		l.Error("failed to get user by id", "error", err)
		return nil, err
	}

	l.Debug("user retrieved successfully", "user_id", user.ID)
	return &user, nil
}

func (d *DatasourceImpl) DeleteUserByClerkID(ctx context.Context, clerkID string) (bool, error) {
	l := d.log.WithContext(ctx).With("operation", "DeleteUserByClerkID")

	if err := assertions.AssertNonEmptyString(clerkID); err != nil {
		l.Debug("invalid clerk user id", "error", err)
		return false, err
	}

	result := d.db.WithContext(ctx).Where("clerk_user_id = ?", clerkID).Delete(&models.User{})
	if result.Error != nil {
		l.Error("failed to delete user", "error", result.Error)
		return false, result.Error
	}

	l.Debug("user deleted successfully", "clerk_user_id", clerkID, "rows_affected", result.RowsAffected)
	return result.RowsAffected > 0, nil
}

func (d *DatasourceImpl) UpdateUserOrganization(ctx context.Context, clerkUserID string, orgID string) (bool, error) {
	l := d.log.WithContext(ctx).With("operation", "UpdateUserOrganization")

	if err := assertions.AssertNonEmptyString(clerkUserID); err != nil {
		l.Debug("invalid clerk user id", "error", err)
		return false, err
	}

	if err := assertions.AssertNonEmptyString(orgID); err != nil {
		l.Debug("invalid organization id", "error", err)
		return false, err
	}

	result := d.db.WithContext(ctx).Model(&models.User{}).Where("clerk_user_id = ?", clerkUserID).Update("organization_id", orgID)
	if result.Error != nil {
		l.Error("failed to update user organization", "error", result.Error)
		return false, result.Error
	}

	l.Debug("user organization updated successfully", "clerk_user_id", clerkUserID, "organization_id", orgID, "rows_affected", result.RowsAffected)
	return result.RowsAffected > 0, nil
}
