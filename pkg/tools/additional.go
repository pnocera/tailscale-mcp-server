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

type AdditionalTools struct {
	client *client.TailscaleClient
}

func NewAdditionalTools(client *client.TailscaleClient) *AdditionalTools {
	return &AdditionalTools{client: client}
}

func (at *AdditionalTools) RegisterTools(mcpServer *server.MCPServer) {
	// Webhook tools
	tool := mcp.NewTool(
		"tailscale_webhooks_list",
		mcp.WithDescription("List all webhook endpoints configured for the tailnet. Returns webhook endpoint URLs, subscription types, and status information. Use this to manage and monitor event notifications sent to external systems. OAuth Scope: webhooks:read."),
	)
	mcpServer.AddTool(tool, at.ListWebhooks)

	tool = mcp.NewTool(
		"tailscale_webhook_create",
		mcp.WithDescription("Create a new webhook endpoint to receive tailnet events. Configure the endpoint URL and specify which event types to subscribe to (e.g., device changes, user events). Essential for integrating Tailscale with external monitoring and automation systems. OAuth Scope: webhooks:write."),
		mcp.WithString("endpoint_url", mcp.Description("The URL where webhook events will be sent"), mcp.Required()),
		mcp.WithArray("subscriptions", mcp.Description("List of event types to subscribe to"), mcp.WithStringItems(), mcp.Required()),
	)
	mcpServer.AddTool(tool, at.CreateWebhook)

	tool = mcp.NewTool(
		"tailscale_webhook_get",
		mcp.WithDescription("Get detailed information about a specific webhook endpoint. Returns endpoint configuration, subscription types, delivery status, and webhook statistics. Use this to monitor webhook performance and troubleshoot delivery issues. OAuth Scope: webhooks:read."),
		mcp.WithString("endpoint_id", mcp.Description("The webhook endpoint ID"), mcp.Required()),
	)
	mcpServer.AddTool(tool, at.GetWebhook)

	tool = mcp.NewTool(
		"tailscale_webhook_delete",
		mcp.WithDescription("Delete a webhook endpoint permanently. This stops all event notifications to the specified endpoint. Use this to remove unused or misconfigured webhooks. Essential for maintaining clean webhook configurations. OAuth Scope: webhooks:write."),
		mcp.WithString("endpoint_id", mcp.Description("The webhook endpoint ID to delete"), mcp.Required()),
	)
	mcpServer.AddTool(tool, at.DeleteWebhook)

	// Logging tools
	tool = mcp.NewTool(
		"tailscale_logging_configuration_get",
		mcp.WithDescription("Get configuration audit logs for the tailnet. Returns log streaming configuration for administrative and policy changes. Essential for compliance, security auditing, and troubleshooting configuration issues. Learn more about logging at /kb/1349/log-events. OAuth Scope: logging:read."),
	)
	mcpServer.AddTool(tool, at.GetConfigurationLogs)

	tool = mcp.NewTool(
		"tailscale_logging_network_get",
		mcp.WithDescription("Get network flow logs for the tailnet. Returns log streaming configuration for network traffic and connection data. Essential for network monitoring, security analysis, and troubleshooting connectivity issues. Learn more about logging at /kb/1349/log-events. OAuth Scope: logging:read."),
	)
	mcpServer.AddTool(tool, at.GetNetworkLogs)

	// Device posture tools
	tool = mcp.NewTool(
		"tailscale_device_posture_integrations_list",
		mcp.WithDescription("List device posture integrations configured for the tailnet. Returns integrations with device posture data providers like CrowdStrike, Microsoft Intune, and others. Essential for managing device security compliance and conditional access policies. Learn more about device posture at /kb/1288/device-posture. OAuth Scope: posture:read."),
	)
	mcpServer.AddTool(tool, at.ListPostureIntegrations)

	tool = mcp.NewTool(
		"tailscale_device_posture_integration_create",
		mcp.WithDescription("Create a new device posture integration with security providers like CrowdStrike, Microsoft Intune, or others. Configure OAuth credentials and provider-specific settings to enable device security data collection. Essential for implementing zero-trust security policies based on device compliance. OAuth Scope: posture:write."),
		mcp.WithString("provider", mcp.Description("The posture provider (e.g., 'crowdstrike', 'intune')"), mcp.Required()),
		mcp.WithString("client_id", mcp.Description("OAuth client ID for the integration"), mcp.Required()),
		mcp.WithString("client_secret", mcp.Description("OAuth client secret for the integration"), mcp.Required()),
		mcp.WithString("tenant_id", mcp.Description("Tenant ID (required for some providers)")),
	)
	mcpServer.AddTool(tool, at.CreatePostureIntegration)

	tool = mcp.NewTool(
		"tailscale_device_posture_integration_get",
		mcp.WithDescription("Get detailed information about a specific device posture integration. Returns integration configuration, connection status, and data collection statistics. Use this to monitor integration health and troubleshoot device posture data issues. OAuth Scope: posture:read."),
		mcp.WithString("id", mcp.Description("The integration ID"), mcp.Required()),
	)
	mcpServer.AddTool(tool, at.GetPostureIntegration)

	tool = mcp.NewTool(
		"tailscale_device_posture_integration_delete",
		mcp.WithDescription("Delete a device posture integration permanently. This stops device security data collection from the specified provider. Use this to remove unused or misconfigured integrations. Note that this may affect security policies that depend on posture data. OAuth Scope: posture:write."),
		mcp.WithString("id", mcp.Description("The integration ID to delete"), mcp.Required()),
	)
	mcpServer.AddTool(tool, at.DeletePostureIntegration)

	// Tailnet settings tools
	tool = mcp.NewTool(
		"tailscale_tailnet_settings_get",
		mcp.WithDescription("Get tailnet settings and configuration. Returns device approval settings, user permissions, key duration, logging preferences, routing options, and posture collection settings. Essential for understanding and managing tailnet policies and behavior. OAuth Scope: settings:read."),
	)
	mcpServer.AddTool(tool, at.GetTailnetSettings)

	tool = mcp.NewTool(
		"tailscale_tailnet_settings_update",
		mcp.WithDescription("Update tailnet settings and configuration. Configure device approval requirements, automatic updates, key durations, user permissions, network logging, regional routing, and posture data collection. Changes affect all devices and users in the tailnet. Use with caution as settings impact security and connectivity. OAuth Scope: settings:write."),
		mcp.WithBoolean("devices_approval_on", mcp.Description("Whether device approval is required")),
		mcp.WithBoolean("devices_auto_updates_on", mcp.Description("Whether devices should auto-update")),
		mcp.WithNumber("devices_key_duration_days", mcp.Description("Default key duration in days")),
		mcp.WithBoolean("users_approval_on", mcp.Description("Whether user approval is required")),
		mcp.WithString("users_role_allowed_to_join_external_tailnets", mcp.Description("Role allowed to join external tailnets")),
		mcp.WithBoolean("network_flow_logging_on", mcp.Description("Whether network flow logging is enabled")),
		mcp.WithBoolean("regional_routing_on", mcp.Description("Whether regional routing is enabled")),
		mcp.WithBoolean("posture_identity_collection_on", mcp.Description("Whether posture identity collection is enabled")),
	)
	mcpServer.AddTool(tool, at.UpdateTailnetSettings)
}

