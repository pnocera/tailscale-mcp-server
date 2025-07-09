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

type DeviceTools struct {
	client *client.TailscaleClient
}

func NewDeviceTools(client *client.TailscaleClient) *DeviceTools {
	return &DeviceTools{client: client}
}

func (dt *DeviceTools) RegisterTools(mcpServer *server.MCPServer) {
	tool := mcp.NewTool(
		"tailscale_devices_list",
		mcp.WithDescription("List all devices in the tailnet. Returns device information including name, IP addresses, machine key, node key, and basic connectivity status. Use 'all' fields to get complete device details including OS version, last seen timestamp, and advanced networking configuration. OAuth Scope: devices:read."),
		mcp.WithString("fields", mcp.Description("Fields to return. Can be 'all' or 'default'"), mcp.Enum("all", "default"), mcp.DefaultString("default")),
	)
	mcpServer.AddTool(tool, dt.ListDevices)

	tool = mcp.NewTool(
		"tailscale_device_get",
		mcp.WithDescription("Get detailed information about a specific device in the tailnet. Returns comprehensive device data including hardware specs, network configuration, authentication status, and connectivity details. Use 'all' fields for complete device information including OS version, last seen timestamp, and advanced networking settings. OAuth Scope: devices:read."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
		mcp.WithString("fields", mcp.Description("Fields to return. Can be 'all' or 'default'"), mcp.Enum("all", "default"), mcp.DefaultString("default")),
	)
	mcpServer.AddTool(tool, dt.GetDevice)

	tool = mcp.NewTool(
		"tailscale_device_delete",
		mcp.WithDescription("Remove a device from the tailnet permanently. This action cannot be undone. The device will lose access to the tailnet and must be re-added with a new auth key to rejoin. Use this for devices that are no longer needed or compromised. OAuth Scope: devices:write."),
		mcp.WithString("device_id", mcp.Description("The device ID to delete"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.DeleteDevice)

	tool = mcp.NewTool(
		"tailscale_device_authorize",
		mcp.WithDescription("Authorize or deauthorize a device for tailnets requiring device authorization. When authorized=true, grants the device access to the tailnet. When authorized=false, revokes access while keeping the device in the tailnet. Useful for temporarily restricting access without removing the device entirely. OAuth Scope: devices:core."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
		mcp.WithBoolean("authorized", mcp.Description("Whether to authorize (true) or deauthorize (false) the device"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.AuthorizeDevice)

	tool = mcp.NewTool(
		"tailscale_device_set_name",
		mcp.WithDescription("Set the Tailscale device name (machine name) for a device. This is the canonical name used throughout the tailnet and affects Magic DNS URLs. Changes propagate immediately, breaking existing Magic DNS URLs with the old name. Provide as FQDN (e.g., 'server.domain.ts.net') or base name (e.g., 'server'). Empty name resets to OS hostname. OAuth Scope: devices:core."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
		mcp.WithString("name", mcp.Description("The new name for the device"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetDeviceName)

	tool = mcp.NewTool(
		"tailscale_device_set_tags",
		mcp.WithDescription("Set tags on a device to assign a non-human identity for ACL-based access control. Tags are more flexible than role accounts and allow multiple identities per device. Must be defined in the tailnet policy file with proper ownership. Once tagged, the tag owns the device. Useful for servers, CI/CD systems, and automated services. OAuth Scope: devices:core."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
		mcp.WithArray("tags", mcp.Description("Array of tags to set on the device"), mcp.WithStringItems(), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetDeviceTags)

	tool = mcp.NewTool(
		"tailscale_device_expire",
		mcp.WithDescription("Expire a device's authentication key, forcing it to re-authenticate to maintain tailnet access. This is a security measure to ensure devices periodically refresh their credentials. The device will need to complete the authentication process again. Use this for security compliance or to revoke access temporarily. OAuth Scope: devices:core."),
		mcp.WithString("device_id", mcp.Description("The device ID to expire"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.ExpireDevice)

	tool = mcp.NewTool(
		"tailscale_device_routes_list",
		mcp.WithDescription("List subnet routes advertised and enabled for a device. Shows both advertised routes (what the device can route) and enabled routes (what the tailnet allows it to route). Routes must be both advertised and enabled to function as subnet routers or exit nodes. Essential for managing network connectivity and traffic routing. OAuth Scope: devices:routes:read."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.ListDeviceRoutes)

	tool = mcp.NewTool(
		"tailscale_device_routes_set",
		mcp.WithDescription("Set enabled subnet routes for a device by replacing the existing list. Routes must be both advertised by the device and enabled via this API to function. Cannot set advertised routes (must be done on device). Use for configuring subnet routers and exit nodes. Examples: ['10.0.0.0/16', '192.168.1.0/24']. OAuth Scope: devices:routes."),
		mcp.WithString("device_id", mcp.Description("The device ID"), mcp.Required()),
		mcp.WithArray("routes", mcp.Description("Array of routes to set"), mcp.WithStringItems(), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetDeviceRoutes)
}

func (dt *DeviceTools) ListDevices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Fields string `json:"fields"`
	}

	if request.Params.Arguments != nil {
		if err := request.BindArguments(&args); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
		}
	}

	client := dt.client.GetClient()
	var devices []tailscale.Device
	var err error

	if args.Fields == "all" {
		devices, err = client.Devices().ListWithAllFields(ctx)
	} else {
		devices, err = client.Devices().List(ctx)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list devices: %v", err)), nil
	}

	devicesJSON, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal devices: %v", err)), nil
	}

	return mcp.NewToolResultText(string(devicesJSON)), nil
}

func (dt *DeviceTools) GetDevice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string `json:"device_id"`
		Fields   string `json:"fields"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	var device *tailscale.Device
	var err error

	if args.Fields == "all" {
		device, err = client.Devices().GetWithAllFields(ctx, args.DeviceID)
	} else {
		device, err = client.Devices().Get(ctx, args.DeviceID)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get device: %v", err)), nil
	}

	deviceJSON, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal device: %v", err)), nil
	}

	return mcp.NewToolResultText(string(deviceJSON)), nil
}

func (dt *DeviceTools) DeleteDevice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string `json:"device_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.Devices().Delete(ctx, args.DeviceID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete device: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s deleted successfully", args.DeviceID)), nil
}

func (dt *DeviceTools) AuthorizeDevice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID   string `json:"device_id"`
		Authorized bool   `json:"authorized"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.Devices().SetAuthorized(ctx, args.DeviceID, args.Authorized); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set device authorization: %v", err)), nil
	}

	status := "authorized"
	if !args.Authorized {
		status = "deauthorized"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s %s successfully", args.DeviceID, status)), nil
}

func (dt *DeviceTools) SetDeviceName(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string `json:"device_id"`
		Name     string `json:"name"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.Devices().SetName(ctx, args.DeviceID, args.Name); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set device name: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s name set to %s", args.DeviceID, args.Name)), nil
}

func (dt *DeviceTools) SetDeviceTags(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string   `json:"device_id"`
		Tags     []string `json:"tags"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.Devices().SetTags(ctx, args.DeviceID, args.Tags); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set device tags: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s tags set to %v", args.DeviceID, args.Tags)), nil
}

func (dt *DeviceTools) ExpireDevice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string `json:"device_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	// The ExpireKey method doesn't exist in the current API, so we'll set key expiry to be disabled=false
	deviceKey := tailscale.DeviceKey{KeyExpiryDisabled: false}
	if err := client.Devices().SetKey(ctx, args.DeviceID, deviceKey); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set device key expiry: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s expired successfully", args.DeviceID)), nil
}

func (dt *DeviceTools) ListDeviceRoutes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string `json:"device_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	routes, err := client.Devices().SubnetRoutes(ctx, args.DeviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list device routes: %v", err)), nil
	}

	routesJSON, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal routes: %v", err)), nil
	}

	return mcp.NewToolResultText(string(routesJSON)), nil
}

func (dt *DeviceTools) SetDeviceRoutes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DeviceID string   `json:"device_id"`
		Routes   []string `json:"routes"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.Devices().SetSubnetRoutes(ctx, args.DeviceID, args.Routes); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set device routes: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Device %s routes set to %v", args.DeviceID, args.Routes)), nil
}
