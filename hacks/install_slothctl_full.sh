#!/bin/bash

set -euo pipefail

echo "Detecting OS and architecture..."

OS=$(uname -s)
ARCH=$(uname -m)

BINARY_NAME="slothctl"
DOWNLOAD_URL=""
INSTALL_DIR="/usr/local/bin"
TEMP_ARCHIVE="/tmp/${BINARY_NAME}_archive.tar.gz"
EXTRACT_DIR="/tmp/${BINARY_NAME}_extracted"

# Get the latest tag from GitHub (or use a specific version)
# For this example, we'll hardcode v1.0.0 as requested.
VERSION="v1.0.0"

# Determine the archive name based on Goreleaser's name_template
ARCHIVE_NAME="${BINARY_NAME}_${VERSION#v}_${OS,,}_${ARCH}"

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${ARCHIVE_NAME}.tar.gz"
                ;;
            aarch64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${ARCHIVE_NAME}.tar.gz"
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
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${ARCHIVE_NAME}.tar.gz"
                ;;
            arm64)
                DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/${VERSION}/${ARCHIVE_NAME}.tar.gz"
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

# --- Start Debugging Block ---
set -x # Enable command tracing
curl -L "$DOWNLOAD_URL" -o "$TEMP_ARCHIVE"
set +x # Disable command tracing

# Check if download was successful and file is valid
if [ ! -f "$TEMP_ARCHIVE" ]; then
    echo "Error: Downloaded file not found at $TEMP_ARCHIVE."
    exit 1
fi

FILE_SIZE=$(stat -c%s "$TEMP_ARCHIVE")
if [ "$FILE_SIZE" -lt 1000 ]; then # Assuming a real tar.gz is > 1KB
    echo "Error: Downloaded file is too small ($FILE_SIZE bytes). Likely an error page."
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
# --- End Debugging Block ---

echo "Extracting $BINARY_NAME..."
mkdir -p "$EXTRACT_DIR"
tar -xzf "$TEMP_ARCHIVE" -C "$EXTRACT_DIR"

# The binary is usually directly inside the extracted directory
# or sometimes in a subdirectory named after the archive.
# We'll assume it's directly in the extracted directory for now.
# If not, a find command might be needed.
EXTRACTED_BINARY_PATH="${EXTRACT_DIR}/${BINARY_NAME}"

echo "Making $BINARY_NAME executable..."
chmod +x "$EXTRACTED_BINARY_PATH"

echo "Moving $BINARY_NAME to $INSTALL_DIR (requires sudo)..."
sudo mv "$EXTRACTED_BINARY_PATH" "$INSTALL_DIR/${BINARY_NAME}"

echo "Cleaning up temporary files..."
rm "$TEMP_ARCHIVE"
rm -rf "$EXTRACT_DIR"

echo "$BINARY_NAME installed successfully to $INSTALL_DIR"
echo "You can now run '$BINARY_NAME' from your terminal."