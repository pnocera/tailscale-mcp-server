# Tailscale MCP Server - Deployment Scripts

This directory contains automation scripts for building, deploying, and managing the Tailscale MCP Server Docker containers.

## Scripts Overview

### Build and Push Scripts

#### `build-and-push.sh` (Bash)
Multi-platform Docker build and registry push script for Unix-like systems.

**Features:**
- Multi-architecture builds (AMD64, ARM64)
- Support for GitHub Container Registry, Docker Hub, and custom registries
- Automatic authentication handling
- Progress monitoring and error handling
- Buildx setup and management

**Usage:**
```bash
# Default build and push
./build-and-push.sh

# Custom registry and tag
./build-and-push.sh ghcr.io/myuser v1.0.0

# Single platform build
./build-and-push.sh docker.io/myuser latest linux/amd64

# Show help
./build-and-push.sh --help
```

**Environment Variables:**
- `GITHUB_TOKEN` - GitHub Personal Access Token (for ghcr.io)
- `GITHUB_ACTOR` - GitHub username (for ghcr.io)
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password

#### `build-and-push.ps1` (PowerShell)
Windows PowerShell version of the build and push script.

**Usage:**
```powershell
# Default build and push
.\\build-and-push.ps1

# Custom registry and tag
.\\build-and-push.ps1 -Registry \"ghcr.io/myuser\" -Tag \"v1.0.0\"

# Single platform build
.\\build-and-push.ps1 -Registry \"docker.io/myuser\" -Platform \"linux/amd64\"

# Show help
.\\build-and-push.ps1 -Help
```

### Quick Run Scripts

#### `quick-run.sh` (Bash)
Instant container deployment script with environment validation.

**Features:**
- Automatic environment variable validation
- Container lifecycle management
- Real-time status monitoring
- Support for both API key and OAuth authentication

**Usage:**
```bash
# Use default image
./quick-run.sh

# Use custom image
./quick-run.sh ghcr.io/myuser/tailscale-mcp-server:v1.0.0

# Show help
./quick-run.sh --help
```

#### `quick-run.ps1` (PowerShell)
Windows PowerShell version of the quick run script.

**Usage:**
```powershell
# Use default image
.\\quick-run.ps1

# Use custom image
.\\quick-run.ps1 -Image \"ghcr.io/myuser/tailscale-mcp-server:v1.0.0\"

# Show help
.\\quick-run.ps1 -Help
```

## Authentication Setup

### GitHub Container Registry (ghcr.io)

1. Create a GitHub Personal Access Token:
   - Go to GitHub → Settings → Developer settings → Personal access tokens
   - Create token with `write:packages` scope

2. Set environment variables:
   ```bash
   export GITHUB_TOKEN=\"ghp_your_token_here\"
   export GITHUB_ACTOR=\"your_github_username\"
   ```

3. Login manually (alternative):
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_ACTOR --password-stdin
   ```

### Docker Hub

1. Set environment variables:
   ```bash
   export DOCKER_USERNAME=\"your_docker_username\"
   export DOCKER_PASSWORD=\"your_docker_password\"
   ```

2. Login manually (alternative):
   ```bash
   docker login -u your_docker_username
   ```

### Custom Registry

1. Login manually:
   ```bash
   docker login your-registry.com
   ```

## Example Workflows

### Development Workflow

1. **Test locally:**
   ```bash
   # Set environment variables
   export TAILSCALE_API_KEY=\"tskey-api-...\"
   
   # Quick run for testing
   ./quick-run.sh
   
   # Check logs
   docker logs -f tailscale-mcp-server
   ```

2. **Build and test custom image:**
   ```bash
   # Build locally
   docker build -t tailscale-mcp-server:test .
   
   # Test with custom image
   ./quick-run.sh tailscale-mcp-server:test
   ```

### Production Deployment

1. **Push to registry:**
   ```bash
   # Set registry credentials
   export GITHUB_TOKEN=\"ghp_...\"
   export GITHUB_ACTOR=\"username\"
   
   # Build and push
   ./build-and-push.sh ghcr.io/myuser v1.0.0
   ```

2. **Deploy on production server:**
   ```bash
   # Set Tailscale credentials
   export TAILSCALE_API_KEY=\"tskey-api-...\"
   
   # Deploy from registry
   ./quick-run.sh ghcr.io/myuser/tailscale-mcp-server:v1.0.0
   ```

### CI/CD Integration

**GitHub Actions example:**
```yaml
- name: Build and Push
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    ./scripts/build-and-push.sh ghcr.io/${{ github.repository_owner }} ${{ github.ref_name }}
```

**GitLab CI example:**
```yaml
script:
  - export DOCKER_USERNAME=$CI_REGISTRY_USER
  - export DOCKER_PASSWORD=$CI_REGISTRY_PASSWORD
  - ./scripts/build-and-push.sh $CI_REGISTRY_IMAGE $CI_COMMIT_TAG
```

## Troubleshooting

### Common Issues

1. **\"Docker not running\"**
   - Start Docker Desktop or Docker service
   - Check: `docker info`

2. **\"buildx not available\"**
   - Install Docker buildx
   - Update Docker Desktop

3. **\"Authentication failed\"**
   - Check environment variables
   - Verify token permissions
   - Try manual login

4. **\"No authentication credentials\"**
   - Set either `TAILSCALE_API_KEY` or both `TAILSCALE_CLIENT_ID` and `TAILSCALE_CLIENT_SECRET`

### Debug Mode

```bash
# Enable debug output
set -x  # In bash scripts

# Or run with verbose docker output
docker build --progress=plain .
```

## Script Customization

All scripts support customization through environment variables and command-line arguments. Check the help output of each script for full details:

```bash
./build-and-push.sh --help
./quick-run.sh --help
```

```powershell
.\\build-and-push.ps1 -Help
.\\quick-run.ps1 -Help
```