package handlers

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/pnocera/tailscale-mcp-server/internal/client"
	"github.com/pnocera/tailscale-mcp-server/pkg/tools"
)

type Handler struct {
	client *client.TailscaleClient
}

func NewHandler(client *client.TailscaleClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) RegisterTools(mcpServer *server.MCPServer) {
	deviceTools := tools.NewDeviceTools(h.client)
	deviceTools.RegisterTools(mcpServer)

	keyTools := tools.NewKeyTools(h.client)
	keyTools.RegisterTools(mcpServer)

	userTools := tools.NewUserTools(h.client)
	userTools.RegisterTools(mcpServer)

	dnsTools := tools.NewDNSTools(h.client)
	dnsTools.RegisterTools(mcpServer)

	additionalTools := tools.NewAdditionalTools(h.client)
	additionalTools.RegisterTools(mcpServer)
}
