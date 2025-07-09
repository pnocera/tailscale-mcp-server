# Quick run script for Tailscale MCP Server (PowerShell)
# Usage: .\quick-run.ps1 [-Image <image_tag>]

param(
    [string]$Image = "ghcr.io/pnocera/tailscale-mcp-server:latest",
    [switch]$Help
)

if ($Help) {
    Write-Host @"
Tailscale MCP Server Quick Run Script (PowerShell)

Usage: .\quick-run.ps1 [-Image <image_tag>] [-Help]

Parameters:
  -Image    Docker image to run (default: ghcr.io/pnocera/tailscale-mcp-server:latest)
  -Help     Show this help message

Examples:
  .\quick-run.ps1                                                    # Use default image
  .\quick-run.ps1 -Image "ghcr.io/myuser/tailscale-mcp-server:v1.0.0"  # Custom image

Environment Variables (required):
  TAILSCALE_API_KEY     Tailscale API key
  OR
  TAILSCALE_CLIENT_ID   OAuth client ID
  TAILSCALE_CLIENT_SECRET   OAuth client secret

Optional:
  TAILSCALE_TAILNET     Tailnet name (default: "-")
"@
    exit 0
}

$ErrorActionPreference = "Stop"

# Color functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

# Check environment variables
function Test-Environment {
    $hasApiKey = $false
    $hasOAuth = $false
    
    if ($env:TAILSCALE_API_KEY) {
        $hasApiKey = $true
        Write-Success "TAILSCALE_API_KEY is set"
    }
    
    if ($env:TAILSCALE_CLIENT_ID -and $env:TAILSCALE_CLIENT_SECRET) {
        $hasOAuth = $true
        Write-Success "OAuth credentials are set"
    }
    
    if (-not $hasApiKey -and -not $hasOAuth) {
        Write-Error "No authentication credentials found!"
        Write-Info "Set either:"
        Write-Info "  `$env:TAILSCALE_API_KEY = 'tskey-api-...'"
        Write-Info "  OR"
        Write-Info "  `$env:TAILSCALE_CLIENT_ID = '...' and `$env:TAILSCALE_CLIENT_SECRET = '...'"
        exit 1
    }
}

# Stop existing container
function Stop-ExistingContainer {
    $existingContainer = docker ps -q -f name=tailscale-mcp-server 2>$null
    if ($existingContainer) {
        Write-Info "Stopping existing container..."
        docker stop tailscale-mcp-server
        docker rm tailscale-mcp-server
    }
}

# Run container
function Start-Container {
    Write-Info "Running Tailscale MCP Server..."
    Write-Info "Image: $Image"
    
    # Build docker run command arguments
    $dockerArgs = @(
        "run", "-d",
        "--name", "tailscale-mcp-server",
        "--restart", "unless-stopped"
    )
    
    # Add environment variables
    if ($env:TAILSCALE_API_KEY) {
        $dockerArgs += @("-e", "TAILSCALE_API_KEY=$($env:TAILSCALE_API_KEY)")
    }
    
    if ($env:TAILSCALE_CLIENT_ID) {
        $dockerArgs += @("-e", "TAILSCALE_CLIENT_ID=$($env:TAILSCALE_CLIENT_ID)")
    }
    
    if ($env:TAILSCALE_CLIENT_SECRET) {
        $dockerArgs += @("-e", "TAILSCALE_CLIENT_SECRET=$($env:TAILSCALE_CLIENT_SECRET)")
    }
    
    if ($env:TAILSCALE_TAILNET) {
        $dockerArgs += @("-e", "TAILSCALE_TAILNET=$($env:TAILSCALE_TAILNET)")
    }
    
    # Add image
    $dockerArgs += $Image
    
    # Execute command
    $containerId = docker @dockerArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Container started successfully!"
        Write-Info "Container ID: $containerId"
    }
    else {
        Write-Error "Failed to start container"
        exit 1
    }
}

# Show status
function Show-Status {
    Write-Info "Container status:"
    docker ps -f name=tailscale-mcp-server --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    Write-Info "Recent logs:"
    docker logs --tail 10 tailscale-mcp-server
}

# Main execution
function Main {
    try {
        Write-Info "Tailscale MCP Server Quick Run"
        Write-Info "=============================="
        
        Test-Environment
        Stop-ExistingContainer
        Start-Container
        
        Start-Sleep -Seconds 2
        Show-Status
        
        Write-Success "Tailscale MCP Server is running!"
        Write-Info "View logs: docker logs -f tailscale-mcp-server"
        Write-Info "Stop server: docker stop tailscale-mcp-server"
    }
    catch {
        Write-Error "Script failed: $_"
        exit 1
    }
}

Main