#!/bin/bash

set -uo pipefail

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

LATEST_VERSION=""
# Try to parse with python if available
if command_exists python3; then
    LATEST_VERSION=$(curl -s "$LATEST_RELEASE_URL" | python3 -c "import sys, json; print(json.load(sys.stdin).get('tag_name', ''))" 2>/dev/null || echo "")
elif command_exists python; then
    # For python 2
    LATEST_VERSION=$(curl -s "$LATEST_RELEASE_URL" | python -c "import sys, json; print json.load(sys.stdin).get('tag_name', '')" 2>/dev/null || echo "")
fi

# If python parsing failed or wasn't available, try grep/sed
if [ -z "$LATEST_VERSION" ]; then
    echo "Python not available or parsing failed. Trying grep/sed as a fallback."
    # The || true is to prevent the script from exiting if grep finds no match
    LATEST_VERSION=$(curl -s "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -e 's/.*"tag_name": *"\([^"]*\)".*/\1/' | tr -d '[:space:]' || true)
fi

# Fallback to default version if all else fails
if [ -z "$LATEST_VERSION" ]; then
    echo "Warning: Could not fetch the latest release version automatically. Falling back to default version: $DEFAULT_VERSION"
    LATEST_VERSION="$DEFAULT_VERSION"
fi

echo "Using version: $LATEST_VERSION"

# 3. Construct the download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}_${LATEST_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# 4. Download and extract the binary
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

# 5. Install the binary
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