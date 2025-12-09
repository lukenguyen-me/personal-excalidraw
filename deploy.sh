#!/bin/bash

# Personal Excalidraw Deployment Script
# Convenient commands for managing production deployment

set -e

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"
PROJECT_NAME="personal-excalidraw"

# Detect $DOCKER_COMPOSE command (supports both old and new syntax)
if $DOCKER_COMPOSE version >/dev/null 2>&1; then
    DOCKER_COMPOSE="$DOCKER_COMPOSE"
elif docker-compose version >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker-compose"
else
    echo "Error: Neither '$DOCKER_COMPOSE' nor 'docker-compose' is available"
    echo "Please install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_env_file() {
    if [ ! -f "$ENV_FILE" ]; then
        print_error "Environment file $ENV_FILE not found!"
        print_info "Run '$0 config' to create it interactively"
        print_info "Or copy .env.production.example to $ENV_FILE and configure it manually"
        exit 1
    fi
}

generate_random_key() {
    openssl rand -base64 32 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
}

# Command functions
config() {
    print_info "Interactive Configuration Setup"
    echo ""

    # Check if .env.production already exists
    if [ -f "$ENV_FILE" ]; then
        print_warning "Configuration file $ENV_FILE already exists!"
        read -p "Do you want to overwrite it? (yes/no): " overwrite
        if [ "$overwrite" != "yes" ]; then
            print_info "Configuration cancelled"
            exit 0
        fi
        # Backup existing file
        cp "$ENV_FILE" "$ENV_FILE.backup.$(date +%Y%m%d_%H%M%S)"
        print_info "Existing configuration backed up"
        echo ""
    fi

    # Check if example file exists
    if [ ! -f ".env.production.example" ]; then
        print_error ".env.production.example file not found!"
        exit 1
    fi

    echo -e "${GREEN}Required Settings${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Database Password
    print_info "Database Password"
    echo "This is the password for your PostgreSQL database."
    echo "Press Enter to generate a secure random password, or type your own:"
    read -p "> " db_password
    if [ -z "$db_password" ]; then
        db_password=$(generate_random_key)
        print_success "Generated random password: $db_password"
    fi
    echo ""

    # Access Key
    print_info "Access Key"
    echo "This is the password you'll use to log into the application."
    echo "Press Enter to generate a secure random key, or type your own:"
    read -p "> " access_key
    if [ -z "$access_key" ]; then
        access_key=$(generate_random_key)
        print_success "Generated random access key: $access_key"
    fi
    echo ""

    echo -e "${YELLOW}Optional Settings${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Application Port
    print_info "Application Port (default: 3000)"
    echo "The port where your application will be accessible via nginx."
    read -p "> " app_port
    app_port=${app_port:-3000}
    echo ""

    # CORS Origins
    print_info "CORS Allowed Origins (default: *)"
    echo "Set this to your domain for better security (e.g., https://yourdomain.com)"
    echo "Use * to allow all origins (less secure but easier for testing)"
    read -p "> " cors_origins
    cors_origins=${cors_origins:-*}
    echo ""

    # Log Level
    print_info "Log Level (default: info)"
    echo "Options: debug, info, warn, error"
    read -p "> " log_level
    log_level=${log_level:-info}
    echo ""

    # Enable Authentication
    print_info "Enable Authentication? (default: yes)"
    echo "Recommended: yes (requires access key to use the app)"
    read -p "> " auth_enabled
    if [ "$auth_enabled" = "no" ] || [ "$auth_enabled" = "n" ] || [ "$auth_enabled" = "false" ]; then
        auth_enabled="false"
    else
        auth_enabled="true"
    fi
    echo ""

    # Create the configuration file
    print_info "Creating configuration file..."
    cat > "$ENV_FILE" << EOF
# Production Environment Configuration
# Generated on $(date)

# Application Port (exposed to host)
APP_PORT=$app_port

# Database Configuration
DB_USER=postgres
DB_PASSWORD=$db_password
DB_NAME=personal_excalidraw
DB_SSLMODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# Server Configuration
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# CORS Configuration
CORS_ALLOWED_ORIGINS=$cors_origins

# Logging
LOG_LEVEL=$log_level

# Authentication
AUTH_ENABLED=$auth_enabled
ACCESS_KEY=$access_key
EOF

    print_success "Configuration file created: $ENV_FILE"
    echo ""

    echo -e "${GREEN}Summary${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "  Port:           ${BLUE}$app_port${NC}"
    echo -e "  Database:       ${BLUE}postgres (password set)${NC}"
    echo -e "  Access Key:     ${BLUE}$access_key${NC}"
    echo -e "  CORS:           ${BLUE}$cors_origins${NC}"
    echo -e "  Log Level:      ${BLUE}$log_level${NC}"
    echo -e "  Auth Enabled:   ${BLUE}$auth_enabled${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    print_success "Configuration complete!"
    print_info "Your access key is: ${GREEN}$access_key${NC}"
    print_warning "Save this access key! You'll need it to log in."
    echo ""
    print_info "Next step: Run '$0 start' to start the application"
}

