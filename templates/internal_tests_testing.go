package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{.Module}}/internal/mocks"
	"{{.Module}}/internal/shared/logger"

	"go.uber.org/mock/gomock"
)

// MockServer represents a test server instance
type MockServer struct {
	Server     *httptest.Server
	Logger     *logger.Logger
	Controller *gomock.Controller
}

// NewMockServer creates a new mock server for testing
func NewMockServer(t *testing.T) *MockServer {
	testLogger := logger.NewLogger()
	ctrl := gomock.NewController(t)

	// Create a simple test server
	router := http.NewServeMux()
	server := httptest.NewServer(router)

	return &MockServer{
		Server:     server,
		Logger:     testLogger,
		Controller: ctrl,
	}
}

// Close shuts down the mock server
func (m *MockServer) Close() {
	m.Server.Close()
	m.Controller.Finish()
}

// GetURL returns the base URL of the mock server
func (m *MockServer) GetURL() string {
	return m.Server.URL
}

// GetClient returns an HTTP client configured for the mock server
func (m *MockServer) GetClient() *http.Client {
	return m.Server.Client()
}

// CreateMockMocks creates mock instances for testing
func CreateMockMocks(ctrl *gomock.Controller) *MockMocks {
	return &MockMocks{
		UsersDatasource:         mocks.NewMockUsersDatasource(ctrl),
		UsersService:            mocks.NewMockUsersService(ctrl),
		OrganizationsDatasource: mocks.NewMockOrganizationsDatasource(ctrl),
		OrganizationsService:    mocks.NewMockOrganizationsService(ctrl),
		HealthService:           mocks.NewMockHealthService(ctrl),
		GormDB:                  mocks.NewMockGormDB(ctrl),
		SQLDB:                   mocks.NewMockSQLDB(ctrl),
	}
}

// MockMocks contains all mock instances
type MockMocks struct {
	UsersDatasource         *mocks.MockUsersDatasource
	UsersService            *mocks.MockUsersService
	OrganizationsDatasource *mocks.MockOrganizationsDatasource
	OrganizationsService    *mocks.MockOrganizationsService
	HealthService           *mocks.MockHealthService
	GormDB                  *mocks.MockGormDB
	SQLDB                   *mocks.MockSQLDB
}
