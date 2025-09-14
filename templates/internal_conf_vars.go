package conf

import (
	"fmt"
	"os"
	"strconv"

	"{{.Module}}/internal/shared/validation"
)

type ClerkVars struct {
	Key    string `validate:"required"`
	Secret string `validate:"required"`
}

type ServerVars struct {
	Name        string `validate:"required"`
	Version     string `validate:"required"`
	Environment string `validate:"required"`
	Host        string `validate:"required"`
	Port        string `validate:"required"`
	Protocol    string `validate:"required"`
}

type DatabaseVars struct {
	DatabaseHost     string `validate:"required"`
	DatabasePort     string `validate:"required"`
	DatabaseName     string `validate:"required"`
	DatabaseUser     string `validate:"required"`
	DatabasePassword string `validate:"required"`
	DatabaseSSLMode  string `validate:"required"`
}

type CSRFVars struct {
	AuthKey string `env:"CSRF_AUTH_KEY"`
	Secure  bool   `env:"CSRF_SECURE" default:"false"`
}

type SecurityVars struct {
	CSPPolicy          string `env:"SECURITY_CSP_POLICY"`
	HSTSMaxAge         int    `env:"SECURITY_HSTS_MAX_AGE" default:"31536000"`
	FrameOptions       string `env:"SECURITY_FRAME_OPTIONS" default:"DENY"`
	ContentTypeOptions bool   `env:"SECURITY_CONTENT_TYPE_OPTIONS" default:"true"`
	ReferrerPolicy     string `env:"SECURITY_REFERRER_POLICY" default:"strict-origin-when-cross-origin"`
	PermissionsPolicy  string `env:"SECURITY_PERMISSIONS_POLICY"`
}

type RequestLimitsVars struct {
	MaxRequestSize    int64 `env:"REQUEST_MAX_SIZE" default:"10485760"`              // 10MB default
	MaxHeaderSize     int64 `env:"REQUEST_MAX_HEADER_SIZE" default:"1048576"`        // 1MB default
	MaxFileUploadSize int64 `env:"REQUEST_MAX_FILE_UPLOAD_SIZE" default:"104857600"` // 100MB for file uploads
	ReadTimeout       int   `env:"REQUEST_READ_TIMEOUT" default:"30"`                // 30 seconds
	WriteTimeout      int   `env:"REQUEST_WRITE_TIMEOUT" default:"30"`               // 30 seconds
}

type ConfigVars struct {
	Clerk         ClerkVars
	Server        ServerVars
	Database      DatabaseVars
	CSRF          CSRFVars
	Security      SecurityVars
	RequestLimits RequestLimitsVars
}

// LoadConfigVarsFromEnv loads and validates all application configuration variables from environment variables.
func LoadConfigVarsFromEnv() (*ConfigVars, error) {
	dbVars := DatabaseVars{
		DatabaseHost:     os.Getenv("{{.Name | upper}}_DATABASE_HOST"),
		DatabasePort:     os.Getenv("{{.Name | upper}}_DATABASE_PORT"),
		DatabaseUser:     os.Getenv("{{.Name | upper}}_DATABASE_USER"),
		DatabaseName:     os.Getenv("{{.Name | upper}}_DATABASE_NAME"),
		DatabasePassword: os.Getenv("{{.Name | upper}}_DATABASE_PASSWORD"),
		DatabaseSSLMode:  os.Getenv("{{.Name | upper}}_DATABASE_SSL_MODE"),
	}

	serverVars := ServerVars{
		Environment: os.Getenv("{{.Name | upper}}_SERVER_ENV"),
		Version:     os.Getenv("{{.Name | upper}}_SERVER_VERSION"),
		Name:        os.Getenv("{{.Name | upper}}_SERVER_NAME"),
		Host:        os.Getenv("{{.Name | upper}}_SERVER_HOST"),
		Port:        os.Getenv("{{.Name | upper}}_SERVER_PORT"),
		Protocol:    os.Getenv("{{.Name | upper}}_SERVER_PROTOCOL"),
	}

	csrfVars := CSRFVars{
		AuthKey: os.Getenv("{{.Name | upper}}_CSRF_AUTH_KEY"),
		Secure:  os.Getenv("{{.Name | upper}}_CSRF_SECURE") == "false",
	}

	clerkVars := ClerkVars{
		Key:    os.Getenv("{{.Name | upper}}_CLERK_KEY"),
		Secret: os.Getenv("{{.Name | upper}}_CLERK_SECRET"),
	}

	securityVars := SecurityVars{
		CSPPolicy: os.Getenv("{{.Name | upper}}_SECURITY_CSP_POLICY"),
		HSTSMaxAge: func() int {
			if val := os.Getenv("{{.Name | upper}}_SECURITY_HSTS_MAX_AGE"); val != "" {
				if parsed, err := strconv.Atoi(val); err == nil {
					return parsed
				}
			}
			return 31536000 // Default 1 year
		}(),
		FrameOptions:       os.Getenv("{{.Name | upper}}_SECURITY_FRAME_OPTIONS"),
		ContentTypeOptions: os.Getenv("{{.Name | upper}}_SECURITY_CONTENT_TYPE_OPTIONS") != "false",
		ReferrerPolicy:     os.Getenv("{{.Name | upper}}_SECURITY_REFERRER_POLICY"),
		PermissionsPolicy:  os.Getenv("{{.Name | upper}}_SECURITY_PERMISSIONS_POLICY"),
	}

	requestLimitsVars := RequestLimitsVars{
		MaxRequestSize: func() int64 {
			if val := os.Getenv("{{.Name | upper}}_REQUEST_MAX_SIZE"); val != "" {
				if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
					return parsed
				}
			}
			return 10485760 // 10MB default
		}(),
		MaxHeaderSize: func() int64 {
			if val := os.Getenv("{{.Name | upper}}_REQUEST_MAX_HEADER_SIZE"); val != "" {
				if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
					return parsed
				}
			}
			return 1048576 // 1MB default
		}(),
		MaxFileUploadSize: func() int64 {
			if val := os.Getenv("{{.Name | upper}}_REQUEST_MAX_FILE_UPLOAD_SIZE"); val != "" {
				if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
					return parsed
				}
			}
			return 104857600 // 100MB default
		}(),
		ReadTimeout: func() int {
			if val := os.Getenv("{{.Name | upper}}_REQUEST_READ_TIMEOUT"); val != "" {
				if parsed, err := strconv.Atoi(val); err == nil {
					return parsed
				}
			}
			return 30 // 30 seconds default
		}(),
		WriteTimeout: func() int {
			if val := os.Getenv("{{.Name | upper}}_REQUEST_WRITE_TIMEOUT"); val != "" {
				if parsed, err := strconv.Atoi(val); err == nil {
					return parsed
				}
			}
			return 30 // 30 seconds default
		}(),
	}

	if err := validation.ValidateStruct(clerkVars); err != nil {
		return nil, fmt.Errorf("invalid clerk vars: %w", err)
	}

	if err := validation.ValidateStruct(serverVars); err != nil {
		return nil, fmt.Errorf("invalid server vars: %w", err)
	}

	if err := validation.ValidateStruct(dbVars); err != nil {
		return nil, fmt.Errorf("invalid db vars: %w", err)
	}

	if err := validation.ValidateStruct(csrfVars); err != nil {
		return nil, fmt.Errorf("invalid csrf vars: %w", err)
	}

	if err := validation.ValidateStruct(securityVars); err != nil {
		return nil, fmt.Errorf("invalid security vars: %w", err)
	}

	if err := validation.ValidateStruct(requestLimitsVars); err != nil {
		return nil, fmt.Errorf("invalid request limits vars: %w", err)
	}

	return &ConfigVars{
		Clerk:         clerkVars,
		Server:        serverVars,
		Database:      dbVars,
		CSRF:          csrfVars,
		Security:      securityVars,
		RequestLimits: requestLimitsVars,
	}, nil
}

