# Multi-stage build for minimal image size

# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for multi-architecture support
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Build the application with version information
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o nexttrace_exporter \
    .

# Stage 2: Install nexttrace
FROM alpine:latest AS nexttrace-installer

# Install curl for downloading nexttrace
RUN apk add --no-cache curl bash

# Download and install nexttrace
RUN curl -sL nxtrace.org/nt | bash

# Stage 3: Final runtime image
FROM alpine:latest

# Add labels for better container management
LABEL maintainer="vinsec" \
      description="Prometheus exporter for NextTrace network path tracing" \
      org.opencontainers.image.source="https://github.com/vinsec/nexttrace_exporter" \
      org.opencontainers.image.licenses="MIT"

# Install runtime dependencies (wget needed for healthcheck)
RUN apk add --no-cache ca-certificates tzdata wget

# Copy nexttrace binary from installer stage
COPY --from=nexttrace-installer /usr/local/bin/nexttrace /usr/local/bin/nexttrace

# Copy exporter binary from builder stage
COPY --from=builder /build/nexttrace_exporter /usr/local/bin/nexttrace_exporter

# Create config directory
RUN mkdir -p /etc/nexttrace_exporter

# Note: Running as root is required for nexttrace to have raw network socket access
# This is necessary for ICMP operations used in traceroute functionality

WORKDIR /root

# Expose metrics port
EXPOSE 9101

# Set default command
ENTRYPOINT ["/usr/local/bin/nexttrace_exporter"]
CMD ["--config.file=/etc/nexttrace_exporter/config.yml"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9101/-/healthy || exit 1
