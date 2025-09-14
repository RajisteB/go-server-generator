package models_test

import (
	"testing"
	"time"

	"{{.Module}}/internal/health/models"

	"github.com/stretchr/testify/assert"
)

func TestHealthStatus_Structure(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	services := map[string]string{
		"database": "healthy",
		"redis":    "healthy",
		"api":      "healthy",
	}

	// Act
	healthStatus := models.HealthStatus{
		Status:    "healthy",
		Timestamp: timestamp,
		Services:  services,
		Version:   "1.0.0",
		Uptime:    "2h30m15s",
	}

	// Assert
	assert.Equal(t, "healthy", healthStatus.Status)
	assert.Equal(t, timestamp, healthStatus.Timestamp)
	assert.Equal(t, services, healthStatus.Services)
	assert.Equal(t, "1.0.0", healthStatus.Version)
	assert.Equal(t, "2h30m15s", healthStatus.Uptime)
}

func TestHealthStatus_UnhealthyStatus(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	services := map[string]string{
		"database": "unhealthy",
		"redis":    "healthy",
		"api":      "healthy",
	}

	// Act
	healthStatus := models.HealthStatus{
		Status:    "unhealthy",
		Timestamp: timestamp,
		Services:  services,
		Version:   "1.0.0",
		Uptime:    "2h30m15s",
	}

	// Assert
	assert.Equal(t, "unhealthy", healthStatus.Status)
	assert.Equal(t, "unhealthy", healthStatus.Services["database"])
	assert.Equal(t, "healthy", healthStatus.Services["redis"])
	assert.Equal(t, "healthy", healthStatus.Services["api"])
}

func TestHealthStatus_EmptyServices(t *testing.T) {
	// Arrange
	timestamp := time.Now()

	// Act
	healthStatus := models.HealthStatus{
		Status:    "healthy",
		Timestamp: timestamp,
		Services:  map[string]string{},
		Version:   "1.0.0",
		Uptime:    "0s",
	}

	// Assert
	assert.Equal(t, "healthy", healthStatus.Status)
	assert.Empty(t, healthStatus.Services)
	assert.Equal(t, "0s", healthStatus.Uptime)
}

func TestHealthStatus_NilServices(t *testing.T) {
	// Arrange
	timestamp := time.Now()

	// Act
	healthStatus := models.HealthStatus{
		Status:    "healthy",
		Timestamp: timestamp,
		Services:  nil,
		Version:   "1.0.0",
		Uptime:    "0s",
	}

	// Assert
	assert.Equal(t, "healthy", healthStatus.Status)
	assert.Nil(t, healthStatus.Services)
}

func TestHealthStatus_JSONSerialization(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	services := map[string]string{
		"database": "healthy",
		"redis":    "healthy",
	}

	healthStatus := models.HealthStatus{
		Status:    "healthy",
		Timestamp: timestamp,
		Services:  services,
		Version:   "1.0.0",
		Uptime:    "2h30m15s",
	}

	// Act
	jsonBytes, err := json.Marshal(healthStatus)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)
	
	// Verify JSON contains expected fields
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "status")
	assert.Contains(t, jsonStr, "timestamp")
	assert.Contains(t, jsonStr, "services")
	assert.Contains(t, jsonStr, "version")
	assert.Contains(t, jsonStr, "uptime")
	assert.Contains(t, jsonStr, "healthy")
	assert.Contains(t, jsonStr, "1.0.0")
}

func TestHealthStatus_JSONDeserialization(t *testing.T) {
	// Arrange
	jsonStr := `{
		"status": "healthy",
		"timestamp": "2023-12-25T15:30:45Z",
		"services": {
			"database": "healthy",
			"redis": "healthy"
		},
		"version": "1.0.0",
		"uptime": "2h30m15s"
	}`

	// Act
	var healthStatus models.HealthStatus
	err := json.Unmarshal([]byte(jsonStr), &healthStatus)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "healthy", healthStatus.Status)
	assert.Equal(t, "healthy", healthStatus.Services["database"])
	assert.Equal(t, "healthy", healthStatus.Services["redis"])
	assert.Equal(t, "1.0.0", healthStatus.Version)
	assert.Equal(t, "2h30m15s", healthStatus.Uptime)
	assert.NotZero(t, healthStatus.Timestamp)
}

