# Tailscale MCP Server - Docker Build and Push Script (PowerShell)
# Usage: .\build-and-push.ps1 [-Registry] [-Tag] [-Platform]
# Example: .\build-and-push.ps1 -Registry "ghcr.io/username" -Tag "v1.0.0" -Platform "linux/amd64,linux/arm64"

param(
    [string]$Registry = "ghcr.io/pnocera",
    [string]$Tag = "latest",
    [string]$Platform = "linux/amd64,linux/arm64",
    [switch]$Help
)

# Show help if requested
if ($Help) {
    Write-Host @"
Tailscale MCP Server - Docker Build and Push Script (PowerShell)

Usage: .\build-and-push.ps1 [-Registry <registry>] [-Tag <tag>] [-Platform <platform>] [-Help]

Parameters:
  -Registry    Docker registry URL (default: ghcr.io/pnocera)
  -Tag         Image tag (default: latest)
  -Platform    Target platforms (default: linux/amd64,linux/arm64)
  -Help        Show this help message

Examples:
  .\build-and-push.ps1                                                     # Use defaults
  .\build-and-push.ps1 -Registry "ghcr.io/myuser" -Tag "v1.0.0"          # Custom registry and tag
  .\build-and-push.ps1 -Registry "docker.io/myuser" -Platform "linux/amd64"  # Docker Hub, single platform

Environment Variables:
  GITHUB_TOKEN      GitHub Personal Access Token (for ghcr.io)
  GITHUB_ACTOR      GitHub username (for ghcr.io)
  DOCKER_USERNAME   Docker Hub username
  DOCKER_PASSWORD   Docker Hub password

Prerequisites:
  - Docker with buildx support
  - Authentication to target registry
  - Network access to registry
"@
    exit 0
}

# Set error action preference
$ErrorActionPreference = "Stop"

# Constants
$ImageName = "tailscale-mcp-server"
$FullImageName = "$Registry/$ImageName"

# Color functions for output
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Check if Docker is running
function Test-Docker {
    Write-Info "Checking Docker status..."
    try {
        docker info *>$null
        Write-Success "Docker is running"
        return $true
    }
    catch {
        Write-Error "Docker is not running or not accessible"
        Write-Info "Please start Docker Desktop or Docker service"
        return $false
    }
}

# Check if buildx is available
function Test-Buildx {
    Write-Info "Checking Docker buildx..."
    try {
        docker buildx version *>$null
        Write-Success "Docker buildx is available"
        return $true
    }
    catch {
        Write-Error "Docker buildx is not available"
        Write-Info "Please install Docker buildx or use Docker Desktop"
        return $false
    }
}

# Create and use buildx builder
function Initialize-Buildx {
    Write-Info "Setting up buildx builder..."
    
    # Check if builder exists
    $builderExists = $false
    try {
        docker buildx inspect tailscale-mcp-builder *>$null
        $builderExists = $true
        Write-Info "Using existing buildx builder..."
    }
    catch {
        Write-Info "Creating new buildx builder..."
    }
    
    if (-not $builderExists) {
        docker buildx create --name tailscale-mcp-builder --use
    }
    else {
        docker buildx use tailscale-mcp-builder
    }
    
    # Bootstrap the builder
    docker buildx inspect --bootstrap
    Write-Success "Buildx builder ready"
}

# Login to registry
function Connect-Registry {
    Write-Info "Checking registry authentication..."
    
    if ($Registry -like "ghcr.io*") {
        Write-Info "GitHub Container Registry detected"
        $githubToken = $env:GITHUB_TOKEN
        $githubActor = $env:GITHUB_ACTOR
        
        if (-not $githubToken) {
            Write-Warning "GITHUB_TOKEN environment variable not set"
            Write-Info "You may need to login manually:"
            Write-Info "echo `$env:GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin"
        }
        else {
            if (-not $githubActor) {
                $githubActor = $env:USERNAME
            }
            $githubToken | docker login ghcr.io -u $githubActor --password-stdin
            Write-Success "Logged into GitHub Container Registry"
        }
    }
    elseif ($Registry -like "*docker.io*" -or $Registry -like "*docker.com*") {
        Write-Info "Docker Hub detected"
        $dockerUsername = $env:DOCKER_USERNAME
        $dockerPassword = $env:DOCKER_PASSWORD
        
        if (-not $dockerPassword) {
            Write-Warning "DOCKER_PASSWORD environment variable not set"
            Write-Info "You may need to login manually:"
            Write-Info "docker login -u USERNAME"
        }
        else {
            $dockerPassword | docker login -u $dockerUsername --password-stdin
            Write-Success "Logged into Docker Hub"
        }
    }
    else {
        Write-Warning "Unknown registry. Make sure you're logged in:"
        Write-Info "docker login $Registry"
    }
}

# Build and push image
function Build-AndPushImage {
    Write-Info "Building and pushing image..."
    Write-Info "Registry: $Registry"
    Write-Info "Image: $ImageName"
    Write-Info "Tag: $Tag"
    Write-Info "Platform: $Platform"
    Write-Info "Full image name: ${FullImageName}:$Tag"
    
    # Build and push
    docker buildx build `
        --platform $Platform `
        --tag "${FullImageName}:$Tag" `
        --tag "${FullImageName}:latest" `
        --push `
        --progress=plain `
        .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Image built and pushed successfully!"
        Write-Info "Image: ${FullImageName}:$Tag"
        Write-Info "Image: ${FullImageName}:latest"
    }
    else {
        Write-Error "Build failed with exit code $LASTEXITCODE"
        throw "Docker build failed"
    }
}

# Get image information
function Get-ImageInfo {
    Write-Info "Getting image information..."
    
    try {
        docker buildx imagetools inspect "${FullImageName}:$Tag"
    }
    catch {
        Write-Warning "Could not retrieve image information"
    }
}

# Cleanup function
function Invoke-Cleanup {
    Write-Info "Cleaning up..."
    # Remove builder if we created it
    # docker buildx rm tailscale-mcp-builder 2>$null
}

# Main execution
function Main {
    try {
        Write-Info "Starting Tailscale MCP Server Docker build and push"
        Write-Info "=============================================="
        
        # Pre-flight checks
        if (-not (Test-Docker)) { exit 1 }
        if (-not (Test-Buildx)) { exit 1 }
        
        # Setup
        Initialize-Buildx
        Connect-Registry
        
        # Build and push
        Build-AndPushImage
        
        # Post-build info
        Get-ImageInfo
        
        Write-Success "Build and push completed successfully!"
        Write-Info "You can now use the image with:"
        Write-Info "docker run -e TAILSCALE_API_KEY=your-key ${FullImageName}:$Tag"
    }
    catch {
        Write-Error "Script failed: $_"
        exit 1
    }
    finally {
        Invoke-Cleanup
    }
}

# Run main function
Main