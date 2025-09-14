//go:generate mockgen -destination=../../mocks/mock_organizations_datasource.go -package=mocks {{.Module}}/internal/organizations/datasource OrganizationsDatasource

package datasource

import (
	"context"

	"{{.Module}}/internal/organizations/models"
	"{{.Module}}/internal/shared/assertions"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/uuid"

	"gorm.io/gorm"
)

const (
	pkgName = "organizations"
	layer   = "datasource"
)

type OrganizationsDatasource interface {
	CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, org *models.Organization) (bool, error)
	GetOrganizationByClerkOrgID(ctx context.Context, clerkOrgID string) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error)
	DeleteOrganizationByClerkID(ctx context.Context, clerkID string) (bool, error)
}

type DatasourceImpl struct {
	log *logger.Logger
	db  *gorm.DB
}

func NewDatasource(logger *logger.Logger, db *gorm.DB) OrganizationsDatasource {
	dsLogger := logger.With("package", pkgName, "layer", layer)
	return &DatasourceImpl{log: dsLogger, db: db}
}

func (d *DatasourceImpl) CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	l := d.log.WithContext(ctx).With("operation", "CreateOrganization")

	// Generate UUID if not provided
	if org.ID == "" {
		org.ID = uuid.GenerateNamespaceUUID("org")
	}

	if err := d.db.WithContext(ctx).Create(org).Error; err != nil {
		l.Error("failed to create organization", "error", err)
		return nil, err
	}

	l.Debug("organization created successfully", "org_id", org.ID)
	return org, nil
}

func (d *DatasourceImpl) UpdateOrganization(ctx context.Context, org *models.Organization) (bool, error) {
	l := d.log.WithContext(ctx).With("operation", "UpdateOrganization")

	result := d.db.WithContext(ctx).Model(org).Where("clerk_org_id = ?", org.ClerkOrgID).Updates(org)
	if result.Error != nil {
		l.Error("failed to update organization", "error", result.Error)
		return false, result.Error
	}

	l.Debug("organization updated successfully", "org_id", org.ID, "rows_affected", result.RowsAffected)
	return result.RowsAffected > 0, nil
}

func (d *DatasourceImpl) GetOrganizationByClerkOrgID(ctx context.Context, clerkOrgID string) (*models.Organization, error) {
	l := d.log.WithContext(ctx).With("operation", "GetOrganizationByClerkOrgID")

	if err := assertions.AssertNonEmptyString(clerkOrgID); err != nil {
		l.Debug("invalid clerk org id", "error", err)
		return nil, err
	}

	var org models.Organization
	if err := d.db.WithContext(ctx).Where("clerk_org_id = ?", clerkOrgID).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Debug("organization not found", "clerk_org_id", clerkOrgID)
			return nil, err
		}
		l.Error("failed to get organization by clerk org id", "error", err)
		return nil, err
	}

	l.Debug("organization retrieved successfully", "org_id", org.ID)
	return &org, nil
}

func (d *DatasourceImpl) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	l := d.log.WithContext(ctx).With("operation", "GetOrganizationByID")

	if err := assertions.AssertNonEmptyString(id); err != nil {
		l.Debug("invalid organization id", "error", err)
		return nil, err
	}

	var org models.Organization
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Debug("organization not found", "org_id", id)
			return nil, err
		}
		l.Error("failed to get organization by id", "error", err)
		return nil, err
	}

	l.Debug("organization retrieved successfully", "org_id", org.ID)
	return &org, nil
}

func (d *DatasourceImpl) DeleteOrganizationByClerkID(ctx context.Context, clerkID string) (bool, error) {
	l := d.log.WithContext(ctx).With("operation", "DeleteOrganizationByClerkID")

	if err := assertions.AssertNonEmptyString(clerkID); err != nil {
		l.Debug("invalid clerk org id", "error", err)
		return false, err
	}

	result := d.db.WithContext(ctx).Where("clerk_org_id = ?", clerkID).Delete(&models.Organization{})
	if result.Error != nil {
		l.Error("failed to delete organization", "error", result.Error)
		return false, result.Error
	}

	l.Debug("organization deleted successfully", "clerk_org_id", clerkID, "rows_affected", result.RowsAffected)
	return result.RowsAffected > 0, nil
}
