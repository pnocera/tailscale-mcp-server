# Tailscale MCP Server

An MCP (Model Context Protocol) server for managing Tailscale resources using the official Tailscale Go client library v2. This server provides complete coverage of the Tailscale API with enhanced, self-descriptive tools powered by OpenAPI documentation.

## 🚀 Features

This MCP server provides **42 comprehensive tools** organized into logical categories, each with detailed descriptions, OAuth scopes, use cases, and security considerations:

### 🖥️ Device Management (9 tools)
- **tailscale_devices_list** - List all devices with optional detailed fields
- **tailscale_device_get** - Get comprehensive device information
- **tailscale_device_delete** - Permanently remove devices from tailnet
- **tailscale_device_authorize** - Authorize/deauthorize devices for access control
- **tailscale_device_set_name** - Set device names (affects Magic DNS)
- **tailscale_device_set_tags** - Assign tags for ACL-based access control
- **tailscale_device_expire** - Force device re-authentication
- **tailscale_device_routes_list** - List subnet routes and exit node configuration
- **tailscale_device_routes_set** - Configure subnet routing and exit nodes

### 🔐 Key Management (4 tools)
- **tailscale_keys_list** - List all authentication keys with capabilities
- **tailscale_key_get** - Get detailed key information and usage statistics
- **tailscale_key_create** - Create reusable, ephemeral, or preauthorized keys
- **tailscale_key_delete** - Revoke authentication keys

### 👥 User Management (8 tools)
- **tailscale_users_list** - List all users with roles and status
- **tailscale_user_get** - Get detailed user profile information
- **tailscale_user_approve** - Approve users for tailnet access
- **tailscale_user_suspend** - Temporarily suspend user access
- **tailscale_user_restore** - Restore suspended users
- **tailscale_user_delete** - Permanently remove users
- **tailscale_contacts_get** - Get tailnet contact preferences
- **tailscale_contact_update** - Update contact information for notifications

### 🌐 DNS Management (9 tools)
- **tailscale_dns_nameservers_get** - Get configured DNS nameservers
- **tailscale_dns_nameservers_set** - Set custom DNS nameservers
- **tailscale_dns_preferences_get** - Get MagicDNS and DNS preferences
- **tailscale_dns_preferences_set** - Configure MagicDNS and DNS behavior
- **tailscale_dns_searchpaths_get** - Get DNS search domain suffixes
- **tailscale_dns_searchpaths_set** - Set DNS search paths for short names
- **tailscale_policy_get** - Get current ACL policy file (HuJSON)
- **tailscale_policy_set** - Update ACL policy with security rules
- **tailscale_policy_validate** - Validate policy files before deployment

### 🔗 Advanced Features (12 tools)
- **tailscale_webhooks_list** - List webhook endpoints for event notifications
- **tailscale_webhook_create** - Create webhooks for external integrations
- **tailscale_webhook_get** - Get webhook configuration and statistics
- **tailscale_webhook_delete** - Remove webhook endpoints
- **tailscale_logging_configuration_get** - Get audit log streaming configuration
- **tailscale_logging_network_get** - Get network flow log configuration
- **tailscale_device_posture_integrations_list** - List security posture integrations
- **tailscale_device_posture_integration_create** - Create posture provider integrations
- **tailscale_device_posture_integration_get** - Get posture integration details
- **tailscale_device_posture_integration_delete** - Remove posture integrations
- **tailscale_tailnet_settings_get** - Get comprehensive tailnet settings
- **tailscale_tailnet_settings_update** - Update tailnet configuration

## 📦 Installation

### Prerequisites
- Go 1.24 or later
- Valid Tailscale account with API access
- Tailscale API key or OAuth client credentials

### Build from Source

1. Clone the repository and navigate to the MCP directory:
```bash
git clone <repository-url>
cd mcp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build -o tailscale-mcp-server ./cmd
```

### Binary Installation
```bash
# Download and install directly
go install github.com/pnocera/tailscale-mcp-server/cmd@latest
```

## ⚙️ Configuration

The server supports both API key and OAuth authentication methods for maximum flexibility.

### Environment Variables

#### API Key Authentication (Recommended for personal use)
```bash
export TAILSCALE_API_KEY="tskey-api-..."
export TAILSCALE_TAILNET="your-tailnet-name"  # Optional, defaults to "-"
```

#### OAuth Authentication (Recommended for applications)
```bash
export TAILSCALE_CLIENT_ID="your-oauth-client-id"
export TAILSCALE_CLIENT_SECRET="your-oauth-client-secret"
export TAILSCALE_TAILNET="your-tailnet-name"  # Optional, defaults to "-"
```

