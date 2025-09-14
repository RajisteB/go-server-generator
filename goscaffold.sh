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
    local default_dir="${NEW_GO_SERVER_DEFAULT_DIR:-$HOME/Projects}"
    
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
    
    # Store the original working directory
    local original_dir="$(pwd)"
    
    # Change to the scaffold directory
    cd "$scaffold_dir" || {
        echo "❌ Error: Could not change to $scaffold_dir"
        return 1
    }
    
    # Build command arguments
    local go_args=()
    
    if [[ -n "$project_name" ]]; then
        go_args+=("-name" "$project_name")
    fi
    
    if [[ -n "$module_name" ]]; then
        go_args+=("-module" "$module_name")
    fi
    
    if [[ -n "$description" ]]; then
        go_args+=("-description" "$description")
    fi
    
    if [[ -n "$port" ]]; then
        go_args+=("-port" "$port")
    fi
    
    if [[ -n "$project_path" ]]; then
        go_args+=("-path" "$project_path")
    fi
    
    # Run the generator with arguments
    echo "🚀 Starting Go Backend Project Generator..."
    go run main.go "${go_args[@]}"
    
    # Change to the created project directory
    local final_project_path=""
    if [[ -n "$project_path" ]]; then
        # If absolute path, use it directly
        if [[ "$project_path" == /* ]]; then
            final_project_path="$project_path"
        else
            # If relative path, make it relative to original directory
            final_project_path="$original_dir/$project_path"
        fi
    else
        # Default: create project in default directory with project name
        final_project_path="$default_dir/$project_name"
    fi
    
    # Convert absolute path to ~/ format if it's under home directory
    local home_dir="$HOME"
    if [[ "$final_project_path" == "$home_dir"* ]]; then
        final_project_path="~${final_project_path#$home_dir}"
    fi
    
    # Always try to change directory - the Go program will have created the project
    # We'll determine the project name from the "Next steps" output or use a fallback
    echo ""
    echo "📁 Attempting to change to project directory..."
    
    # Try to extract project name from the last output or use the provided name
    local actual_project_name="$project_name"
    if [[ -z "$actual_project_name" ]]; then
        # If no project name was provided via args, try to find the most recent directory
        # that was created in the default directory
        actual_project_name=$(find "$default_dir" -maxdepth 1 -type d -newer "$default_dir" 2>/dev/null | head -1 | xargs basename)
    fi
    
    if [[ -n "$actual_project_name" ]]; then
        # Determine the correct path based on whether a custom path was provided
        local expanded_path=""
        if [[ -n "$project_path" ]]; then
            # Use the custom path
            if [[ "$project_path" == /* ]]; then
                expanded_path="$project_path"
            else
                expanded_path="$original_dir/$project_path"
            fi
        else
            # Use the default directory
            expanded_path="$default_dir/$actual_project_name"
        fi
        
        echo "📁 Changing to project directory: ~${expanded_path#$HOME}"
        cd "$expanded_path" || {
            echo "⚠️  Warning: Could not change to project directory: ~${expanded_path#$HOME}"
            echo "   You may need to manually navigate to your project directory"
        }
    else
        echo "⚠️  Warning: Could not determine project directory to change into"
        echo "   Please manually navigate to your project directory"
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
