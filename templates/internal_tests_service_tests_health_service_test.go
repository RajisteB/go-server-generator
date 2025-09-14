package service_tests

import (
	"context"
	"testing"

	healthService "{{.Module}}/internal/health/service"
	"{{.Module}}/internal/mocks"
	"{{.Module}}/internal/shared/logger"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHealthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewLogger()
	mockDB := mocks.NewMockGormDB(ctrl)
	mockSQLDB := mocks.NewMockSQLDB(ctrl)

	service := healthService.NewService(mockLogger, mockDB)

	t.Run("GetHealth_Success", func(t *testing.T) {
		ctx := context.Background()

		// Mock database ping success
		mockDB.EXPECT().DB().Return(mockSQLDB, nil)
		mockSQLDB.EXPECT().Ping().Return(nil)

		// Act
		result, err := service.GetHealth(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "healthy", result.Status)
		assert.Equal(t, "healthy", result.Services["database"])
		assert.Equal(t, "1.0.0", result.Version)
		assert.NotZero(t, result.Timestamp)
	})

	t.Run("GetHealth_DatabaseError", func(t *testing.T) {
		ctx := context.Background()

		// Mock database ping failure
		mockDB.EXPECT().DB().Return(mockSQLDB, nil)
		mockSQLDB.EXPECT().Ping().Return(assert.AnError)

		// Act
		result, err := service.GetHealth(ctx)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "unhealthy", result.Status)
		assert.Equal(t, "unhealthy", result.Services["database"])
	})
}
