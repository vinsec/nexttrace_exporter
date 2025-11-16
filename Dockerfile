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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o nexttrace_exporter \
    .

# Stage 2: Install nexttrace
FROM alpine:latest AS nexttrace-installer

# Install curl for downloading nexttrace
RUN apk add --no-cache curl bash

# Download and install nexttrace
RUN curl -sSL https://raw.githubusercontent.com/sjlleo/nexttrace/main/install.sh | bash

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy nexttrace binary from installer stage
COPY --from=nexttrace-installer /usr/local/bin/nexttrace /usr/local/bin/nexttrace

# Copy exporter binary from builder stage
COPY --from=builder /build/nexttrace_exporter /usr/local/bin/nexttrace_exporter

# Create non-root user (note: exporter needs NET_RAW capability)
RUN addgroup -g 1000 nexttrace && \
    adduser -D -u 1000 -G nexttrace nexttrace

# Create config directory
RUN mkdir -p /etc/nexttrace_exporter && \
    chown -R nexttrace:nexttrace /etc/nexttrace_exporter

# Switch to non-root user
USER nexttrace

WORKDIR /home/nexttrace

# Expose metrics port
EXPOSE 9101

# Set default command
ENTRYPOINT ["/usr/local/bin/nexttrace_exporter"]
CMD ["--config.file=/etc/nexttrace_exporter/config.yml"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9101/-/healthy || exit 1
