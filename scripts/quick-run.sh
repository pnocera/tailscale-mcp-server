#!/bin/bash

# Quick run script for Tailscale MCP Server
# Usage: ./quick-run.sh [image_tag]

set -euo pipefail

# Default values
DEFAULT_IMAGE="ghcr.io/pnocera/tailscale-mcp-server:latest"
IMAGE="${1:-$DEFAULT_IMAGE}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check environment variables
check_env() {
    local has_api_key=false
    local has_oauth=false
    
    if [ -n "${TAILSCALE_API_KEY:-}" ]; then
        has_api_key=true
        log_success "TAILSCALE_API_KEY is set"
    fi
    
    if [ -n "${TAILSCALE_CLIENT_ID:-}" ] && [ -n "${TAILSCALE_CLIENT_SECRET:-}" ]; then
        has_oauth=true
        log_success "OAuth credentials are set"
    fi
    
    if [ "$has_api_key" = false ] && [ "$has_oauth" = false ]; then
        log_error "No authentication credentials found!"
        log_info "Set either:"
        log_info "  TAILSCALE_API_KEY=tskey-api-..."
        log_info "  OR"
        log_info "  TAILSCALE_CLIENT_ID=... and TAILSCALE_CLIENT_SECRET=..."
        exit 1
    fi
}

# Stop existing container
stop_existing() {
    if docker ps -q -f name=tailscale-mcp-server | grep -q .; then
        log_info "Stopping existing container..."
        docker stop tailscale-mcp-server
        docker rm tailscale-mcp-server
    fi
}

# Run container
run_container() {
    log_info "Running Tailscale MCP Server..."
    log_info "Image: $IMAGE"
    
    # Build docker run command
    local cmd="docker run -d --name tailscale-mcp-server --restart unless-stopped"
    
    # Add environment variables
    if [ -n "${TAILSCALE_API_KEY:-}" ]; then
        cmd="$cmd -e TAILSCALE_API_KEY=$TAILSCALE_API_KEY"
    fi
    
    if [ -n "${TAILSCALE_CLIENT_ID:-}" ]; then
        cmd="$cmd -e TAILSCALE_CLIENT_ID=$TAILSCALE_CLIENT_ID"
    fi
    
    if [ -n "${TAILSCALE_CLIENT_SECRET:-}" ]; then
        cmd="$cmd -e TAILSCALE_CLIENT_SECRET=$TAILSCALE_CLIENT_SECRET"
    fi
    
    if [ -n "${TAILSCALE_TAILNET:-}" ]; then
        cmd="$cmd -e TAILSCALE_TAILNET=$TAILSCALE_TAILNET"
    fi
    
    # Add image
    cmd="$cmd $IMAGE"
    
    # Execute command
    eval $cmd
    
    log_success "Container started successfully!"
    log_info "Container ID: $(docker ps -q -f name=tailscale-mcp-server)"
}

# Show status
show_status() {
    log_info "Container status:"
    docker ps -f name=tailscale-mcp-server --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    log_info "Recent logs:"
    docker logs --tail 10 tailscale-mcp-server
}

# Main execution
main() {
    log_info "Tailscale MCP Server Quick Run"
    log_info "=============================="
    
    check_env
    stop_existing
    run_container
    
    sleep 2
    show_status
    
    log_success "Tailscale MCP Server is running!"
    log_info "View logs: docker logs -f tailscale-mcp-server"
    log_info "Stop server: docker stop tailscale-mcp-server"
}

# Show help
if [[ "${1:-}" == "-h" ]] || [[ "${1:-}" == "--help" ]]; then
    cat << EOF
Tailscale MCP Server Quick Run Script

Usage: $0 [image_tag]

Arguments:
  image_tag    Docker image to run (default: $DEFAULT_IMAGE)

Examples:
  $0                                    # Use default image
  $0 ghcr.io/myuser/tailscale-mcp-server:v1.0.0  # Custom image

Environment Variables (required):
  TAILSCALE_API_KEY     Tailscale API key
  OR
  TAILSCALE_CLIENT_ID   OAuth client ID
  TAILSCALE_CLIENT_SECRET   OAuth client secret

Optional:
  TAILSCALE_TAILNET     Tailnet name (default: "-")

EOF
    exit 0
fi

main "$@"