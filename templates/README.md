# {{.Name}}

{{.Description}}

## Features

- ✅ **Clerk Authentication** - Modern authentication with Clerk
- ✅ **PostgreSQL Database** - Robust database with GORM ORM
- ✅ **Clean Architecture** - Controller-Service-Datasource pattern
- ✅ **Structured Logging** - JSON logging with context
- ✅ **Security Middleware** - CSRF, CORS, Rate limiting, Security headers
- ✅ **Request Validation** - Input validation and sanitization
- ✅ **Health Checks** - Application health monitoring
- ✅ **Graceful Shutdown** - Proper server shutdown handling
- ✅ **Configuration Management** - Viper-based config with .env.local support
- ✅ **UUID Generation** - Multiple UUID formats (standard, short, namespaced)
- ✅ **CI/CD Pipeline** - GitHub Actions workflow with tests, coverage, and builds
- ✅ **Development Tools** - Makefile with comprehensive development commands

## Getting Started

### Quick Start

1. **Setup development environment:**
   ```bash
   make setup
   ```

2. **Update configuration:**
   Edit `.env.local` with your configuration (database, Clerk keys, etc.)

3. **Run the application:**
   ```bash
   make dev
   ```
   
   The server will start on port {{.Port}} (configurable via {{.Name | upper}}_SERVER_PORT environment variable)

### Manual Setup

1. **Install dependencies:**
   ```bash
   make install-deps
   ```

2. **Install development tools:**
   ```bash
   make install-tools
   ```

3. **Run tests:**
   ```bash
   make test
   ```

4. **Build application:**
   ```bash
   make build
   ```

## Docker

Build and run with Docker:

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Or manually:
docker build -t {{.Name}} .
docker run -p {{.Port}}:{{.Port}} {{.Name}}
```

## Project Structure

```
.
├── .github/               # GitHub Actions CI/CD
│   └── workflows/
│       └── ci.yml         # CI pipeline
├── cmd/                   # Application entrypoints
├── internal/              # Private application code
│   ├── conf/              # Configuration management
│   ├── handlers/          # HTTP request handlers
│   ├── health/            # Health check endpoints
│   ├── organizations/     # Organization domain
│   │   ├── controller/    # HTTP controllers
│   │   ├── datasource/    # Data access layer
│   │   ├── models/        # Data models
│   │   └── service/       # Business logic
│   ├── shared/            # Shared utilities
│   │   ├── assertions/    # Validation assertions
│   │   ├── constants/     # Application constants
│   │   ├── http/          # HTTP helpers
│   │   ├── logger/        # Structured logging
│   │   ├── middleware/    # HTTP middleware
│   │   ├── uuid/          # UUID generation
│   │   └── validation/    # Input validation
│   ├── tests/             # Test files
│   └── users/             # User domain
│       ├── controller/    # HTTP controllers
│       ├── datasource/    # Data access layer
│       ├── models/        # Data models
│       └── service/       # Business logic
├── .env.local            # Environment variables template
├── .gitignore            # Git ignore rules
├── Makefile              # Development commands
└── README.md             # This file
```

## API Endpoints

### Health
- `GET /api/v1/health` - Health check

### Authentication
- `GET /api/v1/csrf-token` - Get CSRF token

### Users (Webhook endpoints)
- `POST /api/v1/users/clerk` - Create user from Clerk webhook
- `PUT /api/v1/users/clerk` - Update user from Clerk webhook
- `DELETE /api/v1/users/clerk` - Delete user from Clerk webhook

### Organizations (Protected endpoints)
- `GET /api/v1/organizations/{id}` - Get organization by ID
- `GET /api/v1/organizations/clerk/{clerk_id}` - Get organization by Clerk ID

## Environment Variables

See `.env.local` for all available configuration options. The application uses Viper for configuration management, which automatically loads from `.env.local` files with fallback to system environment variables.

### Key Configuration Areas:
- **Server**: Host, port, protocol, environment
- **Database**: PostgreSQL connection settings (default password: `root`)
- **Clerk**: Authentication keys and configuration
- **Security**: CSRF, security headers, request limits

## Development

The project follows clean architecture principles with clear separation of concerns:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Datasources**: Handle data persistence
- **Models**: Define data structures
- **Shared**: Reusable utilities and middleware

### Available Commands

```bash
make help              # Show all available commands
make setup             # Setup development environment
make dev               # Start development server
make dev-watch         # Start with file watching (requires air)
make test              # Run unit tests
make coverage          # Generate coverage report
make build             # Build application
make lint              # Run linter (requires golangci-lint)
make static-analysis   # Run static analysis tools
make check-rules       # Run NASA rule checks
make clean             # Clean build artifacts
```

### CI/CD Pipeline

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that runs:
- **Tests**: Unit tests with coverage
- **Safety Checks**: NASA rule compliance
- **Coverage**: Test coverage reporting
- **Build**: Application compilation
- **Lint**: Code quality checks (commented out, ready to enable)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
