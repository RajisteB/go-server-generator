package health

import (
	"net/http"

	healthService "{{.Module}}/internal/health/service"
	httpHelpers "{{.Module}}/internal/shared/http"
	"{{.Module}}/internal/shared/logger"
)

const (
	pkgName = "health"
	layer   = "controller"
)

type HealthController interface {
	GetHealth(w http.ResponseWriter, r *http.Request) error
}

type ControllerImpl struct {
	log     *logger.Logger
	service healthService.HealthService
}

func NewController(logger *logger.Logger, service healthService.HealthService) HealthController {
	ctrlLogger := logger.With("package", pkgName, "layer", layer)
	return &ControllerImpl{log: ctrlLogger, service: service}
}

func (c *ControllerImpl) GetHealth(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	l := c.log.WithContext(ctx).With("method", "GetHealth")

	health, err := c.service.GetHealth(ctx)
	if err != nil {
		l.Error("failed to get health status", "error", err)
		return httpHelpers.RespondWithError(w, err)
	}

	return httpHelpers.RespondWithJSON(w, http.StatusOK, health)
}
