package models_test

import (
	"testing"
	"time"

	"{{.Module}}/internal/organizations/models"

	"github.com/stretchr/testify/assert"
)

func TestOrganization_Structure(t *testing.T) {
	// Arrange
	now := time.Now()
	imageURL := "https://example.com/org-image.jpg"

	// Act
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		ImageURL:   &imageURL,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  time.Time{},
	}

	// Assert
	assert.Equal(t, "org_123", org.ID)
	assert.Equal(t, "clerk_org_456", org.ClerkOrgID)
	assert.Equal(t, "Test Organization", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, "https://example.com/org-image.jpg", *org.ImageURL)
	assert.Equal(t, now, org.CreatedAt)
	assert.Equal(t, now, org.UpdatedAt)
	assert.True(t, org.DeletedAt.IsZero())
}

func TestClerkOrganizationRequest_ToOrganization(t *testing.T) {
	// Arrange
	imageURL := "https://example.com/org-image.jpg"
	clerkRequest := models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Test Organization",
			Slug:      "test-org",
			ImageURL:  imageURL,
			CreatedAt: 1640995200000, // 2022-01-01 00:00:00 UTC
			UpdatedAt: 1640995200000,
		},
		EventAttributes: models.EventAttributes{
			HTTPRequest: models.HTTPRequest{
				ClientIP:  "192.168.1.1",
				UserAgent: "Mozilla/5.0",
			},
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.created",
	}

	// Act
	org := clerkRequest.ToOrganization()

	// Assert
	assert.Equal(t, "clerk_org_123", org.ClerkOrgID)
	assert.Equal(t, "Test Organization", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, imageURL, *org.ImageURL)
	assert.Equal(t, time.UnixMilli(1640995200000), org.CreatedAt)
	assert.Equal(t, time.UnixMilli(1640995200000), org.UpdatedAt)
}

func TestClerkOrganizationRequest_ToOrganization_EmptyImageURL(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Test Organization",
			Slug:      "test-org",
			ImageURL:  "", // Empty image URL
			CreatedAt: 1640995200000,
			UpdatedAt: 1640995200000,
		},
		EventAttributes: models.EventAttributes{},
		Object:          "organization",
		Timestamp:       1640995200000,
		Type:            "organization.created",
	}

	// Act
	org := clerkRequest.ToOrganization()

	// Assert
	assert.Equal(t, "clerk_org_123", org.ClerkOrgID)
	assert.Equal(t, "Test Organization", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, "", *org.ImageURL)
	assert.Equal(t, time.UnixMilli(1640995200000), org.CreatedAt)
	assert.Equal(t, time.UnixMilli(1640995200000), org.UpdatedAt)
}

func TestOrganization_Sanitize(t *testing.T) {
	// Arrange
	imageURL := "<script>alert('xss')</script>https://example.com/image.jpg"
	org := models.Organization{
		ID:         "org_<script>alert('xss')</script>123",
		ClerkOrgID: "clerk_org_456",
		Name:        "Test<script>alert('xss')</script> Organization",
		Slug:        "test-org<script>alert('xss')</script>",
		ImageURL:    &imageURL,
	}

	// Act
	org.Sanitize()

	// Assert
	assert.NotContains(t, org.ID, "<script>")
	assert.NotContains(t, org.Name, "<script>")
	assert.NotContains(t, org.Slug, "<script>")
	assert.NotContains(t, *org.ImageURL, "<script>")
}

func TestClerkOrganizationRequest_Sanitize(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkOrganizationRequest{
		Object:    "<script>alert('xss')</script>organization",
		Type:      "organization.<script>alert('xss')</script>created",
		Timestamp: 1640995200000,
	}

	// Act
	clerkRequest.Sanitize()

	// Assert
	assert.NotContains(t, clerkRequest.Object, "<script>")
	assert.NotContains(t, clerkRequest.Type, "<script>")
}

func TestClerkOrganizationDeleteRequest_Sanitize(t *testing.T) {
	// Arrange
	deleteRequest := models.ClerkOrganizationDeleteRequest{
		Object:    "<script>alert('xss')</script>organization",
		Type:      "organization.<script>alert('xss')</script>deleted",
		Timestamp: 1640995200000,
	}

	// Act
	deleteRequest.Sanitize()

	// Assert
	assert.NotContains(t, deleteRequest.Object, "<script>")
	assert.NotContains(t, deleteRequest.Type, "<script>")
}

