package config

import (
	"fmt"
	"os"
)

type Config struct {
	TailscaleAPIKey    string
	TailscaleTailnet   string
	TailscaleClientID  string
	TailscaleClientSecret string
	UseOAuth           bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		TailscaleAPIKey:       os.Getenv("TAILSCALE_API_KEY"),
		TailscaleTailnet:      os.Getenv("TAILSCALE_TAILNET"),
		TailscaleClientID:     os.Getenv("TAILSCALE_CLIENT_ID"),
		TailscaleClientSecret: os.Getenv("TAILSCALE_CLIENT_SECRET"),
	}

	if cfg.TailscaleTailnet == "" {
		cfg.TailscaleTailnet = "-"
	}

	cfg.UseOAuth = cfg.TailscaleClientID != "" && cfg.TailscaleClientSecret != ""

	if !cfg.UseOAuth && cfg.TailscaleAPIKey == "" {
		return nil, fmt.Errorf("either TAILSCALE_API_KEY or both TAILSCALE_CLIENT_ID and TAILSCALE_CLIENT_SECRET must be set")
	}

	return cfg, nil
}