start() {
    print_info "Starting Personal Excalidraw..."
    check_env_file

    # Export all variables from the env file so $DOCKER_COMPOSE can use them
    set -a
    source "$ENV_FILE"
    set +a

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d

    print_success "Application started successfully!"
    print_info "Waiting for services to be healthy..."
    sleep 5

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps
    print_info "Application available at: http://localhost:${APP_PORT:-3000}"
}

stop() {
    print_info "Stopping Personal Excalidraw..."

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" down

    print_success "Application stopped successfully!"
}

restart() {
    print_info "Restarting Personal Excalidraw..."
    check_env_file

    # Export all variables from the env file so $DOCKER_COMPOSE can use them
    set -a
    source "$ENV_FILE"
    set +a

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" restart

    print_success "Application restarted successfully!"

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps
}

status() {
    print_info "Checking service status..."

    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps

    echo ""
    print_info "Service health:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}"
}

logs() {
    local service="${1:-}"
    local follow="${2:-}"

    if [ "$follow" = "-f" ] || [ "$follow" = "--follow" ]; then
        if [ -n "$service" ]; then
            print_info "Following logs for service: $service"
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs -f "$service"
        else
            print_info "Following logs for all services..."
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs -f
        fi
    else
        if [ -n "$service" ]; then
            print_info "Showing last 100 lines for service: $service"
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=100 "$service"
        else
            print_info "Showing last 100 lines for all services..."
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=100
        fi
    fi
}

upgrade() {
    print_info "Upgrading Personal Excalidraw..."
    check_env_file

    # Export all variables from the env file so $DOCKER_COMPOSE can use them
    set -a
    source "$ENV_FILE"
    set +a

    # Pull latest changes
    print_info "Pulling latest changes from git..."
    git pull origin master || git pull origin main || {
        print_warning "Git pull failed or not in a git repository"
        print_info "Continuing with local files..."
    }

    # Backup database
    print_info "Creating database backup..."
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" exec -T postgres pg_dump -U postgres personal_excalidraw > "$BACKUP_FILE" 2>/dev/null || {
        print_warning "Database backup failed (database might not be running)"
    }

    if [ -f "$BACKUP_FILE" ]; then
        print_success "Database backed up to: $BACKUP_FILE"
    fi

    # Rebuild and restart
    print_info "Rebuilding and restarting services..."
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d --build

    # Wait for services to be healthy
    print_info "Waiting for services to be healthy..."
    sleep 10

    # Check status
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps

    # Show recent logs
    print_info "Recent application logs:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=20 backend

    print_success "Upgrade completed successfully!"

    # Clean up old images
    print_info "Cleaning up old Docker images..."
    docker image prune -f

    print_success "Cleanup completed!"
}