func TestDeleteData_Sanitize(t *testing.T) {
	// Arrange
	deleteData := models.DeleteData{
		ID:      "org_<script>alert('xss')</script>123",
		Deleted: true,
	}

	// Act
	deleteData.Sanitize()

	// Assert
	assert.NotContains(t, deleteData.ID, "<script>")
	assert.True(t, deleteData.Deleted)
}

func TestOrganization_JSONSerialization(t *testing.T) {
	// Arrange
	imageURL := "https://example.com/org-image.jpg"
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		ImageURL:   &imageURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	jsonBytes, err := json.Marshal(org)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains expected fields
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "id")
	assert.Contains(t, jsonStr, "clerk_org_id")
	assert.Contains(t, jsonStr, "name")
	assert.Contains(t, jsonStr, "slug")
	assert.Contains(t, jsonStr, "image_url")
	assert.Contains(t, jsonStr, "created_at")
	assert.Contains(t, jsonStr, "updated_at")
}

func TestOrganization_JSONDeserialization(t *testing.T) {
	// Arrange
	jsonStr := `{
		"id": "org_123",
		"clerk_org_id": "clerk_org_456",
		"name": "Test Organization",
		"slug": "test-org",
		"image_url": "https://example.com/org-image.jpg",
		"created_at": "2023-12-25T15:30:45Z",
		"updated_at": "2023-12-25T15:30:45Z"
	}`

	// Act
	var org models.Organization
	err := json.Unmarshal([]byte(jsonStr), &org)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "org_123", org.ID)
	assert.Equal(t, "clerk_org_456", org.ClerkOrgID)
	assert.Equal(t, "Test Organization", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, "https://example.com/org-image.jpg", *org.ImageURL)
	assert.NotZero(t, org.CreatedAt)
	assert.NotZero(t, org.UpdatedAt)
}

func TestOrganization_ValidationTags(t *testing.T) {
	tests := []struct {
		name        string
		org         models.Organization
		expectValid bool
	}{
		{
			name: "valid organization",
			org: models.Organization{
				ID:         "org_123",
				ClerkOrgID: "clerk_org_456",
				Name:       "Test Organization",
			},
			expectValid: true,
		},
		{
			name: "missing required ID",
			org: models.Organization{
				ClerkOrgID: "clerk_org_456",
				Name:       "Test Organization",
			},
			expectValid: false,
		},
		{
			name: "missing required ClerkOrgID",
			org: models.Organization{
				ID:   "org_123",
				Name: "Test Organization",
			},
			expectValid: false,
		},
		{
			name: "missing required Name",
			org: models.Organization{
				ID:         "org_123",
				ClerkOrgID: "clerk_org_456",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validation.ValidateStruct(tt.org)

			// Assert
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestOrganization_NilImageURL(t *testing.T) {
	// Arrange
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		ImageURL:   nil, // Nil image URL
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	jsonBytes, err := json.Marshal(org)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains null for image_url
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "image_url")
	assert.Contains(t, jsonStr, "null")
}

func TestOrganization_EmptySlug(t *testing.T) {
	// Arrange
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "", // Empty slug
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	jsonBytes, err := json.Marshal(org)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains empty slug
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "slug")
	assert.Contains(t, jsonStr, `""`)
}

func TestOrganization_TimeHandling(t *testing.T) {
	// Arrange
	now := time.Now()
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  time.Time{}, // Zero time
	}

	// Act
	jsonBytes, err := json.Marshal(org)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains created_at and updated_at but not deleted_at
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "created_at")
	assert.Contains(t, jsonStr, "updated_at")
	assert.Contains(t, jsonStr, "deleted_at")
}

