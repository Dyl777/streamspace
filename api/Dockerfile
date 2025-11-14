# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.mod
COPY go.sum go.sum*

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o api-server cmd/main.go

# Final stage - minimal runtime image
FROM alpine:3.19

# Install ca-certificates for HTTPS and git for repository cloning
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copy binary from builder
COPY --from=builder /workspace/api-server .

# Create directory for repository clones
RUN mkdir -p /tmp/streamspace-repos && chmod 777 /tmp/streamspace-repos

# Run as non-root user
RUN addgroup -g 1000 streamspace && \
    adduser -D -u 1000 -G streamspace streamspace && \
    chown -R streamspace:streamspace /app /tmp/streamspace-repos

USER streamspace

# Expose API port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

ENTRYPOINT ["/app/api-server"]