func (at *AdditionalTools) ListWebhooks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := at.client.GetClient()
	webhooks, err := client.Webhooks().List(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list webhooks: %v", err)), nil
	}

	webhooksJSON, err := json.MarshalIndent(webhooks, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal webhooks: %v", err)), nil
	}

	return mcp.NewToolResultText(string(webhooksJSON)), nil
}

func (at *AdditionalTools) CreateWebhook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		EndpointURL   string   `json:"endpoint_url"`
		Subscriptions []string `json:"subscriptions"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// Convert string subscriptions to WebhookSubscriptionType
	subscriptions := make([]tailscale.WebhookSubscriptionType, len(args.Subscriptions))
	for i, sub := range args.Subscriptions {
		subscriptions[i] = tailscale.WebhookSubscriptionType(sub)
	}

	createReq := tailscale.CreateWebhookRequest{
		EndpointURL:   args.EndpointURL,
		Subscriptions: subscriptions,
	}

	client := at.client.GetClient()
	webhook, err := client.Webhooks().Create(ctx, createReq)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create webhook: %v", err)), nil
	}

	webhookJSON, err := json.MarshalIndent(webhook, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal webhook: %v", err)), nil
	}

	return mcp.NewToolResultText(string(webhookJSON)), nil
}

func (at *AdditionalTools) GetWebhook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		EndpointID string `json:"endpoint_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := at.client.GetClient()
	webhook, err := client.Webhooks().Get(ctx, args.EndpointID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get webhook: %v", err)), nil
	}

	webhookJSON, err := json.MarshalIndent(webhook, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal webhook: %v", err)), nil
	}

	return mcp.NewToolResultText(string(webhookJSON)), nil
}

