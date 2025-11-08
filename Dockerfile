# GoReleaser-compatible Dockerfile
# Uses distroless for minimal attack surface (~20MB total)
# Note: GoReleaser provides pre-built binary in build context

FROM gcr.io/distroless/static:nonroot

# Copy the pre-built binary from GoReleaser build context
COPY jsmon /usr/local/bin/jsmon

# Set working directory
WORKDIR /app

# Run as non-root user (uid/gid 65532)
USER 65532:65532

# OCI labels (some will be overridden by GoReleaser)
LABEL org.opencontainers.image.title="JSMon-Go"
LABEL org.opencontainers.image.description="JavaScript Change Monitor for Bug Bounty Hunting"
LABEL org.opencontainers.image.authors="LuD1161"
LABEL org.opencontainers.image.licenses="MIT"

# Runtime expects these directories to be mounted as volumes:
# - /app/targets    - Target files with URLs to monitor
# - /app/downloads  - Downloaded content storage
# - /app/jsmon.json - State file (endpoint→hash mapping)
# - /app/.env       - Configuration file

ENTRYPOINT ["/usr/local/bin/jsmon"]
