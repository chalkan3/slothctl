#!/bin/bash

set -euo pipefail

REPO_URL="https://github.com/chalkan3/slothctl.git"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
TEMP_DIR="/tmp/slothctl_install"

echo "Starting slothctl installation from source..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

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

echo "Building $BINARY_NAME from source..."
cd "$TEMP_DIR"
go build -o "$BINARY_NAME" ./cmd/slothctl

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."