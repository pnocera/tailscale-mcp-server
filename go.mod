module github.com/pnocera/tailscale-mcp-server

go 1.24.0

require (
	github.com/mark3labs/mcp-go v0.33.0
	tailscale.com/client/tailscale/v2 v2.0.0
)

replace tailscale.com/client/tailscale/v2 => ../

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/tailscale/hujson v0.0.0-20250605163823-992244df8c5a // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
)
