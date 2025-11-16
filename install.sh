#!/bin/bash

# Installation script for nexttrace_exporter

set -e

VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
CONFIG_DIR="${CONFIG_DIR:-/etc/nexttrace_exporter}"
BINARY_NAME="nexttrace_exporter"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}NextTrace Exporter Installation Script${NC}"
echo "========================================"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Error: This script must be run as root${NC}" 
   exit 1
fi

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Detected OS: $OS"
echo "Detected Architecture: $ARCH"

# Check if nexttrace is installed
if ! command -v nexttrace &> /dev/null; then
    echo -e "${YELLOW}Warning: nexttrace is not installed${NC}"
    echo "Please install nexttrace first:"
    echo "  curl -sSL https://raw.githubusercontent.com/sjlleo/nexttrace/main/install.sh | sudo bash"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Download binary
DOWNLOAD_URL="https://github.com/vinsec/nexttrace_exporter/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}"

echo "Downloading ${BINARY_NAME}..."
if command -v curl &> /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "/tmp/${BINARY_NAME}"
elif command -v wget &> /dev/null; then
    wget "$DOWNLOAD_URL" -O "/tmp/${BINARY_NAME}"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

# Install binary
echo "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
chmod +x "/tmp/${BINARY_NAME}"
mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

# Set capabilities (alternative to running as root)
echo "Setting capabilities..."
setcap cap_net_raw+ep "${INSTALL_DIR}/${BINARY_NAME}" || echo -e "${YELLOW}Warning: Failed to set capabilities${NC}"

# Create config directory
echo "Creating configuration directory..."
mkdir -p "$CONFIG_DIR"

# Create default config if it doesn't exist
if [ ! -f "$CONFIG_DIR/config.yml" ]; then
    echo "Creating default configuration..."
    cat > "$CONFIG_DIR/config.yml" <<EOF
targets:
  - host: 8.8.8.8
    name: google_dns
    interval: 5m
    max_hops: 30
    nexttrace_args: []
EOF
    echo -e "${GREEN}Default configuration created at $CONFIG_DIR/config.yml${NC}"
    echo -e "${YELLOW}Please edit the configuration file to add your targets${NC}"
fi

# Install systemd service (if systemd is available)
if command -v systemctl &> /dev/null; then
    echo "Installing systemd service..."
    cat > /etc/systemd/system/nexttrace_exporter.service <<EOF
[Unit]
Description=NextTrace Exporter
Documentation=https://github.com/vinsec/nexttrace_exporter
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${INSTALL_DIR}/${BINARY_NAME} --config.file=${CONFIG_DIR}/config.yml --web.listen-address=:9101
ExecReload=/bin/kill -HUP \$MAINPID
Restart=on-failure
RestartSec=5s
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    echo -e "${GREEN}Systemd service installed${NC}"
    echo "To enable and start the service:"
    echo "  sudo systemctl enable nexttrace_exporter"
    echo "  sudo systemctl start nexttrace_exporter"
fi

echo ""
echo -e "${GREEN}Installation completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "  1. Edit configuration: $CONFIG_DIR/config.yml"
echo "  2. Start the service:"
echo "     sudo systemctl start nexttrace_exporter"
echo "  3. Check status:"
echo "     sudo systemctl status nexttrace_exporter"
echo "  4. View metrics:"
echo "     curl http://localhost:9101/metrics"
echo ""
echo "For more information, visit:"
echo "  https://github.com/vinsec/nexttrace_exporter"
