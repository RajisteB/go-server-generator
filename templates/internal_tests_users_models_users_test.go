package models_test

import (
	"testing"
	"time"

	"{{.Module}}/internal/users/models"

	"github.com/stretchr/testify/assert"
)

func TestUser_Structure(t *testing.T) {
	// Arrange
	now := time.Now()
	businessName := "Test Business"
	profileImageURL := "https://example.com/image.jpg"

	// Act
	user := models.User{
		ID:               "usr_123",
		ClerkUserID:      "clerk_456",
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		IsBusinessAcount: true,
		BusinessName:     &businessName,
		ProfileImageURL:  &profileImageURL,
		MFAEnabled:       true,
		TwoFactorEnabled: false,
		IsBanned:         false,
		LastActiveAt:     now,
		CreatedAt:        now,
		UpdatedAt:        now,
		DeletedAt:        time.Time{},
	}

	// Assert
	assert.Equal(t, "usr_123", user.ID)
	assert.Equal(t, "clerk_456", user.ClerkUserID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "org_789", user.OrganizationID)
	assert.True(t, user.IsBusinessAcount)
	assert.Equal(t, "Test Business", *user.BusinessName)
	assert.Equal(t, "https://example.com/image.jpg", *user.ProfileImageURL)
	assert.True(t, user.MFAEnabled)
	assert.False(t, user.TwoFactorEnabled)
	assert.False(t, user.IsBanned)
	assert.Equal(t, now, user.LastActiveAt)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
	assert.True(t, user.DeletedAt.IsZero())
}

func TestClerkUserRequest_ToUser(t *testing.T) {
	// Arrange
	profileImageURL := "https://example.com/profile.jpg"
	clerkRequest := models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "Jane",
			LastName:         "Smith",
			OrganizationID:   "org_789",
			ProfileImageURL:  profileImageURL,
			PasswordEnabled:  true,
			TwoFactorEnabled: true,
			IsBanned:         false,
			CreatedAt:        1640995200000, // 2022-01-01 00:00:00 UTC
			UpdatedAt:        1640995200000,
			LastSignInAt:     1640995200000,
			EmailAddresses: []models.EmailAddress{
				{
					EmailAddress: "jane@example.com",
					ID:           "email_123",
					Object:       "email_address",
					Verification: models.Verification{
						Status:   "verified",
						Strategy: "email_code",
					},
				},
			},
		},
		EventAttributes: models.EventAttributes{
			HTTPRequest: models.HTTPRequest{
				ClientIP:  "192.168.1.1",
				UserAgent: "Mozilla/5.0",
			},
		},
		Object:    "user",
		Timestamp: 1640995200000,
		Type:      "user.created",
	}

	// Act
	user := clerkRequest.ToUser()

	// Assert
	assert.Equal(t, "usr_456", user.ID)
	assert.Equal(t, "clerk_123", user.ClerkUserID)
	assert.Equal(t, "jane@example.com", user.Email)
	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, "org_789", user.OrganizationID)
	assert.Equal(t, profileImageURL, *user.ProfileImageURL)
	assert.True(t, user.MFAEnabled)
	assert.True(t, user.TwoFactorEnabled)
	assert.False(t, user.IsBanned)
	assert.Equal(t, time.UnixMilli(1640995200000), user.CreatedAt)
	assert.Equal(t, time.UnixMilli(1640995200000), user.UpdatedAt)
	assert.Equal(t, time.UnixMilli(1640995200000), user.LastActiveAt)
}

func TestClerkUserRequest_ToUser_NoEmailAddresses(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "Jane",
			LastName:         "Smith",
			OrganizationID:   "org_789",
			EmailAddresses:   []models.EmailAddress{},
		},
		EventAttributes: models.EventAttributes{},
		Object:          "user",
		Timestamp:       1640995200000,
		Type:            "user.created",
	}

	// Act
	user := clerkRequest.ToUser()

	// Assert
	assert.Equal(t, "", user.Email)
}

