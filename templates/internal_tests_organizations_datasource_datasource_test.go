package datasource_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"{{.Module}}/internal/organizations/datasource"
	"{{.Module}}/internal/organizations/models"
	"{{.Module}}/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestOrganizationsDatasource_CreateOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", org).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.CreateOrganization(ctx, org)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, org.ID, result.ID)
	assert.Equal(t, org.ClerkOrgID, result.ClerkOrgID)
	assert.Equal(t, org.Name, result.Name)
	assert.Equal(t, org.Slug, result.Slug)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_CreateOrganization_GenerateUUID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "", // Empty ID to trigger UUID generation
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.CreateOrganization(ctx, org)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.True(t, len(result.ID) > 0)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_CreateOrganization_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	expectedError := errors.New("database error")
	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", org).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.CreateOrganization(ctx, org)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_UpdateOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Updated Organization",
		Slug:       "updated-org",
		UpdatedAt:  time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", org).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", "clerk_org_456").Return(mockDB)
	mockDB.On("Updates", org).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	// Act
	result, err := ds.UpdateOrganization(ctx, org)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_UpdateOrganization_NoRowsAffected(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Updated Organization",
		Slug:       "updated-org",
		UpdatedAt:  time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", org).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", "clerk_org_456").Return(mockDB)
	mockDB.On("Updates", org).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	// Act
	result, err := ds.UpdateOrganization(ctx, org)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_UpdateOrganization_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Updated Organization",
		Slug:       "updated-org",
		UpdatedAt:  time.Now(),
	}

	expectedError := errors.New("database error")
	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Model", org).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", "clerk_org_456").Return(mockDB)
	mockDB.On("Updates", org).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.UpdateOrganization(ctx, org)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByClerkOrgID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkOrgID := "clerk_org_456"
	expectedOrg := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: clerkOrgID,
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkOrgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByClerkOrgID_NotFound(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkOrgID := "clerk_org_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkOrgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrRecordNotFound)

	// Act
	result, err := ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByClerkOrgID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkOrgID := "" // Empty ID

	// Act
	result, err := ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestOrganizationsDatasource_GetOrganizationByID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	orgID := "org_123"
	expectedOrg := &models.Organization{
		ID:         orgID,
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "id = ?", orgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByID_NotFound(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	orgID := "org_123"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "id = ?", orgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrRecordNotFound)

	// Act
	result, err := ds.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	orgID := "" // Empty ID

	// Act
	result, err := ds.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestOrganizationsDatasource_DeleteOrganizationByClerkID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "clerk_org_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkID).Return(mockDB)
	mockDB.On("Delete", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	// Act
	result, err := ds.DeleteOrganizationByClerkID(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_DeleteOrganizationByClerkID_NoRowsAffected(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "clerk_org_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkID).Return(mockDB)
	mockDB.On("Delete", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	// Act
	result, err := ds.DeleteOrganizationByClerkID(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_DeleteOrganizationByClerkID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "" // Empty ID

	// Act
	result, err := ds.DeleteOrganizationByClerkID(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDB.AssertNotCalled(t, "WithContext")
}

func TestOrganizationsDatasource_WithContext(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	clerkOrgID := "clerk_org_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkOrgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Act
	result, err := ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByClerkOrgID_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkOrgID := "clerk_org_456"
	expectedError := errors.New("database error")

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkOrgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_GetOrganizationByID_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	orgID := "org_123"
	expectedError := errors.New("database error")

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "id = ?", orgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestOrganizationsDatasource_DeleteOrganizationByClerkID_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkID := "clerk_org_456"
	expectedError := errors.New("database error")

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkID).Return(mockDB)
	mockDB.On("Delete", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(expectedError)

	// Act
	result, err := ds.DeleteOrganizationByClerkID(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, result)

	mockDB.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkOrganizationsDatasource_CreateOrganization(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	org := &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Create", org).Return(mockDB)
	mockDB.On("Error").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ds.CreateOrganization(ctx, org)
	}
}

func BenchmarkOrganizationsDatasource_GetOrganizationByClerkOrgID(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDB := mocks.NewMockGormDB()
	ds := datasource.NewDatasource(mockLogger, mockDB)
	ctx := context.Background()

	clerkOrgID := "clerk_org_456"

	mockDB.On("WithContext", ctx).Return(mockDB)
	mockDB.On("Where", "clerk_org_id = ?", clerkOrgID).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Organization")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ds.GetOrganizationByClerkOrgID(ctx, clerkOrgID)
	}
}

// Test helper functions
func createTestOrganization() *models.Organization {
	return &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
