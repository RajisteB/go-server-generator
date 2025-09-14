# Go Backend Project Generator

A CLI tool that generates complete Go backend projects with a clean, scalable architecture.

## 🚀 Quick Start

### Installation

Add the generator to your shell configuration:

```bash
# For bash
echo 'source ~/Projects/go-scaffold/goscaffold.sh' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'source ~/Projects/go-scaffold/goscaffold.sh' >> ~/.zshrc
source ~/.zshrc
```

### Usage

```bash
# Interactive mode
go-server --create

# With flags
go-server --create --name my-api --module github.com/user/my-api --port 3000

# Show help
go-server --help
```

## 📁 Generated Project Structure

```
your-project/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── cmd/
│   └── root.go               # CLI root command
├── internal/
│   ├── conf/                 # Configuration management
│   ├── shared/               # Shared utilities
│   │   ├── logger/           # Logging utilities
│   │   ├── validation/       # Input validation
│   │   ├── environment/      # Environment variables
│   │   ├── constants/        # Application constants
│   │   ├── http/             # HTTP utilities
│   │   ├── uuid/             # UUID generation
│   │   ├── assertions/       # Assertion helpers
│   │   └── middleware/       # HTTP middleware
│   ├── handlers/             # HTTP handlers
│   ├── health/               # Health check module
│   ├── users/                # User management module
│   └── organizations/        # Organization management module
└── .env.example             # Environment variables template
```

## 🛠️ Features

- **Clean Architecture**: Follows DDD pattern with clear separation of concerns
- **Modular Design**: Organized into logical modules (users, organizations, health)
- **Configuration Management**: Centralized config with environment variable support
- **Database Integration**: PostgreSQL integration with connection management
- **HTTP Utilities**: Middleware, validation, and response helpers
- **Logging**: Structured logging with configurable levels
- **Health Checks**: Built-in health monitoring endpoints

## 📋 Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--create` | - | Create a new project | Required |
| `--name` | `-n` | Project name | - |
| `--module` | `-m` | Go module name | - |
| `--description` | `-d` | Project description | - |
| `--port` | `-p` | Server port | 8080 |
| `--path` | - | Project path | Current directory |
| `--help` | `-h` | Show help | - |

## 🔧 After Generation

1. Navigate to your project directory
2. Copy environment template: `cp .env.example .env`
3. Update `.env` with your configuration
4. Install dependencies: `go mod tidy`
5. Run the server: `go run main.go`

## 🏗️ Architecture

The generator creates projects following these principles:

- **Domain-Driven Design**: Clear module boundaries
- **Dependency Injection**: Loose coupling between components
- **Configuration as Code**: Centralized configuration management
- **Middleware Pattern**: Reusable HTTP middleware
- **Service Layer**: Business logic separation

## 📝 License

MIT License - feel free to use and modify as needed.
