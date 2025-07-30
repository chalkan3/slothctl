#!/bin/bash

set -euo pipefail

# --- Configuration ---
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
# Fixed URL as requested by the user
DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/v1.0.2/slothctl_1.0.2_linux_amd64.tar.gz"

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Main Logic ---
echo "Starting slothctl installation from fixed URL..."

echo "Downloading from: $DOWNLOAD_URL"

# 1. Download and extract the binary
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT # Cleanup on exit

DOWNLOAD_CMD=""
if command_exists curl; then
    DOWNLOAD_CMD="curl -sL"
elif command_exists wget; then
    DOWNLOAD_CMD="wget -qO-"
else
    echo "Error: Neither curl nor wget are installed."
    exit 1
fi

if ! $DOWNLOAD_CMD "$DOWNLOAD_URL" | tar -xz -C "$TEMP_DIR"; then
    echo "Error: Failed to download or extract the binary. Please check the URL and your network connection."
    echo "URL: $DOWNLOAD_URL"
    exit 1
fi

# 2. Install the binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
# Find the binary in the temp directory. It might be in the root or a subdirectory.
FOUND_BINARY=$(find "$TEMP_DIR" -type f -name "$BINARY_NAME" | head -n 1)

if [ -n "$FOUND_BINARY" ]; then
    echo "Binary found at: $FOUND_BINARY"
    sudo mv "$FOUND_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Error: Binary '$BINARY_NAME' not found in the downloaded archive."
    echo "Contents of the temporary directory:"
    ls -lR "$TEMP_DIR"
    exit 1
fi

echo "slothctl installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
echo "You can now run 'slothctl' from your terminal."
