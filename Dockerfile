# Multi-stage Dockerfile for JSMon-Go
# Uses distroless for minimal attack surface and image size

# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0 for static binary
# -ldflags="-s -w" to strip debug info and reduce size
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w" \
    -o jsmon \
    ./cmd/jsmon

# Stage 2: Create minimal runtime image with distroless
FROM gcr.io/distroless/static:nonroot

# Import ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data for proper logging
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /build/jsmon /usr/local/bin/jsmon

# Set working directory
WORKDIR /app

# Create directories for runtime (will be mounted as volumes)
# Note: distroless runs as nonroot user (uid/gid 65532)
USER 65532:65532

# Labels for OCI image spec
LABEL org.opencontainers.image.title="JSMon-Go"
LABEL org.opencontainers.image.description="JavaScript Change Monitor for Bug Bounty Hunting"
LABEL org.opencontainers.image.authors="aseemshrey"
LABEL org.opencontainers.image.source="https://github.com/aseemshrey/jsmon-go"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.documentation="https://github.com/aseemshrey/jsmon-go#readme"

# Runtime expects these directories to be mounted:
# - /app/targets    - Target files
# - /app/downloads  - Downloaded content
# - /app/jsmon.json - State file
# - /app/.env       - Configuration

ENTRYPOINT ["/usr/local/bin/jsmon"]
