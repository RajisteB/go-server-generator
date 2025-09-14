package organizations_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"{{.Module}}/internal/organizations"
	"{{.Module}}/internal/organizations/models"
	"{{.Module}}/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrganizationsService_CreateOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
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

	mockDatasource.On("CreateOrganization", ctx, mock.AnythingOfType("*models.Organization")).Return(expectedOrg, nil)

	// Act
	result, err := service.CreateOrganization(ctx, clerkRequest)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrg.ID, result.ID)
	assert.Equal(t, expectedOrg.ClerkOrgID, result.ClerkOrgID)
	assert.Equal(t, expectedOrg.Name, result.Name)
	assert.Equal(t, expectedOrg.Slug, result.Slug)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_CreateOrganization_ValidationError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	// Invalid clerk request with missing required fields
	clerkRequest := &models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "", // Missing name
			Slug:      "",
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.created",
	}

	// Act
	result, err := service.CreateOrganization(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "CreateOrganization")
}

func TestOrganizationsService_CreateOrganization_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
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

	expectedError := errors.New("database error")
	mockDatasource.On("CreateOrganization", ctx, mock.AnythingOfType("*models.Organization")).Return(nil, expectedError)

	// Act
	result, err := service.CreateOrganization(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_UpdateOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Updated Organization",
			Slug:      "updated-org",
			ImageURL:  "https://example.com/updated-image.jpg",
			CreatedAt: 1640995200000,
			UpdatedAt: 1640995200000,
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.updated",
	}

	mockDatasource.On("UpdateOrganization", ctx, mock.AnythingOfType("*models.Organization")).Return(true, nil)

	// Act
	result, err := service.UpdateOrganization(ctx, clerkRequest)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_UpdateOrganization_InvalidClerkOrgID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "", // Empty clerk org ID
			Name:      "Updated Organization",
			Slug:      "updated-org",
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.updated",
	}

	// Act
	result, err := service.UpdateOrganization(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "UpdateOrganization")
}

func TestOrganizationsService_GetOrganizationByClerkOrgID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkOrgID := "clerk_org_123"
	expectedOrg := &models.Organization{
		ID:         "org_456",
		ClerkOrgID: clerkOrgID,
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDatasource.On("GetOrganizationByClerkOrgID", ctx, clerkOrgID).Return(expectedOrg, nil)

	// Act
	result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrg.ID, result.ID)
	assert.Equal(t, expectedOrg.ClerkOrgID, result.ClerkOrgID)
	assert.Equal(t, expectedOrg.Name, result.Name)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_GetOrganizationByClerkOrgID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkOrgID := "" // Empty ID

	// Act
	result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "GetOrganizationByClerkOrgID")
}

func TestOrganizationsService_GetOrganizationByID_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	orgID := "org_456"
	expectedOrg := &models.Organization{
		ID:         orgID,
		ClerkOrgID: "clerk_org_123",
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDatasource.On("GetOrganizationByID", ctx, orgID).Return(expectedOrg, nil)

	// Act
	result, err := service.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrg.ID, result.ID)
	assert.Equal(t, expectedOrg.ClerkOrgID, result.ClerkOrgID)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_GetOrganizationByID_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	orgID := "" // Empty ID

	// Act
	result, err := service.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertNotCalled(t, "GetOrganizationByID")
}

func TestOrganizationsService_DeleteOrganization_Success(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkID := "clerk_org_123"
	mockDatasource.On("DeleteOrganizationByClerkID", ctx, clerkID).Return(true, nil)

	// Act
	result, err := service.DeleteOrganization(ctx, clerkID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_DeleteOrganization_InvalidID(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkID := "" // Empty ID

	// Act
	result, err := service.DeleteOrganization(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)

	mockDatasource.AssertNotCalled(t, "DeleteOrganizationByClerkID")
}

func TestOrganizationsService_WithContext(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	clerkOrgID := "clerk_org_123"
	expectedOrg := &models.Organization{
		ID:         "org_456",
		ClerkOrgID: clerkOrgID,
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDatasource.On("GetOrganizationByClerkOrgID", ctx, clerkOrgID).Return(expectedOrg, nil)

	// Act
	result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrg.ID, result.ID)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_UpdateOrganization_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkRequest := &models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Updated Organization",
			Slug:      "updated-org",
			ImageURL:  "https://example.com/updated-image.jpg",
			CreatedAt: 1640995200000,
			UpdatedAt: 1640995200000,
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.updated",
	}

	expectedError := errors.New("database error")
	mockDatasource.On("UpdateOrganization", ctx, mock.AnythingOfType("*models.Organization")).Return(false, expectedError)

	// Act
	result, err := service.UpdateOrganization(ctx, clerkRequest)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, result)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_GetOrganizationByClerkOrgID_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkOrgID := "clerk_org_123"
	expectedError := errors.New("database error")
	mockDatasource.On("GetOrganizationByClerkOrgID", ctx, clerkOrgID).Return(nil, expectedError)

	// Act
	result, err := service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_GetOrganizationByID_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	orgID := "org_456"
	expectedError := errors.New("database error")
	mockDatasource.On("GetOrganizationByID", ctx, orgID).Return(nil, expectedError)

	// Act
	result, err := service.GetOrganizationByID(ctx, orgID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ID)

	mockDatasource.AssertExpectations(t)
}

func TestOrganizationsService_DeleteOrganization_DatasourceError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkID := "clerk_org_123"
	expectedError := errors.New("database error")
	mockDatasource.On("DeleteOrganizationByClerkID", ctx, clerkID).Return(false, expectedError)

	// Act
	result, err := service.DeleteOrganization(ctx, clerkID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.False(t, result)

	mockDatasource.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkOrganizationsService_CreateOrganization(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
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

	mockDatasource.On("CreateOrganization", ctx, mock.AnythingOfType("*models.Organization")).Return(expectedOrg, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateOrganization(ctx, clerkRequest)
	}
}

func BenchmarkOrganizationsService_GetOrganizationByClerkOrgID(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	mockDatasource := mocks.NewMockOrganizationsDatasource()
	service := organizations.NewService(mockLogger, mockDatasource)
	ctx := context.Background()

	clerkOrgID := "clerk_org_123"
	expectedOrg := &models.Organization{
		ID:         "org_456",
		ClerkOrgID: clerkOrgID,
		Name:       "Test Organization",
		Slug:       "test-org",
	}

	mockDatasource.On("GetOrganizationByClerkOrgID", ctx, clerkOrgID).Return(expectedOrg, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetOrganizationByClerkOrgID(ctx, clerkOrgID)
	}
}

// Test helper functions
func createTestClerkOrganizationRequest() *models.ClerkOrganizationRequest {
	return &models.ClerkOrganizationRequest{
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
}

func createTestOrganization() *models.Organization {
	return &models.Organization{
		ID:         "org_456",
		ClerkOrgID: "clerk_org_123",
		Name:       "Test Organization",
		Slug:       "test-org",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
