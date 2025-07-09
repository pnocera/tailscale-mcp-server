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

type UserTools struct {
	client *client.TailscaleClient
}

func NewUserTools(client *client.TailscaleClient) *UserTools {
	return &UserTools{client: client}
}

func (ut *UserTools) RegisterTools(mcpServer *server.MCPServer) {
	tool := mcp.NewTool(
		"tailscale_users_list",
		mcp.WithDescription("List all users in the tailnet. Returns user information including display name, login name, profile picture, role, status, and last seen timestamp. Essential for user management and access auditing. OAuth Scope: users:read."),
	)
	mcpServer.AddTool(tool, ut.ListUsers)

	tool = mcp.NewTool(
		"tailscale_user_get",
		mcp.WithDescription("Get detailed information about a specific user in the tailnet. Returns comprehensive user data including account details, role assignments, device count, and authentication status. Use this for user profile management and access verification. OAuth Scope: users:read."),
		mcp.WithString("user_id", mcp.Description("The user ID"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.GetUser)

	tool = mcp.NewTool(
		"tailscale_user_approve",
		mcp.WithDescription("Approve a user for tailnet access. This grants the user permission to join the tailnet and access resources according to their role and ACL policies. Use this for tailnets requiring user approval for new members. Note: This functionality may not be available in all API versions. OAuth Scope: users:write."),
		mcp.WithString("user_id", mcp.Description("The user ID to approve"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.ApproveUser)

	tool = mcp.NewTool(
		"tailscale_user_suspend",
		mcp.WithDescription("Suspend a user to temporarily revoke their tailnet access. Suspended users cannot access tailnet resources but remain in the user list for future restoration. Use this for temporary access control without removing the user permanently. Note: This functionality may not be available in all API versions. OAuth Scope: users:write."),
		mcp.WithString("user_id", mcp.Description("The user ID to suspend"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.SuspendUser)

	tool = mcp.NewTool(
		"tailscale_user_restore",
		mcp.WithDescription("Restore a previously suspended user to active status. This re-enables their access to tailnet resources according to their role and ACL policies. Use this to reinstate users after temporary suspension. Note: This functionality may not be available in all API versions. OAuth Scope: users:write."),
		mcp.WithString("user_id", mcp.Description("The user ID to restore"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.RestoreUser)

	tool = mcp.NewTool(
		"tailscale_user_delete",
		mcp.WithDescription("Delete a user from the tailnet permanently. This removes the user and their access to all tailnet resources. Use this for user offboarding or when users no longer need access. Note: This functionality may not be available in all API versions. OAuth Scope: users:write."),
		mcp.WithString("user_id", mcp.Description("The user ID to delete"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.DeleteUser)

	tool = mcp.NewTool(
		"tailscale_contacts_get",
		mcp.WithDescription("Get contact preferences for the tailnet. Returns configured contact information for account notifications, support requests, and security alerts. Essential for maintaining proper communication channels and compliance requirements. OAuth Scope: users:read."),
	)
	mcpServer.AddTool(tool, ut.GetContacts)

	tool = mcp.NewTool(
		"tailscale_contact_update",
		mcp.WithDescription("Update contact preferences for the tailnet. Configure email addresses for different contact types: 'account' for billing/administrative, 'support' for technical issues, and 'security' for security-related notifications. Essential for maintaining proper communication channels and compliance. OAuth Scope: users:write."),
		mcp.WithString("contact_type", mcp.Description("Type of contact (account, support, security)"), mcp.Enum("account", "support", "security"), mcp.Required()),
		mcp.WithString("email", mcp.Description("Email address for the contact"), mcp.Required()),
	)
	mcpServer.AddTool(tool, ut.UpdateContact)
}

func (ut *UserTools) ListUsers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := ut.client.GetClient()
	users, err := client.Users().List(ctx, nil, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list users: %v", err)), nil
	}

	usersJSON, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal users: %v", err)), nil
	}

	return mcp.NewToolResultText(string(usersJSON)), nil
}

func (ut *UserTools) GetUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		UserID string `json:"user_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	client := ut.client.GetClient()
	user, err := client.Users().Get(ctx, args.UserID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get user: %v", err)), nil
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}

	return mcp.NewToolResultText(string(userJSON)), nil
}

func (ut *UserTools) ApproveUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		UserID string `json:"user_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// User approval is not available in the current API
	return mcp.NewToolResultError("User approval functionality is not available in the current API"), nil
}

func (ut *UserTools) SuspendUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		UserID string `json:"user_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// User suspension is not available in the current API
	return mcp.NewToolResultError("User suspension functionality is not available in the current API"), nil
}

func (ut *UserTools) RestoreUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		UserID string `json:"user_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// User restoration is not available in the current API
	return mcp.NewToolResultError("User restoration functionality is not available in the current API"), nil
}

func (ut *UserTools) DeleteUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		UserID string `json:"user_id"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	// User deletion is not available in the current API
	return mcp.NewToolResultError("User deletion functionality is not available in the current API"), nil
}

func (ut *UserTools) GetContacts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := ut.client.GetClient()
	contacts, err := client.Contacts().Get(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get contacts: %v", err)), nil
	}

	contactsJSON, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal contacts: %v", err)), nil
	}

	return mcp.NewToolResultText(string(contactsJSON)), nil
}

func (ut *UserTools) UpdateContact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ContactType string `json:"contact_type"`
		Email       string `json:"email"`
	}

	if err := request.BindArguments(&args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	var contactType tailscale.ContactType
	switch args.ContactType {
	case "account":
		contactType = tailscale.ContactAccount
	case "support":
		contactType = tailscale.ContactSupport
	case "security":
		contactType = tailscale.ContactSecurity
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Invalid contact type: %s", args.ContactType)), nil
	}

	updateReq := tailscale.UpdateContactRequest{
		Email: &args.Email,
	}

	client := ut.client.GetClient()
	if err := client.Contacts().Update(ctx, contactType, updateReq); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update contact: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Contact %s updated to %s", args.ContactType, args.Email)), nil
}