func (at *AdditionalTools) DeleteWebhook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		EndpointID string `json:"endpoint_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := at.client.GetClient()
	if err := client.Webhooks().Delete(ctx, args.EndpointID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete webhook: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Webhook %s deleted successfully", args.EndpointID)), nil
}

func (at *AdditionalTools) GetConfigurationLogs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := at.client.GetClient()
	logs, err := client.Logging().LogstreamConfiguration(ctx, tailscale.LogTypeConfig)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get configuration logs: %v", err)), nil
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal logs: %v", err)), nil
	}

	return mcp.NewToolResultText(string(logsJSON)), nil
}

func (at *AdditionalTools) GetNetworkLogs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := at.client.GetClient()
	logs, err := client.Logging().LogstreamConfiguration(ctx, tailscale.LogTypeNetwork)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get network logs: %v", err)), nil
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal logs: %v", err)), nil
	}

	return mcp.NewToolResultText(string(logsJSON)), nil
}

func (at *AdditionalTools) ListPostureIntegrations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := at.client.GetClient()
	integrations, err := client.DevicePosture().ListIntegrations(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list posture integrations: %v", err)), nil
	}

	integrationsJSON, err := json.MarshalIndent(integrations, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal integrations: %v", err)), nil
	}

	return mcp.NewToolResultText(string(integrationsJSON)), nil
}

func (at *AdditionalTools) CreatePostureIntegration(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Provider     string `json:"provider"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		TenantID     string `json:"tenant_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	createReq := tailscale.CreatePostureIntegrationRequest{
		Provider:     tailscale.PostureIntegrationProvider(args.Provider),
		ClientID:     args.ClientID,
		ClientSecret: args.ClientSecret,
		TenantID:     args.TenantID,
	}

	client := at.client.GetClient()
	integration, err := client.DevicePosture().CreateIntegration(ctx, createReq)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create posture integration: %v", err)), nil
	}

	integrationJSON, err := json.MarshalIndent(integration, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal integration: %v", err)), nil
	}

	return mcp.NewToolResultText(string(integrationJSON)), nil
}

func (at *AdditionalTools) GetPostureIntegration(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := at.client.GetClient()
	integration, err := client.DevicePosture().GetIntegration(ctx, args.ID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get posture integration: %v", err)), nil
	}

	integrationJSON, err := json.MarshalIndent(integration, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal integration: %v", err)), nil
	}

	return mcp.NewToolResultText(string(integrationJSON)), nil
}

func (at *AdditionalTools) DeletePostureIntegration(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := at.client.GetClient()
	if err := client.DevicePosture().DeleteIntegration(ctx, args.ID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete posture integration: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Posture integration %s deleted successfully", args.ID)), nil
}

func (at *AdditionalTools) GetTailnetSettings(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := at.client.GetClient()
	settings, err := client.TailnetSettings().Get(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tailnet settings: %v", err)), nil
	}

	settingsJSON, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal settings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(settingsJSON)), nil
}

func (at *AdditionalTools) UpdateTailnetSettings(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		DevicesApprovalOn                      *bool   `json:"devices_approval_on"`
		DevicesAutoUpdatesOn                   *bool   `json:"devices_auto_updates_on"`
		DevicesKeyDurationDays                 *int    `json:"devices_key_duration_days"`
		UsersApprovalOn                        *bool   `json:"users_approval_on"`
		UsersRoleAllowedToJoinExternalTailnets *string `json:"users_role_allowed_to_join_external_tailnets"`
		NetworkFlowLoggingOn                   *bool   `json:"network_flow_logging_on"`
		RegionalRoutingOn                      *bool   `json:"regional_routing_on"`
		PostureIdentityCollectionOn            *bool   `json:"posture_identity_collection_on"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	updateReq := tailscale.UpdateTailnetSettingsRequest{}
	if args.DevicesApprovalOn != nil {
		updateReq.DevicesApprovalOn = args.DevicesApprovalOn
	}
	if args.DevicesAutoUpdatesOn != nil {
		updateReq.DevicesAutoUpdatesOn = args.DevicesAutoUpdatesOn
	}
	if args.DevicesKeyDurationDays != nil {
		updateReq.DevicesKeyDurationDays = args.DevicesKeyDurationDays
	}
	if args.UsersApprovalOn != nil {
		updateReq.UsersApprovalOn = args.UsersApprovalOn
	}
	if args.UsersRoleAllowedToJoinExternalTailnets != nil {
		role := tailscale.RoleAllowedToJoinExternalTailnets(*args.UsersRoleAllowedToJoinExternalTailnets)
		updateReq.UsersRoleAllowedToJoinExternalTailnets = &role
	}
	if args.NetworkFlowLoggingOn != nil {
		updateReq.NetworkFlowLoggingOn = args.NetworkFlowLoggingOn
	}
	if args.RegionalRoutingOn != nil {
		updateReq.RegionalRoutingOn = args.RegionalRoutingOn
	}
	if args.PostureIdentityCollectionOn != nil {
		updateReq.PostureIdentityCollectionOn = args.PostureIdentityCollectionOn
	}

	client := at.client.GetClient()
	if err := client.TailnetSettings().Update(ctx, updateReq); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update tailnet settings: %v", err)), nil
	}

	// Get the updated settings to return
	settings, err := client.TailnetSettings().Get(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get updated tailnet settings: %v", err)), nil
	}

	settingsJSON, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal settings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(settingsJSON)), nil
}