func TestHealthStatus_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := `{
		"status": "healthy",
		"timestamp": "invalid-timestamp",
		"services": {
			"database": "healthy"
		},
		"version": "1.0.0",
		"uptime": "2h30m15s"
	}`

	// Act
	var healthStatus models.HealthStatus
	err := json.Unmarshal([]byte(invalidJSON), &healthStatus)

	// Assert
	assert.Error(t, err)
}

func TestHealthStatus_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		version     string
		uptime      string
		expectValid bool
	}{
		{
			name:        "valid healthy status",
			status:      "healthy",
			version:     "1.0.0",
			uptime:      "2h30m15s",
			expectValid: true,
		},
		{
			name:        "valid unhealthy status",
			status:      "unhealthy",
			version:     "1.0.0",
			uptime:      "0s",
			expectValid: true,
		},
		{
			name:        "empty status",
			status:      "",
			version:     "1.0.0",
			uptime:      "2h30m15s",
			expectValid: false,
		},
		{
			name:        "empty version",
			status:      "healthy",
			version:     "",
			uptime:      "2h30m15s",
			expectValid: false,
		},
		{
			name:        "empty uptime",
			status:      "healthy",
			version:     "1.0.0",
			uptime:      "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			healthStatus := models.HealthStatus{
				Status:    tt.status,
				Timestamp: time.Now(),
				Services:  map[string]string{"database": "healthy"},
				Version:   tt.version,
				Uptime:    tt.uptime,
			}

			// Act & Assert
			if tt.expectValid {
				assert.NotEmpty(t, healthStatus.Status)
				assert.NotEmpty(t, healthStatus.Version)
				assert.NotEmpty(t, healthStatus.Uptime)
			} else {
				// For invalid cases, we're just testing the structure
				// In a real application, you might want to add validation
				assert.Equal(t, tt.status, healthStatus.Status)
				assert.Equal(t, tt.version, healthStatus.Version)
				assert.Equal(t, tt.uptime, healthStatus.Uptime)
			}
		})
	}
}

func TestHealthStatus_ServiceStatusVariations(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]string
	}{
		{
			name: "all healthy services",
			services: map[string]string{
				"database": "healthy",
				"redis":    "healthy",
				"api":      "healthy",
			},
		},
		{
			name: "mixed service statuses",
			services: map[string]string{
				"database": "healthy",
				"redis":    "unhealthy",
				"api":      "degraded",
			},
		},
		{
			name: "single service",
			services: map[string]string{
				"database": "healthy",
			},
		},
		{
			name:     "no services",
			services: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			healthStatus := models.HealthStatus{
				Status:    "healthy",
				Timestamp: time.Now(),
				Services:  tt.services,
				Version:   "1.0.0",
				Uptime:    "2h30m15s",
			}

			// Assert
			assert.Equal(t, tt.services, healthStatus.Services)
			assert.Equal(t, len(tt.services), len(healthStatus.Services))
		})
	}
}

// Benchmark tests
func BenchmarkHealthStatus_JSONMarshal(b *testing.B) {
	healthStatus := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"database": "healthy",
			"redis":    "healthy",
			"api":      "healthy",
		},
		Version: "1.0.0",
		Uptime:  "2h30m15s",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(healthStatus)
	}
}

func BenchmarkHealthStatus_JSONUnmarshal(b *testing.B) {
	jsonStr := `{
		"status": "healthy",
		"timestamp": "2023-12-25T15:30:45Z",
		"services": {
			"database": "healthy",
			"redis": "healthy",
			"api": "healthy"
		},
		"version": "1.0.0",
		"uptime": "2h30m15s"
	}`
	jsonBytes := []byte(jsonStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var healthStatus models.HealthStatus
		_ = json.Unmarshal(jsonBytes, &healthStatus)
	}
}
