#!/bin/bash

set -euo pipefail

echo "Detecting OS and architecture..."

OS=$(uname -s)
ARCH=$(uname -m)

BINARY_NAME="slothctl"
DOWNLOAD_URL=""
INSTALL_DIR="/usr/local/bin"

# Get the latest tag from GitHub (or use a specific version)
# For this example, we'll hardcode v1.0.0 as requested.
VERSION="v1.0.0"

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}-linux-amd64"
                ;;
            aarch64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}-linux-arm64"
                ;;
            *)
                echo "Unsupported architecture: $ARCH on Linux"
                exit 1
                ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}-darwin-amd64"
                ;;
            arm64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}-darwin-arm64"
                ;;
            *)
                echo "Unsupported architecture: $ARCH on macOS"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

echo "Downloading $BINARY_NAME from $DOWNLOAD_URL..."
curl -L "$DOWNLOAD_URL" -o "/tmp/${BINARY_NAME}"

echo "Making $BINARY_NAME executable..."
chmod +x "/tmp/${BINARY_NAME}"

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "/tmp/${BINARY_NAME}" "$INSTALL_DIR/${BINARY_NAME}"

echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."
