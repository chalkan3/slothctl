#!/bin/bash

set -euo pipefail

# --- Configuration ---
REPO="chalkan3/slothctl"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Main Logic ---
echo "Starting slothctl installation..."

# 1. Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected OS: $OS, Architecture: $ARCH"

# 2. Get the latest release version from GitHub API
LATEST_RELEASE_URL="https://api.github.com/repos/${REPO}/releases/latest"
echo "Fetching latest release from: $LATEST_RELEASE_URL"

# Use curl or wget to fetch release data
if command_exists curl; then
    LATEST_VERSION=$(curl -s "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -E 's/.*"tag_name": "(v[0-9]+\.[0-9]+\.[0-9]+)".*/\1/')
elif command_exists wget; then
    LATEST_VERSION=$(wget -qO- "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -E 's/.*"tag_name": "(v[0-9]+\.[0-9]+\.[0-9]+)".*/\1/')
else
    echo "Error: Neither curl nor wget are installed. Please install one and try again."
    exit 1
fi

if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not fetch the latest release version."
    exit 1
fi

echo "Latest version is: $LATEST_VERSION"

# 3. Construct the download URL
# slothctl_0.1.0_linux_amd64.tar.gz
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}_${LATEST_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# 4. Download and extract the binary
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT # Cleanup on exit

if command_exists curl; then
    curl -sL "$DOWNLOAD_URL" | tar -xz -C "$TEMP_DIR"
elif command_exists wget; then
    wget -qO- "$DOWNLOAD_URL" | tar -xz -C "$TEMP_DIR"
else
    echo "Error: Neither curl nor wget are installed."
    exit 1
fi

# 5. Install the binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -f "${TEMP_DIR}/${BINARY_NAME}" ]; then
    sudo mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Error: Binary not found in the downloaded archive."
    exit 1
fi

echo "slothctl installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
echo "You can now run 'slothctl' from your terminal."