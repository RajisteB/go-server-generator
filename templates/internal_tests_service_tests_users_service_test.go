package service_tests

import (
	"context"
	"testing"
	"time"

	"{{.Module}}/internal/mocks"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/users/models"
	usersService "{{.Module}}/internal/users/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUsersService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewLogger()
	mockDatasource := mocks.NewMockUsersDatasource(ctrl)

	service := usersService.NewService(mockLogger, mockDatasource)

	t.Run("CreateUser_Success", func(t *testing.T) {
		ctx := context.Background()

		clerkRequest := &models.ClerkUserRequest{
			Data: models.UserData{
				ID:             "clerk_123",
				ExternalID:     "usr_456",
				FirstName:      "John",
				LastName:       "Doe",
				OrganizationID: "org_789",
				EmailAddresses: []models.EmailAddress{
					{EmailAddress: "john@example.com"},
				},
			},
			Object:    "user",
			Timestamp: 1640995200000,
			Type:      "user.created",
		}

		expectedUser := &models.User{
			ID:           "usr_generated_uuid",
			ClerkUserID:  "clerk_123",
			Email:        "john@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			CreatedAt:    time.UnixMilli(1640995200000),
			UpdatedAt:    time.UnixMilli(1640995200000),
			LastActiveAt: time.UnixMilli(1640995200000),
		}

		mockDatasource.EXPECT().CreateUser(ctx, gomock.Any()).Return(expectedUser, nil)

		// Act
		result, err := service.CreateUser(ctx, clerkRequest)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.ClerkUserID, result.ClerkUserID)
		assert.Equal(t, expectedUser.Email, result.Email)
	})

	t.Run("GetUserByClerkUserID_Success", func(t *testing.T) {
		ctx := context.Background()
		clerkUserID := "clerk_123"

		expectedUser := &models.User{
			ID:          "usr_456",
			ClerkUserID: clerkUserID,
			Email:       "john@example.com",
			FirstName:   "John",
			LastName:    "Doe",
		}

		mockDatasource.EXPECT().GetUserByClerkUserID(ctx, clerkUserID).Return(expectedUser, nil)

		// Act
		result, err := service.GetUserByClerkUserID(ctx, clerkUserID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.ClerkUserID, result.ClerkUserID)
	})

	t.Run("GetUserByClerkUserID_InvalidID", func(t *testing.T) {
		ctx := context.Background()
		clerkUserID := "" // Empty ID

		// Act
		result, err := service.GetUserByClerkUserID(ctx, clerkUserID)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.ID)
	})
}