func TestClerkUserRequest_ToUser_MultipleEmailAddresses(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "Jane",
			LastName:         "Smith",
			OrganizationID:   "org_789",
			EmailAddresses: []models.EmailAddress{
				{
					EmailAddress: "jane@example.com",
					ID:           "email_123",
				},
				{
					EmailAddress: "jane.work@company.com",
					ID:           "email_456",
				},
			},
		},
		EventAttributes: models.EventAttributes{},
		Object:          "user",
		Timestamp:       1640995200000,
		Type:            "user.created",
	}

	// Act
	user := clerkRequest.ToUser()

	// Assert
	// Should use the first email address
	assert.Equal(t, "jane@example.com", user.Email)
}

func TestUser_Sanitize(t *testing.T) {
	// Arrange
	businessName := "<script>alert('xss')</script>Business Name"
	profileImageURL := "https://example.com/image.jpg"
	user := models.User{
		ID:               "usr_<script>alert('xss')</script>123",
		ClerkUserID:      "clerk_456",
		Email:            "TEST@EXAMPLE.COM",
		FirstName:        "John<script>alert('xss')</script>",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		BusinessName:     &businessName,
		ProfileImageURL:  &profileImageURL,
	}

	// Act
	user.Sanitize()

	// Assert
	assert.NotContains(t, user.ID, "<script>")
	assert.NotContains(t, user.FirstName, "<script>")
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotContains(t, *user.BusinessName, "<script>")
	assert.Equal(t, "https://example.com/image.jpg", *user.ProfileImageURL)
}

func TestClerkUserRequest_Sanitize(t *testing.T) {
	// Arrange
	clerkRequest := models.ClerkUserRequest{
		Object:    "<script>alert('xss')</script>user",
		Type:      "user.<script>alert('xss')</script>created",
		Timestamp: 1640995200000,
	}

	// Act
	clerkRequest.Sanitize()

	// Assert
	assert.NotContains(t, clerkRequest.Object, "<script>")
	assert.NotContains(t, clerkRequest.Type, "<script>")
}

func TestClerkUserDeleteRequest_Sanitize(t *testing.T) {
	// Arrange
	deleteRequest := models.ClerkUserDeleteRequest{
		Object:    "<script>alert('xss')</script>user",
		Type:      "user.<script>alert('xss')</script>deleted",
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
		ID:      "usr_<script>alert('xss')</script>123",
		Deleted: true,
	}

	// Act
	deleteData.Sanitize()

	// Assert
	assert.NotContains(t, deleteData.ID, "<script>")
	assert.True(t, deleteData.Deleted)
}