### Authentication Priority
1. If both `TAILSCALE_CLIENT_ID` and `TAILSCALE_CLIENT_SECRET` are set, OAuth is used
2. Otherwise, API key authentication is used with `TAILSCALE_API_KEY`

## 🚀 Usage

### Running the Server

```bash
# Using API key authentication
TAILSCALE_API_KEY="tskey-api-..." ./tailscale-mcp-server

# Using OAuth authentication
TAILSCALE_CLIENT_ID="..." TAILSCALE_CLIENT_SECRET="..." ./tailscale-mcp-server

# With custom tailnet
TAILSCALE_API_KEY="tskey-api-..." TAILSCALE_TAILNET="mycompany.com" ./tailscale-mcp-server
```

### MCP Client Integration

#### Claude Code Integration
```bash
# Add to Claude Code
claude mcp add tailscale /path/to/tailscale-mcp-server
```

#### Generic MCP Client Configuration
```json
{
  "mcpServers": {
    "tailscale": {
      "command": "/path/to/tailscale-mcp-server",
      "env": {
        "TAILSCALE_API_KEY": "tskey-api-...",
        "TAILSCALE_TAILNET": "your-tailnet"
      }
    }
  }
}
```

## 📚 Tool Examples

### Device Management
```json
// List all devices with detailed information
{
  "name": "tailscale_devices_list",
  "arguments": {
    "fields": "all"
  }
}

// Get specific device details
{
  "name": "tailscale_device_get",
  "arguments": {
    "device_id": "device-id-here",
    "fields": "all"
  }
}

// Authorize a device
{
  "name": "tailscale_device_authorize",
  "arguments": {
    "device_id": "device-id-here",
    "authorized": true
  }
}

// Set device tags for ACL-based access control
{
  "name": "tailscale_device_set_tags",
  "arguments": {
    "device_id": "device-id-here",
    "tags": ["tag:server", "tag:production"]
  }
}
```

### Key Management
```json
// Create a reusable, preauthorized key with tags
{
  "name": "tailscale_key_create",
  "arguments": {
    "reusable": true,
    "ephemeral": false,
    "preauthorized": true,
    "description": "CI/CD deployment key",
    "tags": ["tag:ci", "tag:automated"],
    "expiry_seconds": 86400
  }
}

// List all authentication keys
{
  "name": "tailscale_keys_list",
  "arguments": {}
}
```

### DNS Configuration
```json
// Set custom DNS nameservers
{
  "name": "tailscale_dns_nameservers_set",
  "arguments": {
    "nameservers": ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
  }
}

// Enable MagicDNS for easy device connectivity
{
  "name": "tailscale_dns_preferences_set",
  "arguments": {
    "magic_dns": true
  }
}

// Set DNS search paths for short hostnames
{
  "name": "tailscale_dns_searchpaths_set",
  "arguments": {
    "search_paths": ["company.com", "internal.local"]
  }
}
```

### Policy Management
```json
// Get current ACL policy
{
  "name": "tailscale_policy_get",
  "arguments": {}
}

// Validate policy before applying
{
  "name": "tailscale_policy_validate",
  "arguments": {
    "policy": "{\n  \"acls\": [\n    {\n      \"action\": \"accept\",\n      \"src\": [\"tag:server\"],\n      \"dst\": [\"tag:database:5432\"]\n    }\n  ]\n}"
  }
}

// Update ACL policy
{
  "name": "tailscale_policy_set",
  "arguments": {
    "policy": "{\n  \"acls\": [\n    {\n      \"action\": \"accept\",\n      \"src\": [\"*\"],\n      \"dst\": [\"*:*\"]\n    }\n  ]\n}"
  }
}
```

### Webhooks & Integrations
```json
// Create webhook for device events
{
  "name": "tailscale_webhook_create",
  "arguments": {
    "endpoint_url": "https://your-app.com/webhook",
    "subscriptions": ["device.created", "device.deleted", "user.approved"]
  }
}

// Create device posture integration
{
  "name": "tailscale_device_posture_integration_create",
  "arguments": {
    "provider": "crowdstrike",
    "client_id": "your-client-id",
    "client_secret": "your-client-secret",
    "tenant_id": "your-tenant-id"
  }
}
```

## 🏗️ Architecture

The server follows a clean, modular architecture:

```
├── cmd/
│   └── main.go                 # Entry point and server setup
├── internal/
│   ├── config/                 # Configuration management
│   ├── client/                 # Tailscale client wrapper
│   └── handlers/               # MCP request handlers
├── pkg/
│   └── tools/                  # Tool implementations
│       ├── devices.go          # Device management (9 tools)
│       ├── keys.go             # Key management (4 tools)
│       ├── users.go            # User & contact management (8 tools)
│       ├── dns.go              # DNS & policy management (9 tools)
│       └── additional.go       # Advanced features (12 tools)
├── tailscale_api_docs/         # OpenAPI documentation
├── .gitignore                  # Git ignore rules
├── LICENSE.md                  # MIT License
└── README.md                   # This file
```

