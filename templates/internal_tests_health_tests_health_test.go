package health_tests

import (
	"testing"
	"time"

	"{{.Module}}/internal/health/models"

	"github.com/stretchr/testify/assert"
)

func TestHealthModels(t *testing.T) {
	t.Run("HealthStatus_Structure", func(t *testing.T) {
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
	})

	t.Run("HealthStatus_UnhealthyStatus", func(t *testing.T) {
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
	})

	t.Run("HealthStatus_EmptyServices", func(t *testing.T) {
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
	})
}
