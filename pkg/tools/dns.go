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

type DNSTools struct {
	client *client.TailscaleClient
}

func NewDNSTools(client *client.TailscaleClient) *DNSTools {
	return &DNSTools{client: client}
}

func (dt *DNSTools) RegisterTools(mcpServer *server.MCPServer) {
	tool := mcp.NewTool(
		"tailscale_dns_nameservers_get",
		mcp.WithDescription("Get DNS nameservers configured for the tailnet. Returns the list of DNS servers that devices will use for domain resolution. Essential for understanding and troubleshooting DNS configuration. Learn more about DNS in Tailscale at /kb/1054/dns. OAuth Scope: dns:read."),
	)
	mcpServer.AddTool(tool, dt.GetNameservers)

	tool = mcp.NewTool(
		"tailscale_dns_nameservers_set",
		mcp.WithDescription("Set DNS nameservers for the tailnet. Configure which DNS servers devices will use for domain resolution. Provide IP addresses of DNS servers (e.g., ['8.8.8.8', '1.1.1.1']). Changes apply to all devices in the tailnet. Learn more about DNS in Tailscale at /kb/1054/dns. OAuth Scope: dns:write."),
		mcp.WithArray("nameservers", mcp.Description("List of DNS nameserver addresses"), mcp.WithStringItems(), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetNameservers)

	tool = mcp.NewTool(
		"tailscale_dns_preferences_get",
		mcp.WithDescription("Get DNS preferences for the tailnet. Returns MagicDNS configuration and other DNS settings. MagicDNS enables automatic DNS resolution for device names within the tailnet (e.g., 'device-name.tailnet.ts.net'). Essential for understanding DNS behavior. OAuth Scope: dns:read."),
	)
	mcpServer.AddTool(tool, dt.GetPreferences)

	tool = mcp.NewTool(
		"tailscale_dns_preferences_set",
		mcp.WithDescription("Set DNS preferences for the tailnet. Enable or disable MagicDNS, which provides automatic DNS resolution for device names within the tailnet. When enabled, devices can reach each other using names like 'device-name.tailnet.ts.net'. Essential for easy device connectivity. OAuth Scope: dns:write."),
		mcp.WithBoolean("magic_dns", mcp.Description("Enable MagicDNS"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetPreferences)

	tool = mcp.NewTool(
		"tailscale_dns_searchpaths_get",
		mcp.WithDescription("Get DNS search paths for the tailnet. Returns the list of domain suffixes that will be appended to short hostnames during DNS resolution. For example, with search path 'company.com', 'server' resolves to 'server.company.com'. Essential for understanding DNS resolution behavior. OAuth Scope: dns:read."),
	)
	mcpServer.AddTool(tool, dt.GetSearchPaths)

	tool = mcp.NewTool(
		"tailscale_dns_searchpaths_set",
		mcp.WithDescription("Set DNS search paths for the tailnet. Configure domain suffixes that will be appended to short hostnames during DNS resolution. For example, with search path 'company.com', typing 'server' will resolve to 'server.company.com'. Improves user experience by enabling short hostname usage. OAuth Scope: dns:write."),
		mcp.WithArray("search_paths", mcp.Description("List of DNS search paths"), mcp.WithStringItems(), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetSearchPaths)

	tool = mcp.NewTool(
		"tailscale_policy_get",
		mcp.WithDescription("Get the current policy file (ACL) for the tailnet. Returns the access control list in HuJSON format that defines who can access what resources. The policy file controls device access, user permissions, and network routing rules. Essential for understanding and managing security policies. Learn more about ACLs at /kb/1018/acls. OAuth Scope: acl:read."),
	)
	mcpServer.AddTool(tool, dt.GetPolicy)

	tool = mcp.NewTool(
		"tailscale_policy_set",
		mcp.WithDescription("Set the policy file (ACL) for the tailnet. Upload a new access control list in HuJSON format to define security policies. Controls device access, user permissions, SSH access, and network routing. Changes apply immediately to all devices. Validate policy first using tailscale_policy_validate. Learn more about ACLs at /kb/1018/acls. OAuth Scope: acl:write."),
		mcp.WithString("policy", mcp.Description("Policy file content in HuJSON format"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.SetPolicy)

	tool = mcp.NewTool(
		"tailscale_policy_validate",
		mcp.WithDescription("Validate a policy file (ACL) without applying it to the tailnet. Checks the HuJSON syntax and policy rules for errors before deployment. Essential for safe policy management - always validate before setting a new policy. Prevents accidental misconfigurations that could disrupt network access. Learn more about ACLs at /kb/1018/acls. OAuth Scope: acl:read."),
		mcp.WithString("policy", mcp.Description("Policy file content in HuJSON format to validate"), mcp.Required()),
	)
	mcpServer.AddTool(tool, dt.ValidatePolicy)
}

func (dt *DNSTools) GetNameservers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := dt.client.GetClient()
	nameservers, err := client.DNS().Nameservers(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get nameservers: %v", err)), nil
	}

	nameserversJSON, err := json.MarshalIndent(nameservers, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal nameservers: %v", err)), nil
	}

	return mcp.NewToolResultText(string(nameserversJSON)), nil
}

func (dt *DNSTools) SetNameservers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Nameservers []string `json:"nameservers"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.DNS().SetNameservers(ctx, args.Nameservers); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set nameservers: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("DNS nameservers set to: %v", args.Nameservers)), nil
}

func (dt *DNSTools) GetPreferences(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := dt.client.GetClient()
	preferences, err := client.DNS().Preferences(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get DNS preferences: %v", err)), nil
	}

	preferencesJSON, err := json.MarshalIndent(preferences, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal preferences: %v", err)), nil
	}

	return mcp.NewToolResultText(string(preferencesJSON)), nil
}

func (dt *DNSTools) SetPreferences(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		MagicDNS bool `json:"magic_dns"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	preferences := tailscale.DNSPreferences{
		MagicDNS: args.MagicDNS,
	}

	client := dt.client.GetClient()
	if err := client.DNS().SetPreferences(ctx, preferences); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set DNS preferences: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("DNS preferences updated: MagicDNS=%v", args.MagicDNS)), nil
}

func (dt *DNSTools) GetSearchPaths(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := dt.client.GetClient()
	searchPaths, err := client.DNS().SearchPaths(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get search paths: %v", err)), nil
	}

	searchPathsJSON, err := json.MarshalIndent(searchPaths, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal search paths: %v", err)), nil
	}

	return mcp.NewToolResultText(string(searchPathsJSON)), nil
}

func (dt *DNSTools) SetSearchPaths(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		SearchPaths []string `json:"search_paths"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.DNS().SetSearchPaths(ctx, args.SearchPaths); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set search paths: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("DNS search paths set to: %v", args.SearchPaths)), nil
}

func (dt *DNSTools) GetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := dt.client.GetClient()
	policy, err := client.PolicyFile().Raw(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get policy: %v", err)), nil
	}

	return mcp.NewToolResultText(policy.HuJSON), nil
}

func (dt *DNSTools) SetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Policy string `json:"policy"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.PolicyFile().Set(ctx, args.Policy, ""); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set policy: %v", err)), nil
	}

	return mcp.NewToolResultText("Policy file updated successfully"), nil
}

func (dt *DNSTools) ValidatePolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Policy string `json:"policy"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := dt.client.GetClient()
	if err := client.PolicyFile().Validate(ctx, args.Policy); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Policy validation failed: %v", err)), nil
	}

	return mcp.NewToolResultText("Policy validation passed"), nil
}
