#!/bin/bash

# Go Server Generator Function
# Compatible with both bash and zsh
# Add this function to your ~/.bashrc or ~/.zshrc file
# Usage: go-server --create [options]

go-server() {
    local scaffold_dir="$HOME/Projects/go-scaffold"
    local project_name=""
    local module_name=""
    local description=""
    local port="8080"
    local project_path=""
    local interactive_mode=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --create)
                interactive_mode=false
                shift
                ;;
            --name|-n)
                project_name="$2"
                shift 2
                ;;
            --module|-m)
                module_name="$2"
                shift 2
                ;;
            --description|-d)
                description="$2"
                shift 2
                ;;
            --port|-p)
                port="$2"
                shift 2
                ;;
            --path)
                project_path="$2"
                shift 2
                ;;
            --help|-h)
                go-server-help
                return 0
                ;;
            *)
                # If no flag, treat as project name for backward compatibility
                if [[ -z "$project_name" ]]; then
                    project_name="$1"
                fi
                shift
                ;;
        esac
    done
    
    # Check if the scaffold directory exists
    if [ ! -d "$scaffold_dir" ]; then
        echo "❌ Error: Go scaffold directory not found at $scaffold_dir"
        echo "Please ensure the go-scaffold project is located at ~/Projects/go-scaffold"
        return 1
    fi
    
    # Check if main.go exists in the scaffold directory
    if [ ! -f "$scaffold_dir/main.go" ]; then
        echo "❌ Error: main.go not found in $scaffold_dir"
        echo "Please ensure the go-scaffold project is properly set up"
        return 1
    fi
    
    # Change to the scaffold directory
    cd "$scaffold_dir" || {
        echo "❌ Error: Could not change to $scaffold_dir"
        return 1
    }
    
    # If interactive mode or missing required parameters, run interactively
    if [[ "$interactive_mode" == true ]] || [[ -z "$project_name" ]] || [[ -z "$module_name" ]]; then
        echo "🚀 Starting Go Backend Project Generator..."
        go run main.go
    else
        # Non-interactive mode with provided parameters
        echo "🚀 Creating Go Backend Project: $project_name"
        
        # Create a temporary input file for non-interactive execution
        local temp_input=$(mktemp)
        {
            echo "$project_name"
            echo "$module_name"
            echo "$description"
            echo "$port"
            echo "$project_path"
        } > "$temp_input"
        
        # Run the generator with the input file
        go run main.go < "$temp_input"
        
        # Clean up
        rm "$temp_input"
    fi
}

# Help function
go-server-help() {
    echo "Go Server Generator Commands:"
    echo "============================"
    echo ""
    echo "go-server --create [options]  - Create a new Go backend project"
    echo "go-server --help              - Show this help message"
    echo ""
    echo "Options:"
    echo "  --name, -n <name>           - Project name"
    echo "  --module, -m <module>       - Go module name (e.g., github.com/username/project)"
    echo "  --description, -d <desc>    - Project description"
    echo "  --port, -p <port>           - Server port (default: 8080)"
    echo "  --path <path>               - Project path (default: current directory)"
    echo "  --help, -h                  - Show this help message"
    echo ""
    echo "Examples:"
    echo "  go-server --create --name my-api --module github.com/user/my-api --port 3000"
    echo "  go-server --create -n my-app -m github.com/user/my-app -d \"My awesome API\""
    echo "  go-server --create  # Interactive mode"
    echo ""
    echo "Interactive mode will prompt you for all required information."
}
