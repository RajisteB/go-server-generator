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

## Getting Started

1. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

2. Update the .env file with your configuration

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```
   
   The server will start on port {{.Port}} (configurable via {{.Name | upper}}_SERVER_PORT environment variable)

## Docker

Build and run with Docker:

```bash
docker build -t {{.Name}} .
docker run -p {{.Port}}:{{.Port}} {{.Name}}
```

## Project Structure

```
.
├── cmd/                    # Application entrypoints
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
│   │   ├── environment/   # Environment utilities
│   │   ├── http/          # HTTP helpers
│   │   ├── logger/        # Structured logging
│   │   ├── middleware/    # HTTP middleware
│   │   ├── uuid/          # UUID generation
│   │   └── validation/    # Input validation
│   └── users/             # User domain
│       ├── controller/    # HTTP controllers
│       ├── datasource/    # Data access layer
│       ├── models/        # Data models
│       └── service/       # Business logic
├── .env.example          # Environment variables template
└── README.md            # This file
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

See `.env.example` for all available configuration options.

## Development

The project follows clean architecture principles with clear separation of concerns:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Datasources**: Handle data persistence
- **Models**: Define data structures
- **Shared**: Reusable utilities and middleware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
