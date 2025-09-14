package health_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"{{.Module}}/internal/health"
	"{{.Module}}/internal/health/models"
	"{{.Module}}/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestHealthService_GetHealth_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database ping success
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(nil)

	// Act
	result, err := service.GetHealth(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "healthy", result.Status)
	assert.Equal(t, "healthy", result.Services["database"])
	assert.Equal(t, "1.0.0", result.Version)
	assert.NotZero(t, result.Timestamp)

	mockDB.AssertExpectations(t)
	mockSQLDB.AssertExpectations(t)
}

func TestHealthService_GetHealth_DatabaseConnectionError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database connection error
	expectedErr := errors.New("database connection failed")
	mockDB.On("DB").Return(nil, expectedErr)

	// Act
	result, err := service.GetHealth(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, result)
	assert.Equal(t, "unhealthy", result.Status)
	assert.Equal(t, "unhealthy", result.Services["database"])
	assert.NotZero(t, result.Timestamp)

	mockDB.AssertExpectations(t)
}

func TestHealthService_GetHealth_DatabasePingError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database ping failure
	expectedErr := errors.New("ping failed")
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(expectedErr)

	// Act
	result, err := service.GetHealth(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, result)
	assert.Equal(t, "unhealthy", result.Status)
	assert.Equal(t, "unhealthy", result.Services["database"])
	assert.NotZero(t, result.Timestamp)

	mockDB.AssertExpectations(t)
	mockSQLDB.AssertExpectations(t)
}

func TestHealthService_GetHealth_WithContext(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	// Mock database ping success
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(nil)

	// Act
	result, err := service.GetHealth(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "healthy", result.Status)

	mockDB.AssertExpectations(t)
	mockSQLDB.AssertExpectations(t)
}

func TestHealthService_GetHealth_TimestampAccuracy(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database ping success
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(nil)

	// Act
	start := time.Now()
	result, err := service.GetHealth(ctx)
	end := time.Now()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Check that timestamp is within reasonable bounds
	assert.True(t, result.Timestamp.After(start.Add(-time.Second)))
	assert.True(t, result.Timestamp.Before(end.Add(time.Second)))

	mockDB.AssertExpectations(t)
	mockSQLDB.AssertExpectations(t)
}

func TestHealthService_GetHealth_ServiceStructure(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database ping success
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(nil)

	// Act
	result, err := service.GetHealth(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify all required fields are present
	assert.NotEmpty(t, result.Status)
	assert.NotZero(t, result.Timestamp)
	assert.NotNil(t, result.Services)
	assert.Contains(t, result.Services, "database")
	assert.NotEmpty(t, result.Version)
	assert.NotEmpty(t, result.Uptime)

	mockDB.AssertExpectations(t)
	mockSQLDB.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkHealthService_GetHealth(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	mockSQLDB := mocks.NewMockSQLDB()
	
	service := health.NewService(mockLogger, mockDB)
	ctx := context.Background()

	// Mock database ping success
	mockDB.On("DB").Return(mockSQLDB, nil)
	mockSQLDB.On("Ping").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetHealth(ctx)
	}
}

// Test helper functions
func createTestHealthStatus() *models.HealthStatus {
	return &models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"database": "healthy",
		},
		Version: "1.0.0",
		Uptime:  "unknown",
	}
}

func createTestUnhealthyStatus() *models.HealthStatus {
	return &models.HealthStatus{
		Status:    "unhealthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"database": "unhealthy",
		},
		Version: "1.0.0",
		Uptime:  "unknown",
	}
}
