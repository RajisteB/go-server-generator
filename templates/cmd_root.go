package cmd

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"{{.Module}}/internal/conf"
	"{{.Module}}/internal/handlers"
	"{{.Module}}/internal/shared/constants"
	"{{.Module}}/internal/shared/logger"

	"github.com/rs/cors"
	"gorm.io/gorm"
)

type RootConfig struct {
	Logger *logger.Logger
	Config *conf.ConfigVars
	DB     *gorm.DB
}

func loadRootConfig() *RootConfig {
	vars, err := conf.LoadConfigVarsFromEnv()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewLogger(logger.DevelopmentConfig(vars.Server.Name, vars.Server.Version))

	return &RootConfig{
		Logger: appLogger,
		Config: vars,
	}
}

func (root *RootConfig) loadDatabase() *RootConfig {
	db, err := conf.InitConnectionPool(conf.PGConfig{
		Host:     root.Config.Database.DatabaseHost,
		Port:     root.Config.Database.DatabasePort,
		User:     root.Config.Database.DatabaseUser,
		DBName:   root.Config.Database.DatabaseName,
		Password: root.Config.Database.DatabasePassword,
		SSLMode:  root.Config.Database.DatabaseSSLMode,
		MaxConns: 10,
		Logger:   root.Logger,
	})
	if err != nil {
		root.Logger.Error("database connection failed", "error", err)
		panic(err)
	}

	root.DB = db
	root.Logger.Info("database connected")
	return root
}

func (root *RootConfig) exec() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	root.Logger.Info("application context created")

	root.loadDatabase()
	root.Logger.Info("database loaded")

	dependencies := conf.LoadDependencies(root.Logger, root.Config, root.DB)
	root.Logger.Info("dependencies loaded")

	handler := handlers.NewHandler(root.Logger, dependencies).Register()
	root.Logger.Info("handler registered")

	crossOrigin := cors.New(cors.Options{
		AllowedOrigins: func() []string {
			env := os.Getenv("{{.Name | upper}}_ENVIRONMENT")
			switch env {
			case "production":
				return []string{
					"https://yourdomain.com",
					"https://www.yourdomain.com",
				}
			case "staging":
				return []string{
					"https://staging.yourdomain.com",
					"https://staging-app.yourdomain.com",
				}
			default:
				return []string{
					constants.ServerAllowedOriginLocal,
					constants.ServerAllowedOriginVite,
					constants.ServerAllowedOriginReact,
					constants.ServerAllowedOriginReactNative,
					constants.ServerAllowedOriginPostman,
				}
			}
		}(),
		AllowedMethods: []string{
			string(constants.AllowedMethodGET),
			string(constants.AllowedMethodPOST),
			string(constants.AllowedMethodPUT),
			string(constants.AllowedMethodPATCH),
			string(constants.AllowedMethodDELETE),
		},
		AllowCredentials: true,
	})
	root.Logger.Info("cors middleware generated")

	appHandler := crossOrigin.Handler(handler)
	root.Logger.Info("app handler generated")

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", root.Config.Server.Port),
		Handler:        appHandler,
		WriteTimeout:   constants.WriteTimeout,
		ReadTimeout:    constants.ReadTimeout,
		IdleTimeout:    constants.IdleTimeout,
		MaxHeaderBytes: int(root.Config.RequestLimits.MaxHeaderSize),
	}

	root.Logger.Info("starting server",
		"port", root.Config.Server.Port,
		"maxRequestSize", root.Config.RequestLimits.MaxRequestSize,
		"readTimeout", root.Config.RequestLimits.ReadTimeout,
		"writeTimeout", root.Config.RequestLimits.WriteTimeout,
	)

	defer func() {
		root.Logger.Info("closing database connection")
		sqlDB, err := root.DB.DB()
		if err != nil {
			root.Logger.Error("failed to get sql db", "error", err)
		}
		if err := sqlDB.Close(); err != nil {
			root.Logger.Error("failed to close database connection", "error", err)
		}
		root.Logger.Info("database connection closed")
	}()

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		constants.ShutdownGracePeriod,
		"duration for which the server gracefully waits for existing connections to finish",
	)
	flag.Parse()

	// run server in goroutine to prevent blocking
	go func() {
		root.Logger.Info("{{.Name}} service running", "port", root.Config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			root.Logger.Error("unexpected server error", "error", err)
		}
	}()

	// Block until we receive our signal
	<-ctx.Done()
	root.Logger.Info("received shutdown signal, shutting down {{.Name}} service gracefully")

	// Create a deadline to wait for
	cx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := server.Shutdown(cx); err != nil {
		root.Logger.Error("error during server shutdown")
	}

	root.Logger.Info("application successfully shutdown")
	os.Exit(0)
}

func Run() {
	root := loadRootConfig()
	root.exec()
}