backup() {
    print_info "Creating database backup..."

    BACKUP_DIR="backups"
    mkdir -p "$BACKUP_DIR"

    BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" exec -T postgres pg_dump -U postgres personal_excalidraw > "$BACKUP_FILE"

    print_success "Database backed up to: $BACKUP_FILE"

    # Keep only last 10 backups
    print_info "Keeping only last 10 backups..."
    ls -t "$BACKUP_DIR"/backup_*.sql | tail -n +11 | xargs -r rm

    print_info "Available backups:"
    ls -lh "$BACKUP_DIR"/backup_*.sql 2>/dev/null || print_warning "No backups found"
}

restore() {
    local backup_file="$1"

    if [ -z "$backup_file" ]; then
        print_error "Please specify a backup file to restore"
        print_info "Usage: $0 restore <backup_file.sql>"
        print_info "Available backups:"
        ls -lh backups/backup_*.sql 2>/dev/null || print_warning "No backups found"
        exit 1
    fi

    if [ ! -f "$backup_file" ]; then
        print_error "Backup file not found: $backup_file"
        exit 1
    fi

    print_warning "This will restore the database from: $backup_file"
    read -p "Are you sure? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        print_info "Restore cancelled"
        exit 0
    fi

    print_info "Restoring database from backup..."
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" exec -T postgres psql -U postgres personal_excalidraw < "$backup_file"

    print_success "Database restored successfully!"

    # Restart application
    print_info "Restarting application..."
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" restart backend

    print_success "Application restarted!"
}

clean() {
    print_warning "This will remove all containers, networks, and images for $PROJECT_NAME"
    print_warning "Database volumes will be preserved"
    read -p "Are you sure? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        print_info "Clean cancelled"
        exit 0
    fi

    print_info "Stopping and removing containers..."
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" down

    print_info "Removing images..."
    docker images | grep personal-excalidraw | awk '{print $3}' | xargs -r docker rmi -f

    print_info "Cleaning up Docker system..."
    docker system prune -f

    print_success "Cleanup completed!"
}

# Usage/help function
usage() {
    cat << EOF
${GREEN}Personal Excalidraw Deployment Script${NC}

${YELLOW}Usage:${NC}
  $0 <command> [options]

${YELLOW}Commands:${NC}
  ${BLUE}config${NC}             Interactive configuration setup (recommended for first-time setup)
  ${BLUE}start${NC}              Start the application
  ${BLUE}stop${NC}               Stop the application
  ${BLUE}restart${NC}            Restart the application
  ${BLUE}status${NC}             Show service status and health
  ${BLUE}logs${NC} [service]     Show logs (last 100 lines)
                       - Add service name: frontend, backend, nginx, postgres
                       - Add -f to follow logs
  ${BLUE}upgrade${NC}            Pull latest changes, backup database, rebuild and restart
  ${BLUE}backup${NC}             Create database backup
  ${BLUE}restore${NC} <file>     Restore database from backup file
  ${BLUE}clean${NC}              Remove all containers and images (keeps database)

${YELLOW}Examples:${NC}
  $0 config                 # Interactive setup (first time)
  $0 start                  # Start all services
  $0 logs backend -f        # Follow backend logs
  $0 logs postgres          # Show last 100 lines of postgres logs
  $0 upgrade                # Upgrade to latest version
  $0 backup                 # Create database backup
  $0 restore backups/backup_20240101_120000.sql
  $0 status                 # Check service health

${YELLOW}First Time Setup:${NC}
  1. Run: $0 config
  2. Follow the prompts
  3. Run: $0 start

${YELLOW}Configuration:${NC}
  Environment file: $ENV_FILE
  Compose file: $COMPOSE_FILE

EOF
}

# Main script logic
case "${1:-}" in
    config)
        config
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs "${2:-}" "${3:-}"
        ;;
    upgrade)
        upgrade
        ;;
    backup)
        backup
        ;;
    restore)
        restore "${2:-}"
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        usage
        ;;
    *)
        print_error "Unknown command: ${1:-}"
        echo ""
        usage
        exit 1
        ;;
esac
