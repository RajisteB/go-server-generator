package datasource_tests

import (
	"context"
	"testing"
	"time"

	"{{.Module}}/internal/mocks"
	organizationsDatasource "{{.Module}}/internal/organizations/datasource"
	"{{.Module}}/internal/organizations/models"
	"{{.Module}}/internal/shared/logger"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestOrganizationsDatasource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewLogger()
	mockDB := mocks.NewMockGormDB(ctrl)

	datasource := organizationsDatasource.NewDatasource(mockLogger, mockDB)

	t.Run("CreateOrganization_Success", func(t *testing.T) {
		ctx := context.Background()

		org := &models.Organization{
			ID:         "org_123",
			ClerkOrgID: "clerk_org_456",
			Name:       "Test Organization",
			Slug:       "test-org",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Create(org).Return(mockDB)
		mockDB.EXPECT().Error().Return(nil)

		// Act
		result, err := datasource.CreateOrganization(ctx, org)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, org.ID, result.ID)
		assert.Equal(t, org.ClerkOrgID, result.ClerkOrgID)
		assert.Equal(t, org.Name, result.Name)
	})

	t.Run("GetOrganizationByClerkOrgID_Success", func(t *testing.T) {
		ctx := context.Background()
		clerkOrgID := "clerk_org_456"

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Where("clerk_org_id = ?", clerkOrgID).Return(mockDB)
		mockDB.EXPECT().First(gomock.Any()).Return(mockDB)
		mockDB.EXPECT().Error().Return(nil)

		// Act
		result, err := datasource.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetOrganizationByClerkOrgID_NotFound", func(t *testing.T) {
		ctx := context.Background()
		clerkOrgID := "clerk_org_456"

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Where("clerk_org_id = ?", clerkOrgID).Return(mockDB)
		mockDB.EXPECT().First(gomock.Any()).Return(mockDB)
		mockDB.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		// Act
		result, err := datasource.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("GetOrganizationByClerkOrgID_InvalidID", func(t *testing.T) {
		ctx := context.Background()
		clerkOrgID := "" // Empty ID

		// Act
		result, err := datasource.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
