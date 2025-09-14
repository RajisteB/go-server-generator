package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"{{.Module}}/internal/users"
	"{{.Module}}/internal/users/models"
	"{{.Module}}/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsersService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
			OrganizationID:   "org_789",
			EmailAddresses: []models.EmailAddress{
				{EmailAddress: "john@example.com"},
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}

	expectedUser := &models.User{
		ID:               "usr_generated_uuid",
		ClerkUserID:      "clerk_123",
		Email:            "john@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		CreatedAt:        time.UnixMilli(1640995200000),
		UpdatedAt:        time.UnixMilli(1640995200000),
		LastActiveAt:     time.UnixMilli(1640995200000),
	}

	mockDatasource.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)

	// Act
	result, err := service.CreateUser(ctx, clerkRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.ClerkUserID, result.ClerkUserID)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.FirstName, result.FirstName)
	assert.Equal(t, expectedUser.LastName, result.LastName)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_CreateUser_ValidationError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	// Invalid clerk request with missing required fields
	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "", // Missing external ID
			FirstName:        "", // Missing first name
			LastName:         "", // Missing last name
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}

	// Act
	result, err := service.CreateUser(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "CreateUser")
}

func TestUsersService_CreateUser_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
			EmailAddresses: []models.EmailAddress{
				{EmailAddress: "john@example.com"},
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}

	expectedError := errors.New("database error")
	mockDatasource.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil, expectedError)

	// Act
	result, err := service.CreateUser(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
			EmailAddresses: []models.EmailAddress{
				{EmailAddress: "john@example.com"},
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.updated",
	}

	mockDatasource.On("UpdateUser", ctx, mock.AnythingOfType("*models.User")).Return(true, nil)

	// Act
	result, err := service.UpdateUser(ctx, clerkRequest)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_UpdateUser_InvalidClerkUserID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "", // Empty clerk user ID
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.updated",
	}

	// Act
	result, err := service.UpdateUser(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "UpdateUser")
}

func TestUsersService_GetUserByClerkUserID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "clerk_123"
	expectedUser := &models.User{
		ID:          "usr_456",
		ClerkUserID: clerkUserID,
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDatasource.On("GetUserByClerkUserID", ctx, clerkUserID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.ClerkUserID, result.ClerkUserID)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_GetUserByClerkUserID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "" // Empty ID

	// Act
	result, err := service.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "GetUserByClerkUserID")
}

func TestUsersService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	userID := "usr_456"
	expectedUser := &models.User{
		ID:          userID,
		ClerkUserID: "clerk_123",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDatasource.On("GetUserByID", ctx, userID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.ClerkUserID, result.ClerkUserID)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_GetUserByID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	userID := "" // Empty ID

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "GetUserByID")
}

func TestUsersService_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkID := "clerk_123"
	mockDatasource.On("DeleteUserByClerkID", ctx, clerkID).Return(true, nil)

	// Act
	result, err := service.DeleteUser(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_DeleteUser_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkID := "" // Empty ID

	// Act
	result, err := service.DeleteUser(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "DeleteUserByClerkID")
}

func TestUsersService_UpdateUserOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "clerk_123"
	orgID := "org_456"
	mockDatasource.On("UpdateUserOrganization", ctx, clerkUserID, orgID).Return(true, nil)

	// Act
	result, err := service.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestUsersService_UpdateUserOrganization_InvalidClerkUserID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "" // Empty ID
	orgID := "org_456"

	// Act
	result, err := service.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "UpdateUserOrganization")
}

func TestUsersService_UpdateUserOrganization_InvalidOrgID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "clerk_123"
	orgID := "" // Empty ID

	// Act
	result, err := service.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "UpdateUserOrganization")
}

func TestUsersService_WithContext(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	clerkUserID := "clerk_123"
	expectedUser := &models.User{
		ID:          "usr_456",
		ClerkUserID: clerkUserID,
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDatasource.On("GetUserByClerkUserID", ctx, clerkUserID).Return(expectedUser, nil)

	// Act
	result, err := service.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	mockDatasource.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkUsersService_CreateUser(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
			EmailAddresses: []models.EmailAddress{
				{EmailAddress: "john@example.com"},
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}

	expectedUser := &models.User{
		ID:          "usr_generated_uuid",
		ClerkUserID: "clerk_123",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDatasource.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateUser(ctx, clerkRequest)
	}
}

func BenchmarkUsersService_GetUserByClerkUserID(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockUsersDatasource()
	service := users.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkUserID := "clerk_123"
	expectedUser := &models.User{
		ID:          "usr_456",
		ClerkUserID: clerkUserID,
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDatasource.On("GetUserByClerkUserID", ctx, clerkUserID).Return(expectedUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserByClerkUserID(ctx, clerkUserID)
	}
}

// Test helper functions
func createTestClerkUserRequest() *models.ClerkUserRequest {
	return &models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "John",
			LastName:         "Doe",
			OrganizationID:   "org_789",
			EmailAddresses: []models.EmailAddress{
				{EmailAddress: "john@example.com"},
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}
}

func createTestUser() *models.User {
	return &models.User{
		ID:          "usr_456",
		ClerkUserID: "clerk_123",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
