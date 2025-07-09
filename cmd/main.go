package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pnocera/tailscale-mcp-server/internal/client"
	"github.com/pnocera/tailscale-mcp-server/internal/config"
	"github.com/pnocera/tailscale-mcp-server/internal/handlers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	tailscaleClient, err := client.NewTailscaleClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Tailscale client: %v", err)
	}

	if err := tailscaleClient.ValidateConnection(context.Background()); err != nil {
		log.Fatalf("Failed to validate Tailscale connection: %v", err)
	}

	mcpServer := server.NewMCPServer(
		"tailscale-mcp-server",
		"1.0.0",
		server.WithLogging(),
	)

	handler := handlers.NewHandler(tailscaleClient)
	handler.RegisterTools(mcpServer)

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
