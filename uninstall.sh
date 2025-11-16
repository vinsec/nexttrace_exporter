#!/bin/bash

# Uninstallation script for nexttrace_exporter

set -e

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
CONFIG_DIR="${CONFIG_DIR:-/etc/nexttrace_exporter}"
BINARY_NAME="nexttrace_exporter"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}NextTrace Exporter Uninstallation Script${NC}"
echo "=========================================="

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Error: This script must be run as root${NC}" 
   exit 1
fi

# Stop and disable systemd service
if command -v systemctl &> /dev/null; then
    if systemctl is-active --quiet nexttrace_exporter; then
        echo "Stopping nexttrace_exporter service..."
        systemctl stop nexttrace_exporter
    fi
    
    if systemctl is-enabled --quiet nexttrace_exporter 2>/dev/null; then
        echo "Disabling nexttrace_exporter service..."
        systemctl disable nexttrace_exporter
    fi
    
    if [ -f /etc/systemd/system/nexttrace_exporter.service ]; then
        echo "Removing systemd service file..."
        rm /etc/systemd/system/nexttrace_exporter.service
        systemctl daemon-reload
    fi
fi

# Remove binary
if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    echo "Removing binary..."
    rm "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${YELLOW}Binary not found at ${INSTALL_DIR}/${BINARY_NAME}${NC}"
fi

# Ask about config files
if [ -d "$CONFIG_DIR" ]; then
    echo ""
    read -p "Remove configuration directory $CONFIG_DIR? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo "Configuration directory removed"
    else
        echo "Configuration directory preserved"
    fi
fi

echo ""
echo -e "${GREEN}Uninstallation completed!${NC}"