func TestUser_JSONSerialization(t *testing.T) {
	// Arrange
	businessName := "Test Business"
	profileImageURL := "https://example.com/image.jpg"
	user := models.User{
		ID:               "usr_123",
		ClerkUserID:      "clerk_456",
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		IsBusinessAcount: true,
		BusinessName:     &businessName,
		ProfileImageURL:  &profileImageURL,
		MFAEnabled:       true,
		TwoFactorEnabled: false,
		IsBanned:         false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Act
	jsonBytes, err := json.Marshal(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains expected fields
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "id")
	assert.Contains(t, jsonStr, "clerk_user_id")
	assert.Contains(t, jsonStr, "email")
	assert.Contains(t, jsonStr, "first_name")
	assert.Contains(t, jsonStr, "last_name")
	assert.Contains(t, jsonStr, "organization_id")
	assert.Contains(t, jsonStr, "is_business_account")
	assert.Contains(t, jsonStr, "business_name")
	assert.Contains(t, jsonStr, "profile_image_url")
	assert.Contains(t, jsonStr, "mfa_enabled")
	assert.Contains(t, jsonStr, "two_factor_enabled")
	assert.Contains(t, jsonStr, "is_banned")
	assert.Contains(t, jsonStr, "created_at")
	assert.Contains(t, jsonStr, "updated_at")
}

func TestUser_JSONDeserialization(t *testing.T) {
	// Arrange
	jsonStr := `{
		"id": "usr_123",
		"clerk_user_id": "clerk_456",
		"email": "test@example.com",
		"first_name": "John",
		"last_name": "Doe",
		"organization_id": "org_789",
		"is_business_account": true,
		"business_name": "Test Business",
		"profile_image_url": "https://example.com/image.jpg",
		"mfa_enabled": true,
		"two_factor_enabled": false,
		"is_banned": false,
		"created_at": "2023-12-25T15:30:45Z",
		"updated_at": "2023-12-25T15:30:45Z"
	}`

	// Act
	var user models.User
	err := json.Unmarshal([]byte(jsonStr), &user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "usr_123", user.ID)
	assert.Equal(t, "clerk_456", user.ClerkUserID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "org_789", user.OrganizationID)
	assert.True(t, user.IsBusinessAcount)
	assert.Equal(t, "Test Business", *user.BusinessName)
	assert.Equal(t, "https://example.com/image.jpg", *user.ProfileImageURL)
	assert.True(t, user.MFAEnabled)
	assert.False(t, user.TwoFactorEnabled)
	assert.False(t, user.IsBanned)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUser_ValidationTags(t *testing.T) {
	tests := []struct {
		name        string
		user        models.User
		expectValid bool
	}{
		{
			name: "valid user",
			user: models.User{
				ID:          "usr_123",
				ClerkUserID: "clerk_456",
				FirstName:   "John",
				LastName:    "Doe",
			},
			expectValid: true,
		},
		{
			name: "missing required ID",
			user: models.User{
				ClerkUserID: "clerk_456",
				FirstName:   "John",
				LastName:    "Doe",
			},
			expectValid: false,
		},
		{
			name: "missing required ClerkUserID",
			user: models.User{
				ID:        "usr_123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectValid: false,
		},
		{
			name: "missing required FirstName",
			user: models.User{
				ID:          "usr_123",
				ClerkUserID: "clerk_456",
				LastName:    "Doe",
			},
			expectValid: false,
		},
		{
			name: "missing required LastName",
			user: models.User{
				ID:          "usr_123",
				ClerkUserID: "clerk_456",
				FirstName:   "John",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validation.ValidateStruct(tt.user)

			// Assert
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkUser_JSONMarshal(b *testing.B) {
	businessName := "Test Business"
	profileImageURL := "https://example.com/image.jpg"
	user := models.User{
		ID:               "usr_123",
		ClerkUserID:      "clerk_456",
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		IsBusinessAcount: true,
		BusinessName:     &businessName,
		ProfileImageURL:  &profileImageURL,
		MFAEnabled:       true,
		TwoFactorEnabled: false,
		IsBanned:         false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(user)
	}
}

func BenchmarkUser_Sanitize(b *testing.B) {
	businessName := "<script>alert('xss')</script>Business Name"
	profileImageURL := "https://example.com/image.jpg"
	user := models.User{
		ID:               "usr_<script>alert('xss')</script>123",
		ClerkUserID:      "clerk_456",
		Email:            "TEST@EXAMPLE.COM",
		FirstName:        "John<script>alert('xss')</script>",
		LastName:         "Doe",
		OrganizationID:   "org_789",
		BusinessName:     &businessName,
		ProfileImageURL:  &profileImageURL,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.Sanitize()
	}
}

func BenchmarkClerkUserRequest_ToUser(b *testing.B) {
	clerkRequest := models.ClerkUserRequest{
		Data: models.UserData{
			ID:               "clerk_123",
			ExternalID:       "usr_456",
			FirstName:        "Jane",
			LastName:         "Smith",
			OrganizationID:   "org_789",
			ProfileImageURL:  "https://example.com/profile.jpg",
			PasswordEnabled:  true,
			TwoFactorEnabled: true,
			IsBanned:         false,
			CreatedAt:        1640995200000,
			UpdatedAt:        1640995200000,
			LastSignInAt:     1640995200000,
			EmailAddresses: []models.EmailAddress{
				{
					EmailAddress: "jane@example.com",
					ID:           "email_123",
					Object:       "email_address",
					Verification: models.Verification{
						Status:   "verified",
						Strategy: "email_code",
					},
				},
			},
		},
		EventAttributes: models.EventAttributes{},
		Object:          "user",
		Timestamp:       1640995200000,
		Type:            "user.created",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clerkRequest.ToUser()
	}
}
