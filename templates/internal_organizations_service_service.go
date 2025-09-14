package organizations

import (
	"context"

	organizationsDatasource "{{.Module}}/internal/organizations/datasource"
	"{{.Module}}/internal/organizations/models"
	"{{.Module}}/internal/shared/assertions"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/uuid"
	"{{.Module}}/internal/shared/validation"
)

const (
	pkgName = "organizations"
	layer   = "service"
)

type OrganizationsService interface {
	CreateOrganization(ctx context.Context, org *models.ClerkOrganizationRequest) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, org *models.ClerkOrganizationRequest) (bool, error)
	GetOrganizationByClerkOrgID(ctx context.Context, clerkOrgID string) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, clerkID string) (bool, error)
}

type ServiceImpl struct {
	log  *logger.Logger
	data organizationsDatasource.OrganizationsDatasource
}

func NewService(logger *logger.Logger, datasource organizationsDatasource.OrganizationsDatasource) OrganizationsService {
	serviceLogger := logger.With("package", pkgName, "layer", layer)
	return &ServiceImpl{log: serviceLogger, data: datasource}
}

func (s *ServiceImpl) CreateOrganization(ctx context.Context, request *models.ClerkOrganizationRequest) (*models.Organization, error) {
	l := s.log.WithContext(ctx).With("operation", "CreateOrganization")
	orgID := uuid.GenerateNamespaceUUID("org")
	org := request.ToOrganization()
	org.ID = orgID

	if err := validation.ValidateStruct(org); err != nil {
		l.Error("failed to parse clerk organization request to organization", "error", err)
		return &models.Organization{}, err
	}

	created, err := s.data.CreateOrganization(ctx, &org)
	if err != nil {
		return &models.Organization{}, err
	}

	return created, nil
}

func (s *ServiceImpl) UpdateOrganization(ctx context.Context, request *models.ClerkOrganizationRequest) (bool, error) {
	l := s.log.WithContext(ctx).With("operation", "UpdateOrganization")
	org := request.ToOrganization()

	if err := assertions.AssertNonEmptyString(org.ClerkOrgID); err != nil {
		l.Error("failed to validate clerk org id", "error", err)
		return false, err
	}

	updated, err := s.data.UpdateOrganization(ctx, &org)
	if err != nil {
		return false, err
	}

	return updated, nil
}

func (s *ServiceImpl) GetOrganizationByClerkOrgID(ctx context.Context, clerkOrgID string) (*models.Organization, error) {
	l := s.log.WithContext(ctx).With("operation", "GetOrganizationByClerkOrgID")

	if err := assertions.AssertNonEmptyString(clerkOrgID); err != nil {
		l.Debug("failed to validate clerk org id", "error", err)
		return &models.Organization{}, err
	}

	org, err := s.data.GetOrganizationByClerkOrgID(ctx, clerkOrgID)
	if err != nil {
		return &models.Organization{}, err
	}

	return org, nil
}

func (s *ServiceImpl) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	l := s.log.WithContext(ctx).With("operation", "GetOrganizationByID")

	if err := assertions.AssertNonEmptyString(id); err != nil {
		l.Debug("failed to validate organization id", "error", err)
		return &models.Organization{}, err
	}

	org, err := s.data.GetOrganizationByID(ctx, id)
	if err != nil {
		return &models.Organization{}, err
	}

	return org, nil
}

func (s *ServiceImpl) DeleteOrganization(ctx context.Context, clerkID string) (bool, error) {
	l := s.log.WithContext(ctx).With("operation", "DeleteOrganization")

	if err := assertions.AssertNonEmptyString(clerkID); err != nil {
		l.Debug("failed to validate clerk org id", "error", err)
		return false, err
	}

	deleted, err := s.data.DeleteOrganizationByClerkID(ctx, clerkID)
	if err != nil {
		return false, err
	}

	return deleted, nil
}
