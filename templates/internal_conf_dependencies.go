package conf

import (
	healthController "{{.Module}}/internal/health/controller"
	healthService "{{.Module}}/internal/health/service"
	organizationsController "{{.Module}}/internal/organizations/controller"
	organizationsDatasource "{{.Module}}/internal/organizations/datasource"
	organizationsService "{{.Module}}/internal/organizations/service"
	"{{.Module}}/internal/shared/logger"
	"{{.Module}}/internal/shared/middleware"
	usersController "{{.Module}}/internal/users/controller"
	usersDatasource "{{.Module}}/internal/users/datasource"
	usersService "{{.Module}}/internal/users/service"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"gorm.io/gorm"
)

type ExternalDependencies struct {
	Clerk clerk.Client
}

type Controllers struct {
	Users         usersController.UsersController
	Organizations organizationsController.OrganizationsController
	Health        healthController.HealthController
}

type Dependencies struct {
	Config               *ConfigVars
	ExternalDependencies ExternalDependencies
	Controllers          Controllers
	Middleware           *middleware.Middleware
}

func LoadDependencies(logger *logger.Logger, config *ConfigVars, db *gorm.DB) *Dependencies {
	// Initialize Clerk client
	clerkClient, _ := clerk.NewClient(config.Clerk.Key)

	// Initialize datasources
	usersDS := usersDatasource.NewDatasource(logger, db)
	organizationsDS := organizationsDatasource.NewDatasource(logger, db)

	// Initialize services
	usersSvc := usersService.NewService(logger, usersDS)
	organizationsSvc := organizationsService.NewService(logger, organizationsDS)
	healthSvc := healthService.NewService(logger, db)

	// Initialize controllers
	usersCtrl := usersController.NewController(logger, usersSvc)
	organizationsCtrl := organizationsController.NewController(logger, organizationsSvc)
	healthCtrl := healthController.NewController(logger, healthSvc)

	// Initialize middleware
	mw := middleware.NewMiddleware(clerkClient, config.Clerk.Secret)

	return &Dependencies{
		Config: config,
		ExternalDependencies: ExternalDependencies{
			Clerk: clerkClient,
		},
		Controllers: Controllers{
			Users:         usersCtrl,
			Organizations: organizationsCtrl,
			Health:        healthCtrl,
		},
		Middleware: mw,
	}
}
