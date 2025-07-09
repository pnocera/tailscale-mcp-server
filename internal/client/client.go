package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/pnocera/tailscale-mcp-server/internal/config"
	"tailscale.com/client/tailscale/v2"
)

type TailscaleClient struct {
	client *tailscale.Client
	mu     sync.RWMutex
}

func NewTailscaleClient(cfg *config.Config) (*TailscaleClient, error) {
	client := &tailscale.Client{
		Tailnet: cfg.TailscaleTailnet,
	}

	if cfg.UseOAuth {
		oauthConfig := tailscale.OAuthConfig{
			ClientID:     cfg.TailscaleClientID,
			ClientSecret: cfg.TailscaleClientSecret,
			Scopes:       []string{"all:read", "all:write"},
		}
		client.HTTP = oauthConfig.HTTPClient()
	} else {
		client.APIKey = cfg.TailscaleAPIKey
	}

	return &TailscaleClient{
		client: client,
	}, nil
}

func (tc *TailscaleClient) GetClient() *tailscale.Client {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.client
}

func (tc *TailscaleClient) ValidateConnection(ctx context.Context) error {
	client := tc.GetClient()
	_, err := client.Devices().List(ctx)
	if err != nil {
		return fmt.Errorf("failed to validate Tailscale connection: %w", err)
	}
	return nil
}
