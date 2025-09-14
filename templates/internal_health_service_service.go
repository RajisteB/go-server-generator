//go:generate mockgen -destination=../../mocks/mock_health_service.go -package=mocks {{.Module}}/internal/health/service HealthService

package health

import (
	"context"
	"time"

	"{{.Module}}/internal/health/models"
	"{{.Module}}/internal/shared/logger"

	"gorm.io/gorm"
)

const (
	pkgName = "health"
	layer   = "service"
)

type HealthService interface {
	GetHealth(ctx context.Context) (*models.HealthStatus, error)
}

type ServiceImpl struct {
	log *logger.Logger
	db  *gorm.DB
}

func NewService(logger *logger.Logger, db *gorm.DB) HealthService {
	serviceLogger := logger.With("package", pkgName, "layer", layer)
	return &ServiceImpl{log: serviceLogger, db: db}
}

func (s *ServiceImpl) GetHealth(ctx context.Context) (*models.HealthStatus, error) {
	l := s.log.WithContext(ctx).With("operation", "GetHealth")

	// Check database connectivity
	sqlDB, err := s.db.DB()
	if err != nil {
		l.Error("failed to get sql db", "error", err)
		return &models.HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Services: map[string]string{
				"database": "unhealthy",
			},
		}, err
	}

	if err := sqlDB.Ping(); err != nil {
		l.Error("database ping failed", "error", err)
		return &models.HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Services: map[string]string{
				"database": "unhealthy",
			},
		}, err
	}

	return &models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"database": "healthy",
		},
		Version: "1.0.0",
		Uptime:  "unknown", // You can implement uptime tracking
	}, nil
}
