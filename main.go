package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type ProjectConfig struct {
	Name        string
	Module      string
	Description string
	Port        string
	ProjectPath string
}

// TemplateFile represents a template file mapping
type TemplateFile struct {
	SourcePath string // Path in templates directory
	TargetPath string // Path in generated project
}

// Define all template files based on clean architecture
var templateFiles = []TemplateFile{
	// Core application files
	{"main.go", "main.go"},
	{"go.mod", "go.mod"},
	{"cmd_root.go", "cmd/root.go"},
	{"env_local", ".env.local"},
	{"gitignore", ".gitignore"},
	{"github_yml", ".github/workflows/ci.yml"},
	{"README.md", "README.md"},
	{"Makefile", "Makefile"},

	// Configuration
	{"internal_conf_vars.go", "internal/conf/vars.go"},
	{"internal_conf_pg.go", "internal/conf/pg.go"},
	{"internal_conf_dependencies.go", "internal/conf/dependencies.go"},

	// Shared utilities
	{"internal_shared_logger_logger.go", "internal/shared/logger/logger.go"},
	{"internal_shared_validation_validation.go", "internal/shared/validation/validation.go"},
	{"internal_shared_constants_constants.go", "internal/shared/constants/constants.go"},
	{"internal_shared_http_http.go", "internal/shared/http/http.go"},
	{"internal_shared_uuid_uuid.go", "internal/shared/uuid/uuid.go"},
	{"internal_shared_assertions_assertions.go", "internal/shared/assertions/assertions.go"},
	{"internal_shared_middleware_middleware.go", "internal/shared/middleware/middleware.go"},

	// Shared utilities tests
	{"internal_tests_shared_assertions_assertions_test.go", "internal/tests/shared/assertions/assertions_test.go"},
	{"internal_tests_shared_validation_validation_test.go", "internal/tests/shared/validation/validation_test.go"},
	{"internal_tests_shared_logger_logger_test.go", "internal/tests/shared/logger/logger_test.go"},
	{"internal_tests_shared_constants_constants_test.go", "internal/tests/shared/constants/constants_test.go"},
	{"internal_tests_shared_http_http_test.go", "internal/tests/shared/http/http_test.go"},
	{"internal_tests_shared_uuid_uuid_test.go", "internal/tests/shared/uuid/uuid_test.go"},
	{"internal_tests_shared_middleware_middleware_test.go", "internal/tests/shared/middleware/middleware_test.go"},

	// Handlers
	{"internal_handlers_handlers.go", "internal/handlers/handlers.go"},

	// Health module
	{"internal_health_models_health.go", "internal/health/models/health.go"},
	{"internal_health_service_service.go", "internal/health/service/service.go"},
	{"internal_health_controller_controller.go", "internal/health/controller/controller.go"},

	// Users module
	{"internal_users_models_users.go", "internal/users/models/users.go"},
	{"internal_users_datasource_datasource.go", "internal/users/datasource/datasource.go"},
	{"internal_users_service_service.go", "internal/users/service/service.go"},
	{"internal_users_controller_controller.go", "internal/users/controller/controller.go"},

	// Organizations module
	{"internal_organizations_models_organizations.go", "internal/organizations/models/organizations.go"},
	{"internal_organizations_datasource_datasource.go", "internal/organizations/datasource/datasource.go"},
	{"internal_organizations_service_service.go", "internal/organizations/service/service.go"},
	{"internal_organizations_controller_controller.go", "internal/organizations/controller/controller.go"},
}

