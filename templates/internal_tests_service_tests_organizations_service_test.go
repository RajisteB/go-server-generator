package service_tests

import (
	"context"
	"testing"
	"time"

	"{{.Module}}/internal/mocks"
	"{{.Module}}/internal/organizations/models"
	organizationsService "{{.Module}}/internal/organizations/service"
	"{{.Module}}/internal/shared/logger"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrganizationsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource(ctrl)

	service := organizationsService.NewService(mockLogger, mockDatasource)

	t.Run("CreateOrganization_Success", func(t *testing.T) {
		ctx := context.Background()

		clerkRequest := &models.ClerkOrganizationRequest{
			Data: models.OrganizationData{
				ID:        "clerk_org_123",
				Name:      "Test Organization",
				Slug:      "test-org",
				ImageURL:  "https://example.com/image.jpg",
				CreatedAt: 1640995200000,
				UpdatedAt: 1640995200000,
			},
			Object:    "organization",
			Timestamp: 1640995200000,
			Type:      "organization.created",
		}

		expectedOrg := &models.Organization{
			ID:         "org_generated_uuid",
			ClerkOrgID: "clerk_org_123",
			Name:       "Test Organization",
			Slug:       "test-org",
			CreatedAt:  time.UnixMilli(1640995200000),
			UpdatedAt:  time.UnixMilli(1640995200000),
		}

		mockDatasource.EXPECT().CreateOrganization(ctx, gomock.Any()).Return(expectedOrg, nil)

		// Act
		result, err := service.CreateOrganization(ctx, clerkRequest)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedOrg.ID, result.ID)
		assert.Equal(t, expectedOrg.ClerkOrgID, result.ClerkOrgID)
		assert.Equal(t, expectedOrg.Name, result.Name)
	})

	t.Run("GetOrganizationByClerkOrgID_Success", func(t *testing.T) {
		ctx := context.Background()
		clerkOrgID := "clerk_org_123"

		expectedOrg := &models.Organization{
			ID:         "org_456",
			ClerkOrgID: clerkOrgID,
			Name:       "Test Organization",
			Slug:       "test-org",
		}

		mockDatasource.EXPECT().GetOrganizationByClerkOrgID(ctx, clerkOrgID).Return(expectedOrg, nil)

		// Act
		result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedOrg.ID, result.ID)
		assert.Equal(t, expectedOrg.ClerkOrgID, result.ClerkOrgID)
	})

	t.Run("GetOrganizationByClerkOrgID_InvalidID", func(t *testing.T) {
		ctx := context.Background()
		clerkOrgID := "" // Empty ID

		// Act
		result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.ID)
	})
}
