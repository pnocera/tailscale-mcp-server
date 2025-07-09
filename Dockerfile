# Multi-stage build for optimal image size
FROM golang:1.24-alpine AS builder

# Install required packages for building
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (handle local replace directives)
RUN go mod edit -dropreplace tailscale.com/client/tailscale/v2 && \
    go mod download && \
    go mod edit -replace tailscale.com/client/tailscale/v2=../

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o tailscale-mcp-server \
    ./cmd

# Final stage - minimal runtime image
FROM scratch

# Copy ca-certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/tailscale-mcp-server /tailscale-mcp-server

# Expose the default MCP port (stdio doesn't need ports, but useful for HTTP mode)
EXPOSE 8080

# Set default environment variables
ENV TAILSCALE_TAILNET="-"

# Health check to ensure the server is responsive
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/tailscale-mcp-server", "--health-check"] || exit 1

# Run as non-root user for security
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/tailscale-mcp-server"]