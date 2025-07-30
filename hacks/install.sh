#!/bin/bash

set -euo pipefail

# --- Configuration ---
REPO="chalkan3/slothctl"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
DEFAULT_VERSION="v1.0.2" # Fallback version as requested

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
    # A more robust way to parse JSON with grep/sed
    LATEST_VERSION=$(curl -s "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -e 's/.*"tag_name": *"\([^"]*\)".*/\1/')
elif command_exists wget; then
    LATEST_VERSION=$(wget -qO- "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -e 's/.*"tag_name": *"\([^"]*\)".*/\1/')
else
    echo "Error: Neither curl nor wget are installed. Please install one and try again."
    exit 1
fi

# Fallback to default version if API call fails or returns no version
if [ -z "$LATEST_VERSION" ]; then
    echo "Warning: Could not fetch the latest release version. Falling back to default version: $DEFAULT_VERSION"
    LATEST_VERSION="$DEFAULT_VERSION"
fi

echo "Using version: $LATEST_VERSION"

# 3. Construct the download URL
# The asset name format is slothctl_1.0.2_linux_amd64.tar.gz (version without 'v')
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
    # Handle cases where the binary is inside a subdirectory in the tarball
    FOUND_BINARY=$(find "$TEMP_DIR" -type f -name "$BINARY_NAME" | head -n 1)
    if [ -n "$FOUND_BINARY" ]; then
        sudo mv "$FOUND_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "Error: Binary not found in the downloaded archive."
        echo "Contents of the temporary directory:"
        ls -lR "$TEMP_DIR"
        exit 1
    fi
fi

echo "slothctl installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
echo "You can now run 'slothctl' from your terminal."
