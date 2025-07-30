#!/bin/bash

set -euo pipefail

echo "Detecting OS and architecture..."

OS=$(uname -s)
ARCH=$(uname -m)

BINARY_NAME="slothctl"
DOWNLOAD_URL=""
INSTALL_DIR="/usr/local/bin"
TEMP_FILE="/tmp/${BINARY_NAME}.tar.gz"
EXTRACTED_BINARY="/tmp/${BINARY_NAME}"

# Get the latest tag from GitHub (or use a specific version)
# For this example, we'll hardcode v1.0.0 as requested.
VERSION="v1.0.0"

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}_${VERSION}_linux_amd64.tar.gz"
                ;;
            aarch64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}_${VERSION}_linux_arm64.tar.gz"
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
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}_${VERSION}_darwin_amd64.tar.gz"
                ;;
            arm64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${BINARY_NAME}_${VERSION}_darwin_arm64.tar.gz"
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
curl -L "$DOWNLOAD_URL" -o "$TEMP_FILE"

echo "Extracting $BINARY_NAME..."
tar -xzf "$TEMP_FILE" -C "/tmp/"

echo "Making $BINARY_NAME executable..."
chmod +x "$EXTRACTED_BINARY"

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$EXTRACTED_BINARY" "$INSTALL_DIR/${BINARY_NAME}"

echo "Cleaning up temporary files..."
rm "$TEMP_FILE"

echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."