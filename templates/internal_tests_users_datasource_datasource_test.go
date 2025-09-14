package datasource_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"{{.Module}}/internal/users/datasource"
	"{{.Module}}/internal/users/models"
	"{{.Module}}/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsersDatasource_CreateUser_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
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

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", user).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.CreateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.ClerkUserID, result.ClerkUserID)
	assert.Equal(t, user.Email, result.Email)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_CreateUser_GenerateUUID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	user := &models.User{
		ID:          "", // Empty ID to trigger UUID generation
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.CreateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.True(t, len(result.ID) > 0)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_CreateUser_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
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

	expectedError := errors.New("database error")
	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", user).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.CreateUser(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	user := &models.User{
		ID:          "usr_123",
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		UpdatedAt:   time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", user).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", "clerk_456").Return(mockDB)
	mockDB.On("Updates", user).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	// Act
	result, err := ds.UpdateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_UpdateUser_NoRowsAffected(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	user := &models.User{
		ID:          "usr_123",
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		UpdatedAt:   time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", user).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", "clerk_456").Return(mockDB)
	mockDB.On("Updates", user).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	// Act
	result, err := ds.UpdateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_UpdateUser_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	user := &models.User{
		ID:          "usr_123",
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		UpdatedAt:   time.Now(),
	}

	expectedError := errors.New("database error")
	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", user).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", "clerk_456").Return(mockDB)
	mockDB.On("Updates", user).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.UpdateUser(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_GetUserByClerkUserID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"
	expectedUser := &models.User{
		ID:          "usr_123",
		ClerkUserID: clerkUserID,
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_GetUserByClerkUserID_NotFound(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrRecordNotFound)

	// Act
	result, err := ds.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_GetUserByClerkUserID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "" // Empty ID

	// Act
	result, err := ds.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestUsersDatasource_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	userID := "usr_123"
	expectedUser := &models.User{
		ID:          userID,
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "id = ?", userID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	userID := "usr_123"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "id = ?", userID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrRecordNotFound)

	// Act
	result, err := ds.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_GetUserByID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	userID := "" // Empty ID

	// Act
	result, err := ds.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestUsersDatasource_DeleteUserByClerkID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "clerk_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkID).Return(mockDB)
	mockDB.On("Delete", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	// Act
	result, err := ds.DeleteUserByClerkID(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_DeleteUserByClerkID_NoRowsAffected(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "clerk_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkID).Return(mockDB)
	mockDB.On("Delete", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	// Act
	result, err := ds.DeleteUserByClerkID(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_DeleteUserByClerkID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "" // Empty ID

	// Act
	result, err := ds.DeleteUserByClerkID(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestUsersDatasource_UpdateUserOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"
	orgID := "org_789"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("Update", "organization_id", orgID).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	// Act
	result, err := ds.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_UpdateUserOrganization_NoRowsAffected(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"
	orgID := "org_789"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("Update", "organization_id", orgID).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	// Act
	result, err := ds.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestUsersDatasource_UpdateUserOrganization_InvalidClerkUserID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "" // Empty ID
	orgID := "org_789"

	// Act
	result, err := ds.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestUsersDatasource_UpdateUserOrganization_InvalidOrgID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"
	orgID := "" // Empty ID

	// Act
	result, err := ds.UpdateUserOrganization(ctx, clerkUserID, orgID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestUsersDatasource_WithContext(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	clerkUserID := "clerk_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetUserByClerkUserID(ctx, clerkUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkUsersDatasource_CreateUser(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
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

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", user).Return(mockDB)
	mockDB.On("Error").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ds.CreateUser(ctx, user)
	}
}

func BenchmarkUsersDatasource_GetUserByClerkUserID(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkUserID := "clerk_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_user_id = ?", clerkUserID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ds.GetUserByClerkUserID(ctx, clerkUserID)
	}
}

// Test helper functions
func createTestUser() *models.User {
	return &models.User{
		ID:          "usr_123",
		ClerkUserID: "clerk_456",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
