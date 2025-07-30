#!/bin/bash

set -euo pipefail

# --- Configuration ---
BINARY_NAME="slothctl"
INSTALL_DIR="/usr/local/bin"
REPO_URL="https://github.com/chalkan3/slothctl.git"
TEMP_DIR="/tmp/${BINARY_NAME}_install"
BUILD_LOG="/tmp/${BINARY_NAME}_build.log"

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Main Installation Logic ---

echo "Starting $BINARY_NAME installation from source..."

# Detect OS and architecture
echo "Detecting OS and architecture..."
OS=$(uname -s)
ARCH=$(uname -m)

# Install dependencies
echo "Updating system and installing necessary tools (file, go, git)..."
sudo pacman -Syu --noconfirm || { echo "Error: Failed to update system or install pacman dependencies."; exit 1; }

if ! command_exists file; then
    echo "'file' command not found. Installing..."
    sudo pacman -Sy --noconfirm file || { echo "Error: Failed to install 'file' command."; exit 1; }
fi

if ! command_exists go; then
    echo "Go not found. Installing Go..."
    sudo pacman -Sy --noconfirm go || { echo "Error: Failed to install Go."; exit 1; }
else
    echo "Go is already installed."
fi

if ! command_exists git; then
    echo "Git not found. Installing Git..."
    sudo pacman -Sy --noconfirm git || { echo "Error: Failed to install Git."; exit 1; }
else
    echo "Git is already installed."
fi

# Clone repository
echo "Cloning repository $REPO_URL..."
mkdir -p "$TEMP_DIR"
git clone "$REPO_URL" "$TEMP_DIR" || { echo "Error: Failed to clone repository."; exit 1; }

# Build from source
echo "Building $BINARY_NAME from source... (logging to $BUILD_LOG)"
cd "$TEMP_DIR"
go build -o "$BINARY_NAME" ./cmd/slothctl > "$BUILD_LOG" 2>&1 || {
    echo "Error: go build failed! Check $BUILD_LOG for details."
    cat "$BUILD_LOG"
    exit 1
}

# Verify generated binary
echo "Build completed. Verifying generated binary..."
if [ ! -f "$BINARY_NAME" ]; then
    echo "Error: Generated binary not found at $BINARY_NAME."
    exit 1
fi

FILE_TYPE=$(file -b "$BINARY_NAME")
if [[ ! "$FILE_TYPE" =~ "executable" ]]; then
    echo "Error: Generated file is not an executable binary. Detected type: $FILE_TYPE."
    echo "Content of generated file (first 20 lines):"
    head -n 20 "$BINARY_NAME"
    echo "Full build log in $BUILD_LOG."
    exit 1
fi

# Install binary
echo "Making $BINARY_NAME executable..."
chmod +x "$BINARY_NAME" || { echo "Error: Failed to make binary executable."; exit 1; }

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/${BINARY_NAME}" || { echo "Error: Failed to move binary to install directory."; exit 1; }

# Cleanup
echo "Cleaning up temporary files..."
rm -rf "$TEMP_DIR" || { echo "Warning: Failed to remove temporary directory."; }
rm "$BUILD_LOG" || { echo "Warning: Failed to remove build log."; }

# Success
echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."