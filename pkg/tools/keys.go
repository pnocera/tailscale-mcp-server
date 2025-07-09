package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pnocera/tailscale-mcp-server/internal/client"
	"tailscale.com/client/tailscale/v2"
)

type KeyTools struct {
	client *client.TailscaleClient
}

func NewKeyTools(client *client.TailscaleClient) *KeyTools {
	return &KeyTools{client: client}
}

func (kt *KeyTools) RegisterTools(mcpServer *server.MCPServer) {
	tool := mcp.NewTool(
		"tailscale_keys_list",
		mcp.WithDescription("List all authentication keys for the tailnet. Returns all auth keys including reusable keys, ephemeral keys, and tagged keys. Shows key status, expiration times, usage counts, and associated capabilities. Essential for managing device onboarding and access control. OAuth Scope: keys:read."),
	)
	mcpServer.AddTool(tool, kt.ListKeys)

	tool = mcp.NewTool(
		"tailscale_key_get",
		mcp.WithDescription("Get detailed information about a specific authentication key. Returns key capabilities, creation time, expiration status, usage count, and associated tags. Use this to verify key permissions and monitor key usage for security auditing. OAuth Scope: keys:read."),
		mcp.WithString("key_id", mcp.Description("The key ID"), mcp.Required()),
	)
	mcpServer.AddTool(tool, kt.GetKey)

	tool = mcp.NewTool(
		"tailscale_key_create",
		mcp.WithDescription("Create a new authentication key for device onboarding. Configure key as reusable (multiple devices), ephemeral (temporary devices), or preauthorized (automatic approval). Set expiration time and assign tags for ACL-based access control. Essential for automated device deployment and CI/CD integration. OAuth Scope: keys:write."),
		mcp.WithBoolean("reusable", mcp.Description("Whether the key can be reused"), mcp.DefaultBool(false)),
		mcp.WithBoolean("ephemeral", mcp.Description("Whether devices using this key will be ephemeral"), mcp.DefaultBool(false)),
		mcp.WithBoolean("preauthorized", mcp.Description("Whether devices using this key will be pre-authorized"), mcp.DefaultBool(false)),
		mcp.WithString("description", mcp.Description("Description of the key")),
		mcp.WithArray("tags", mcp.Description("Tags to apply to devices using this key"), mcp.WithStringItems()),
		mcp.WithNumber("expiry_seconds", mcp.Description("Expiry time in seconds from now")),
	)
	mcpServer.AddTool(tool, kt.CreateKey)

	tool = mcp.NewTool(
		"tailscale_key_delete",
		mcp.WithDescription("Delete an authentication key to revoke its ability to add new devices. This does not affect devices already authenticated with this key. Use this to clean up unused keys or revoke compromised keys. Essential for maintaining security hygiene and key lifecycle management. OAuth Scope: keys:write."),
		mcp.WithString("key_id", mcp.Description("The key ID to delete"), mcp.Required()),
	)
	mcpServer.AddTool(tool, kt.DeleteKey)
}

func (kt *KeyTools) ListKeys(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := kt.client.GetClient()
	keys, err := client.Keys().List(ctx, false)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list keys: %v", err)), nil
	}

	keysJSON, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal keys: %v", err)), nil
	}

	return mcp.NewToolResultText(string(keysJSON)), nil
}

func (kt *KeyTools) GetKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		KeyID string `json:"key_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := kt.client.GetClient()
	key, err := client.Keys().Get(ctx, args.KeyID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get key: %v", err)), nil
	}

	keyJSON, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal key: %v", err)), nil
	}

	return mcp.NewToolResultText(string(keyJSON)), nil
}

func (kt *KeyTools) CreateKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Reusable      bool     `json:"reusable"`
		Ephemeral     bool     `json:"ephemeral"`
		Preauthorized bool     `json:"preauthorized"`
		Description   string   `json:"description"`
		Tags          []string `json:"tags"`
		ExpirySeconds int      `json:"expiry_seconds"`
	}

	if request.Params.Arguments != nil {
		if err := request.BindArguments(&args); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
		}
	}

	createReq := tailscale.CreateKeyRequest{
		Capabilities: tailscale.KeyCapabilities{
			Devices: struct {
				Create struct {
					Reusable      bool     `json:"reusable"`
					Ephemeral     bool     `json:"ephemeral"`
					Tags          []string `json:"tags"`
					Preauthorized bool     `json:"preauthorized"`
				} `json:"create"`
			}{
				Create: struct {
					Reusable      bool     `json:"reusable"`
					Ephemeral     bool     `json:"ephemeral"`
					Tags          []string `json:"tags"`
					Preauthorized bool     `json:"preauthorized"`
				}{
					Reusable:      args.Reusable,
					Ephemeral:     args.Ephemeral,
					Tags:          args.Tags,
					Preauthorized: args.Preauthorized,
				},
			},
		},
		Description: args.Description,
	}

	if args.ExpirySeconds > 0 {
		expiry := int64(args.ExpirySeconds)
		createReq.ExpirySeconds = expiry
	}

	client := kt.client.GetClient()
	key, err := client.Keys().Create(ctx, createReq)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create key: %v", err)), nil
	}

	keyJSON, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal key: %v", err)), nil
	}

	return mcp.NewToolResultText(string(keyJSON)), nil
}

func (kt *KeyTools) DeleteKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		KeyID string `json:"key_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := kt.client.GetClient()
	if err := client.Keys().Delete(ctx, args.KeyID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete key: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Key %s deleted successfully", args.KeyID)), nil
}