// Sanitize methods for all structs
func (cv *ClerkVars) Sanitize() {
	cv.Key = validation.SanitizeString(cv.Key)
	cv.Secret = validation.SanitizeString(cv.Secret)
}

func (sv *ServerVars) Sanitize() {
	sv.Name = validation.SanitizeString(sv.Name)
	sv.Version = validation.SanitizeString(sv.Version)
	sv.Environment = validation.SanitizeString(sv.Environment)
	sv.Host = validation.SanitizeString(sv.Host)
	sv.Port = validation.SanitizeString(sv.Port)
	sv.Protocol = validation.SanitizeString(sv.Protocol)
}

func (dv *DatabaseVars) Sanitize() {
	dv.DatabaseHost = validation.SanitizeString(dv.DatabaseHost)
	dv.DatabasePort = validation.SanitizeString(dv.DatabasePort)
	dv.DatabaseName = validation.SanitizeString(dv.DatabaseName)
	dv.DatabaseUser = validation.SanitizeString(dv.DatabaseUser)
	dv.DatabasePassword = validation.SanitizeString(dv.DatabasePassword)
	dv.DatabaseSSLMode = validation.SanitizeString(dv.DatabaseSSLMode)
}

func (cv *CSRFVars) Sanitize() {
	cv.AuthKey = validation.SanitizeString(cv.AuthKey)
}

func (sv *SecurityVars) Sanitize() {
	sv.CSPPolicy = validation.SanitizeString(sv.CSPPolicy)
	sv.FrameOptions = validation.SanitizeString(sv.FrameOptions)
	sv.ReferrerPolicy = validation.SanitizeString(sv.ReferrerPolicy)
	sv.PermissionsPolicy = validation.SanitizeString(sv.PermissionsPolicy)
}

func (rlv *RequestLimitsVars) Sanitize() {
	// Ensure reasonable limits
	if rlv.MaxRequestSize <= 0 {
		rlv.MaxRequestSize = 10485760 // 10MB
	}
	if rlv.MaxHeaderSize <= 0 {
		rlv.MaxHeaderSize = 1048576 // 1MB
	}
	if rlv.MaxFileUploadSize <= 0 {
		rlv.MaxFileUploadSize = 104857600 // 100MB
	}
	if rlv.ReadTimeout <= 0 {
		rlv.ReadTimeout = 30
	}
	if rlv.WriteTimeout <= 0 {
		rlv.WriteTimeout = 30
	}
}

func (cv *ConfigVars) Sanitize() {
	cv.Clerk.Sanitize()
	cv.Server.Sanitize()
	cv.Database.Sanitize()
	cv.CSRF.Sanitize()
	cv.Security.Sanitize()
	cv.RequestLimits.Sanitize()
}