func main() {
	fmt.Println("🚀 Go Backend Project Generator")
	fmt.Println("================================================")

	config := ProjectConfig{}

	// Parse command line flags
	var (
		name        = flag.String("name", "", "Project name")
		module      = flag.String("module", "", "Go module name (e.g., github.com/username/project)")
		description = flag.String("description", "", "Project description")
		port        = flag.String("port", "8080", "Server port")
		projectPath = flag.String("path", "", "Project path (leave empty to create in current directory)")
	)
	flag.Parse()

	// Use command line arguments if provided, otherwise prompt for input
	if *name != "" {
		config.Name = *name
	} else {
		config.Name = getUserInput("Project name: ")
	}

	if *module != "" {
		config.Module = *module
	} else {
		config.Module = getUserInput("Go module name (e.g., github.com/username/project): ")
	}

	// Auto-convert simple names to project-name/module format
	if !strings.Contains(config.Module, "/") {
		config.Module = config.Name + "/" + config.Module
	}


	if *description != "" {
		config.Description = *description
	} else {
		config.Description = getUserInput("Project description: ")
	}

	if *port != "" {
		config.Port = *port
	} else {
		config.Port = getUserInput("Server port [8080]: ")
	}

	if *projectPath != "" {
		config.ProjectPath = *projectPath
	} else {
		config.ProjectPath = getUserInput("Project path (leave empty to create in current directory, or specify absolute/relative path): ")
	}

	// Set default port if empty
	if config.Port == "" {
		config.Port = "8080"
	}

	// Set default project path if empty
	if config.ProjectPath == "" {
		// Use environment variable or default to current directory
		if defaultDir := os.Getenv("NEW_GO_SERVER_DEFAULT_DIR"); defaultDir != "" {
			config.ProjectPath = filepath.Join(defaultDir, config.Name)
		} else {
			config.ProjectPath = config.Name
		}
	}

	// Validate and expand the project path
	if err := validateProjectPath(config.ProjectPath); err != nil {
		fmt.Printf("Error with project path: %v\n", err)
		os.Exit(1)
	}

	// Create project
	fmt.Printf("\nCreating project '%s'...\n", config.Name)
	if err := createProject(config); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Project '%s' created successfully!\n", config.Name)
	fmt.Println("\nNext steps:")

	// Determine the correct directory to show in the cd command
	var cdPath string
	if config.ProjectPath == "." {
		cdPath = config.Name
	} else {
		cdPath = config.ProjectPath
	}

	// Convert to ~/ format if it's under home directory
	homeDir := os.Getenv("HOME")
	if strings.HasPrefix(cdPath, homeDir) {
		cdPath = "~" + strings.TrimPrefix(cdPath, homeDir)
	}

	fmt.Printf("  cd %s\n", cdPath)
	fmt.Println("  cp .env.local .env.local")
	fmt.Println("  # Update .env.local with your configuration")
	fmt.Println("  go mod tidy")
	fmt.Println("  go run main.go")
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}


func validateProjectPath(path string) error {
	if path == "." {
		return nil // Current directory is always valid
	}

	// Check if the path is absolute or relative
	if filepath.IsAbs(path) {
		// For absolute paths, check if parent directory exists
		parentDir := filepath.Dir(path)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			return fmt.Errorf("parent directory does not exist: %s", parentDir)
		}
	} else {
		// For relative paths, check if the path is valid
		if _, err := filepath.Abs(path); err != nil {
			return fmt.Errorf("invalid path: %s", path)
		}
	}

	return nil
}

func createProject(config ProjectConfig) error {
	// Determine the full project path
	var projectDir string
	if config.ProjectPath == "." {
		// If using current directory, create a subdirectory with project name
		projectDir = config.Name
	} else {
		// If custom path provided, use it directly
		projectDir = config.ProjectPath
	}

	// Create project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}

	// Create all template files
	for _, templateFile := range templateFiles {
		if err := createFileFromTemplateFile(config, templateFile, projectDir); err != nil {
			return err
		}
	}

	// Initialize go module
	if err := initializeGoModule(projectDir, config.Module); err != nil {
		return err
	}

	return nil
}

func createFileFromTemplateFile(config ProjectConfig, templateFile TemplateFile, projectDir string) error {
	// Read template content
	templatePath := filepath.Join("templates", templateFile.SourcePath)
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Create target file
	targetPath := filepath.Join(projectDir, templateFile.TargetPath)
	if err := createFileFromTemplate(targetPath, string(templateContent), config); err != nil {
		return err
	}

	return nil
}

func initializeGoModule(projectPath, moduleName string) error {
	// Change to project directory
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(projectPath); err != nil {
		return err
	}

	// The go.mod file is already created from template, so we don't need to initialize it
	// Just verify it exists
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("go.mod file not found after template creation")
	}

	// Download dependencies
	cmd := exec.Command("go", "mod", "tidy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download dependencies: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func createFileFromTemplate(filename, templateStr string, config ProjectConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create template with custom functions
	tmpl, err := template.New("file").Funcs(template.FuncMap{
		"upper": strings.ToUpper,
	}).Parse(templateStr)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}
