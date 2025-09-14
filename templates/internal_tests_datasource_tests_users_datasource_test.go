package datasource_tests

import (
	"context"
	"testing"
	"time"

	"{{.Module}}/internal/mocks"
	"{{.Module}}/internal/shared/logger"
	usersDatasource "{{.Module}}/internal/users/datasource"
	"{{.Module}}/internal/users/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUsersDatasource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewLogger()
	mockDB := mocks.NewMockGormDB(ctrl)

	datasource := usersDatasource.NewDatasource(mockLogger, mockDB)

	t.Run("CreateUser_Success", func(t *testing.T) {
		ctx := context.Background()

		user := &models.User{
			ID:          "usr_123",
			ClerkUserID: "clerk_456",
			Email:       "john@example.com",
			FirstName:   "John",
			LastName:    "Doe",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Create(user).Return(mockDB)
		mockDB.EXPECT().Error().Return(nil)

		// Act
		result, err := datasource.CreateUser(ctx, user)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.ClerkUserID, result.ClerkUserID)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("GetUserByClerkUserID_Success", func(t *testing.T) {
		ctx := context.Background()
		clerkUserID := "clerk_456"

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Where("clerk_user_id = ?", clerkUserID).Return(mockDB)
		mockDB.EXPECT().First(gomock.Any()).Return(mockDB)
		mockDB.EXPECT().Error().Return(nil)

		// Act
		result, err := datasource.GetUserByClerkUserID(ctx, clerkUserID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetUserByClerkUserID_NotFound", func(t *testing.T) {
		ctx := context.Background()
		clerkUserID := "clerk_456"

		mockDB.EXPECT().WithContext(ctx).Return(mockDB)
		mockDB.EXPECT().Where("clerk_user_id = ?", clerkUserID).Return(mockDB)
		mockDB.EXPECT().First(gomock.Any()).Return(mockDB)
		mockDB.EXPECT().Error().Return(gorm.ErrRecordNotFound)

		// Act
		result, err := datasource.GetUserByClerkUserID(ctx, clerkUserID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("GetUserByClerkUserID_InvalidID", func(t *testing.T) {
		ctx := context.Background()
		clerkUserID := "" // Empty ID

		// Act
		result, err := datasource.GetUserByClerkUserID(ctx, clerkUserID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
