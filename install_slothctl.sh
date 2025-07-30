#!/bin/bash

set -euo pipefail

# --- Configuration ---
BINARY_NAME="slothctl"
INSTALL_DIR="/usr/local/bin"
VERSION="v1.0.1" # The version to install

# --- Temporary Files ---
TEMP_ARCHIVE="/tmp/${BINARY_NAME}_archive.tar.gz"
EXTRACT_DIR="/tmp/${BINARY_NAME}_extracted"
BUILD_LOG="/tmp/${BINARY_NAME}_build.log"

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

install_dependencies() {
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
}

# --- Main Installation Logic ---

echo "Starting $BINARY_NAME installation..."

# Detect OS and architecture
echo "Detecting OS and architecture..."
OS=$(uname -s)
ARCH=$(uname -m)

# Determine the archive name based on Goreleaser's name_template
# This now correctly removes the 'v' from the version for the archive name
ARCHIVE_NAME="${BINARY_NAME}_${VERSION#v}_${OS,,}_${ARCH}"

# Construct download URL
DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${ARCHIVE_NAME}.tar.gz"

# --- Download and Verify ---
echo "Downloading $BINARY_NAME from $DOWNLOAD_URL..."

set -x # Enable command tracing for download
curl -L "$DOWNLOAD_URL" -o "$TEMP_ARCHIVE" || { echo "Error: curl download failed."; exit 1; }
set +x # Disable command tracing

if [ ! -f "$TEMP_ARCHIVE" ]; then
    echo "Error: Downloaded file not found at $TEMP_ARCHIVE."
    exit 1
fi

FILE_SIZE=$(stat -c%s "$TEMP_ARCHIVE")
if [ "$FILE_SIZE" -lt 1000 ]; then # Assuming a real tar.gz is > 1KB
    echo "Error: Downloaded file is too small ($FILE_SIZE bytes). Likely an error page or incomplete download."
    echo "Content of downloaded file:"
    cat "$TEMP_ARCHIVE"
    exit 1
fi

FILE_TYPE=$(file -b "$TEMP_ARCHIVE")
if [[ ! "$FILE_TYPE" =~ "gzip compressed data" && ! "$FILE_TYPE" =~ "Zip archive" ]]; then
    echo "Error: Downloaded file is not a valid archive. Detected type: $FILE_TYPE."
    echo "Content of downloaded file:"
    cat "$TEMP_ARCHIVE"
    exit 1
fi

# --- Extraction ---
echo "Extracting $BINARY_NAME..."
mkdir -p "$EXTRACT_DIR"
tar -xzf "$TEMP_ARCHIVE" -C "$EXTRACT_DIR" || { echo "Error: Failed to extract archive."; exit 1; }

# The binary is usually directly inside the extracted directory
# We'll assume it's directly in the extracted directory for now.
EXTRACTED_BINARY_PATH="${EXTRACT_DIR}/${BINARY_NAME}"

if [ ! -f "$EXTRACTED_BINARY_PATH" ]; then
    echo "Error: Extracted binary not found at $EXTRACTED_BINARY_PATH."
    echo "Contents of extraction directory ($EXTRACT_DIR):"
    ls -la "$EXTRACT_DIR"
    exit 1
fi

# --- Installation ---
echo "Making $BINARY_NAME executable..."
chmod +x "$EXTRACTED_BINARY_PATH" || { echo "Error: Failed to make binary executable."; exit 1; }

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$EXTRACTED_BINARY_PATH" "$INSTALL_DIR/${BINARY_NAME}" || { echo "Error: Failed to move binary to install directory."; exit 1; }

# --- Cleanup ---
echo "Cleaning up temporary files..."
rm "$TEMP_ARCHIVE" || { echo "Warning: Failed to remove temporary archive."; }
rm -rf "$EXTRACT_DIR" || { echo "Warning: Failed to remove extraction directory."; }

# --- Success ---
echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."