### Key Design Principles
- **Modular**: Each tool category is organized in separate files
- **Self-descriptive**: Tools include comprehensive descriptions from OpenAPI docs
- **Type-safe**: Full Go type safety with structured request/response handling
- **Error-resilient**: Comprehensive error handling with informative messages
- **OAuth-ready**: Support for both API key and OAuth authentication

## 🔐 Authentication & Security

### OAuth Scopes
Each tool specifies the required OAuth scope in its description:
- `devices:read` / `devices:write` - Device management
- `keys:read` / `keys:write` - Authentication key management
- `users:read` / `users:write` - User management
- `dns:read` / `dns:write` - DNS configuration
- `acl:read` / `acl:write` - ACL policy management
- `webhooks:read` / `webhooks:write` - Webhook management
- `logging:read` - Log configuration access
- `posture:read` / `posture:write` - Device posture management
- `settings:read` / `settings:write` - Tailnet settings

### Security Best Practices
- Store API keys and OAuth credentials securely
- Use environment variables for sensitive configuration
- Implement proper access controls in your MCP client
- Regularly rotate API keys and OAuth credentials
- Monitor API usage through Tailscale admin console

## 🛠️ Development

### Adding New Tools

1. **Identify the OpenAPI endpoint** in `tailscale_api_docs/tailscaleapi.yaml`
2. **Choose the appropriate file** in `pkg/tools/` based on functionality
3. **Add the tool definition** in the `RegisterTools` method:
```go
tool := mcp.NewTool(
    "tailscale_new_tool",
    mcp.WithDescription("Detailed description with OAuth scope and use cases"),
    mcp.WithString("param", mcp.Description("Parameter description"), mcp.Required()),
)
mcpServer.AddTool(tool, dt.NewToolHandler)
```
4. **Implement the handler function** following existing patterns
5. **Test thoroughly** and update documentation

### Enhanced Tool Descriptions
All tools include:
- **Detailed functionality** description
- **OAuth scope** requirements
- **Use cases** and examples
- **Security considerations**
- **Links to Tailscale documentation**

### Testing

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o tailscale-mcp-server-linux ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o tailscale-mcp-server-macos ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o tailscale-mcp-server.exe ./cmd
```

## 📊 Monitoring & Observability

### Built-in Logging
The server provides structured logging for:
- Authentication attempts
- API requests and responses
- Error conditions
- Performance metrics

### Integration with Tailscale
- Monitor API usage in the Tailscale admin console
- Track OAuth token usage and refresh cycles
- Review audit logs for security compliance

## 🔗 Dependencies

- **[tailscale.com/client/tailscale/v2](https://pkg.go.dev/tailscale.com/client/tailscale/v2)** - Official Tailscale Go client library
- **[github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)** - MCP protocol implementation for Go
- **[golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)** - OAuth 2.0 client library
- **Standard Go libraries** - JSON, HTTP, context, logging

## 📄 License

This project is licensed under the MIT License. See [LICENSE.md](LICENSE.md) for details.

## 🤝 Contributing

Contributions are welcome! Please ensure all new tools include:

- **Complete input validation** with proper error messages
- **Comprehensive error handling** for all failure scenarios
- **Detailed descriptions** following the OpenAPI documentation pattern
- **JSON response formatting** consistent with existing tools
- **OAuth scope specifications** in tool descriptions
- **Unit tests** for core functionality
- **Documentation updates** in this README

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-tool`
3. Implement your changes with tests
4. Run the test suite: `go test ./...`
5. Update documentation as needed
6. Submit a pull request with a clear description

## 📚 Resources

- **[Tailscale API Documentation](https://tailscale.com/api)** - Official API reference
- **[MCP Protocol Specification](https://spec.modelcontextprotocol.io/)** - MCP protocol details
- **[Tailscale Knowledge Base](https://tailscale.com/kb/)** - Comprehensive guides and tutorials
- **[Go Client Library Documentation](https://pkg.go.dev/tailscale.com/client/tailscale/v2)** - Official Go client docs

## 🆘 Support

- **Issues**: Report bugs and request features on GitHub
- **Documentation**: Refer to the Tailscale Knowledge Base
- **Community**: Join the Tailscale community forums

---

Made with ❤️ for the Tailscale and MCP communities. This server provides the most comprehensive Tailscale MCP integration available, with self-descriptive tools powered by official OpenAPI documentation.