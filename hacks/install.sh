#!/bin/bash

set -euo pipefail

# --- Configuration ---
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/v1.0.2/slothctl_1.0.2_linux_amd64.tar.gz"

# --- Colors ---
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RESET='\033[0m'

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

print_banner() {
    echo -e "${COLOR_GREEN}"
    echo ' SSS  L      OOO  TTTTT H   H   CCC  TTTTT L     '
    echo 'S     L     O   O   T   H   H  C       T   L     '
    echo ' SSS  L     O   O   T   HHHHH  C       T   L     '
    echo '    S L     O   O   T   H   H  C       T   L     '
    echo ' SSS  LLLLL  OOO    T   H   H   CCC    T   LLLLL '
    echo -e "${COLOR_RESET}"
    echo "--- Sloth Control Installer ---"
    echo
}

# --- Main Logic ---
print_banner

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

echo
echo -e "${COLOR_GREEN}slothctl installed successfully! ${COLOR_RESET}"
echo
echo "Next step: Initialize the configuration by running:"
_YELLOW}"
echo -e "  ${COLOR_YELLOW}slothctl configure init${COLOR_RESET}"
