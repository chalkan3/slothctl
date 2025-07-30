#!/bin/bash

set -euo pipefail

REPO_URL="https://github.com/chalkan3/slothctl.git"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
TEMP_DIR="/tmp/slothctl_install"
BUILD_LOG="/tmp/slothctl_build.log"

echo "Starting slothctl installation from source..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Update system and install basic tools
echo "Updating system and installing necessary tools (file, go, git)..."
sudo pacman -Syu --noconfirm

# Install 'file' command if not present (used for binary verification)
if ! command_exists file; then
    echo "'file' command not found. Installing..."
    sudo pacman -Sy --noconfirm file
fi

# Install Go
if ! command_exists go; then
    echo "Go not found. Installing Go..."
    sudo pacman -Sy --noconfirm go
else
    echo "Go is already installed."
fi

# Install Git
if ! command_exists git; then
    echo "Git not found. Installing Git..."
    sudo pacman -Sy --noconfirm git
else
    echo "Git is already installed."
}

echo "Cloning repository $REPO_URL..."
mkdir -p "$TEMP_DIR"
git clone "$REPO_URL" "$TEMP_DIR"

echo "Building $BINARY_NAME from source... (logging to $BUILD_LOG)"
cd "$TEMP_DIR"

# Execute go build and capture its output and exit code
if ! go build -o "$BINARY_NAME" ./cmd/slothctl > "$BUILD_LOG" 2>&1; then
    echo "Error: go build failed! Check $BUILD_LOG for details."
    cat "$BUILD_LOG"
    exit 1
fi

echo "Build completed. Verifying generated binary..."

# Check if the generated file is actually a binary
if ! file "$BINARY_NAME" | grep -q "executable"; then
    echo "Error: Generated file is not an executable binary! Content of /tmp/${BINARY_NAME} (first 20 lines):"
    head -n 20 "$BINARY_NAME"
    echo "Full build log in $BUILD_LOG."
    exit 1
fi

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."