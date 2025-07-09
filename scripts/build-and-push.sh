#!/bin/bash

# Tailscale MCP Server - Docker Build and Push Script
# Usage: ./build-and-push.sh [registry] [tag] [platform]
# Example: ./build-and-push.sh ghcr.io/username v1.0.0 linux/amd64,linux/arm64

set -euo pipefail

# Default values
DEFAULT_REGISTRY="ghcr.io/pnocera"
DEFAULT_TAG="latest"
DEFAULT_PLATFORM="linux/amd64,linux/arm64"

# Parse arguments
REGISTRY="${1:-$DEFAULT_REGISTRY}"
TAG="${2:-$DEFAULT_TAG}"
PLATFORM="${3:-$DEFAULT_PLATFORM}"

# Image name
IMAGE_NAME="tailscale-mcp-server"
FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker is not running or not accessible"
        exit 1
    fi
    log_success "Docker is running"
}

# Check if buildx is available
check_buildx() {
    if ! docker buildx version >/dev/null 2>&1; then
        log_error "Docker buildx is not available"
        log_info "Please install Docker buildx or use Docker Desktop"
        exit 1
    fi
    log_success "Docker buildx is available"
}

# Create and use buildx builder
setup_buildx() {
    log_info "Setting up buildx builder..."
    
    # Create builder if it doesn't exist
    if ! docker buildx inspect tailscale-mcp-builder >/dev/null 2>&1; then
        log_info "Creating new buildx builder..."
        docker buildx create --name tailscale-mcp-builder --use
    else
        log_info "Using existing buildx builder..."
        docker buildx use tailscale-mcp-builder
    fi
    
    # Bootstrap the builder
    docker buildx inspect --bootstrap
    log_success "Buildx builder ready"
}

# Login to registry
login_registry() {
    log_info "Checking registry authentication..."
    
    if [[ "$REGISTRY" == ghcr.io* ]]; then
        log_info "GitHub Container Registry detected"
        if [ -z "${GITHUB_TOKEN:-}" ]; then
            log_warning "GITHUB_TOKEN not set. You may need to login manually:"
            log_info "echo \$GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin"
        else
            echo "$GITHUB_TOKEN" | docker login ghcr.io -u "${GITHUB_ACTOR:-$(whoami)}" --password-stdin
            log_success "Logged into GitHub Container Registry"
        fi
    elif [[ "$REGISTRY" == *docker.io* ]] || [[ "$REGISTRY" == *docker.com* ]]; then
        log_info "Docker Hub detected"
        if [ -z "${DOCKER_PASSWORD:-}" ]; then
            log_warning "DOCKER_PASSWORD not set. You may need to login manually:"
            log_info "docker login -u USERNAME"
        else
            echo "$DOCKER_PASSWORD" | docker login -u "${DOCKER_USERNAME:-}" --password-stdin
            log_success "Logged into Docker Hub"
        fi
    else
        log_warning "Unknown registry. Make sure you're logged in:"
        log_info "docker login $REGISTRY"
    fi
}

# Build and push image
build_and_push() {
    log_info "Building and pushing image..."
    log_info "Registry: $REGISTRY"
    log_info "Image: $IMAGE_NAME"
    log_info "Tag: $TAG"
    log_info "Platform: $PLATFORM"
    log_info "Full image name: $FULL_IMAGE_NAME:$TAG"
    
    # Build and push
    docker buildx build \
        --platform "$PLATFORM" \
        --tag "$FULL_IMAGE_NAME:$TAG" \
        --tag "$FULL_IMAGE_NAME:latest" \
        --push \
        --progress=plain \
        .
    
    log_success "Image built and pushed successfully!"
    log_info "Image: $FULL_IMAGE_NAME:$TAG"
    log_info "Image: $FULL_IMAGE_NAME:latest"
}

# Get image information
get_image_info() {
    log_info "Getting image information..."
    
    # Get image size and manifest
    docker buildx imagetools inspect "$FULL_IMAGE_NAME:$TAG" || true
}

# Cleanup function
cleanup() {
    log_info "Cleaning up..."
    # Remove builder if we created it
    # docker buildx rm tailscale-mcp-builder || true
}

# Main execution
main() {
    log_info "Starting Tailscale MCP Server Docker build and push"
    log_info "=============================================="
    
    # Pre-flight checks
    check_docker
    check_buildx
    
    # Setup
    setup_buildx
    login_registry
    
    # Build and push
    build_and_push
    
    # Post-build info
    get_image_info
    
    log_success "Build and push completed successfully!"
    log_info "You can now use the image with:"
    log_info "docker run -e TAILSCALE_API_KEY=your-key $FULL_IMAGE_NAME:$TAG"
}

# Show help
show_help() {
    cat << EOF
Tailscale MCP Server - Docker Build and Push Script

Usage: $0 [registry] [tag] [platform]

Arguments:
  registry    Docker registry URL (default: $DEFAULT_REGISTRY)
  tag         Image tag (default: $DEFAULT_TAG)
  platform    Target platforms (default: $DEFAULT_PLATFORM)

Examples:
  $0                                                    # Use defaults
  $0 ghcr.io/myuser v1.0.0                            # Custom registry and tag
  $0 docker.io/myuser latest linux/amd64              # Docker Hub, single platform
  $0 myregistry.com/myuser v2.0.0 linux/amd64,linux/arm64  # Custom registry, multi-platform

Environment Variables:
  GITHUB_TOKEN      GitHub Personal Access Token (for ghcr.io)
  GITHUB_ACTOR      GitHub username (for ghcr.io)
  DOCKER_USERNAME   Docker Hub username
  DOCKER_PASSWORD   Docker Hub password

Prerequisites:
  - Docker with buildx support
  - Authentication to target registry
  - Network access to registry

EOF
}

# Handle script arguments
if [[ "${1:-}" == "-h" ]] || [[ "${1:-}" == "--help" ]]; then
    show_help
    exit 0
fi

# Set trap for cleanup
trap cleanup EXIT

# Run main function
main "$@"