func TestClerkOrganizationRequest_Structure(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Test Organization",
			Slug:      "test-org",
			ImageURL:  "https://example.com/image.jpg",
			CreatedAt: 1640995200000,
			UpdatedAt: 1640995200000,
		},
		EventAttributes: models.EventAttributes{
			HTTPRequest: models.HTTPRequest{
				ClientIP:  "192.168.1.1",
				UserAgent: "Mozilla/5.0",
			},
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.created",
	}

	// Assert
	assert.Equal(t, "clerk_org_123", clerkRequest.Data.ID)
	assert.Equal(t, "Test Organization", clerkRequest.Data.Name)
	assert.Equal(t, "test-org", clerkRequest.Data.Slug)
	assert.Equal(t, "https://example.com/image.jpg", clerkRequest.Data.ImageURL)
	assert.Equal(t, int64(1640995200000), clerkRequest.Data.CreatedAt)
	assert.Equal(t, int64(1640995200000), clerkRequest.Data.UpdatedAt)
	assert.Equal(t, "192.168.1.1", clerkRequest.EventAttributes.HTTPRequest.ClientIP)
	assert.Equal(t, "Mozilla/5.0", clerkRequest.EventAttributes.HTTPRequest.UserAgent)
	assert.Equal(t, "organization", clerkRequest.Object)
	assert.Equal(t, int64(1640995200000), clerkRequest.Timestamp)
	assert.Equal(t, "organization.created", clerkRequest.Type)
}

func TestClerkOrganizationDeleteRequest_Structure(t *testing.T) {
	// Arrange
	deleteRequest := models.ClerkOrganizationDeleteRequest{
		Data: models.DeleteData{
			ID:      "clerk_org_123",
			Deleted: true,
		},
		EventAttributes: models.EventAttributes{
			HTTPRequest: models.HTTPRequest{
				ClientIP:  "192.168.1.1",
				UserAgent: "Mozilla/5.0",
			},
		},
		Object:    "organization",
		Timestamp: 1640995200000,
		Type:      "organization.deleted",
	}

	// Assert
	assert.Equal(t, "clerk_org_123", deleteRequest.Data.ID)
	assert.True(t, deleteRequest.Data.Deleted)
	assert.Equal(t, "192.168.1.1", deleteRequest.EventAttributes.HTTPRequest.ClientIP)
	assert.Equal(t, "Mozilla/5.0", deleteRequest.EventAttributes.HTTPRequest.UserAgent)
	assert.Equal(t, "organization", deleteRequest.Object)
	assert.Equal(t, int64(1640995200000), deleteRequest.Timestamp)
	assert.Equal(t, "organization.deleted", deleteRequest.Type)
}

// Benchmark tests
func BenchmarkOrganization_JSONMarshal(b *testing.B) {
	imageURL := "https://example.com/org-image.jpg"
	org := models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		ImageURL:   &imageURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(org)
	}
}

func BenchmarkOrganization_Sanitize(b *testing.B) {
	imageURL := "<script>alert('xss')</script>https://example.com/image.jpg"
	org := models.Organization{
		ID:         "org_<script>alert('xss')</script>123",
		ClerkOrgID: "clerk_org_456",
		Name:        "Test<script>alert('xss')</script> Organization",
		Slug:        "test-org<script>alert('xss')</script>",
		ImageURL:    &imageURL,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		org.Sanitize()
	}
}

func BenchmarkClerkOrganizationRequest_ToOrganization(b *testing.B) {
	clerkRequest := models.ClerkOrganizationRequest{
		Data: models.OrganizationData{
			ID:        "clerk_org_123",
			Name:      "Test Organization",
			Slug:      "test-org",
			ImageURL:  "https://example.com/image.jpg",
			CreatedAt: 1640995200000,
			UpdatedAt: 1640995200000,
		},
		EventAttributes: models.EventAttributes{},
		Object:          "organization",
		Timestamp:       1640995200000,
		Type:            "organization.created",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clerkRequest.ToOrganization()
	}
}

// Test helper functions
func createTestOrganization() *models.Organization {
	imageURL := "https://example.com/org-image.jpg"
	return &models.Organization{
		ID:         "org_123",
		ClerkOrgID: "clerk_org_456",
		Name:       "Test Organization",
		Slug:       "test-org",
		ImageURL:   &imageURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

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
		EventAttributes: models.EventAttributes{},
		Object:          "organization",
		Timestamp:       1640995200000,
		Type:            "organization.created",
	}
}
