#!/bin/bash

set -euo pipefail

# --- Configuration ---
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="slothctl"
DOWNLOAD_URL="https://github.com/chalkan3/slothctl/releases/download/v1.0.2/slothctl_1.0.2_linux_amd64.tar.gz"

# --- Emojis and Colors ---
DOWNLOAD_ICON="⬇️"
INSTALL_ICON="⚙️"
SUCCESS_ICON="✅"

C_RESET='\033[0m'
C_GREEN='\033[0;32m'
C_YELLOW='\033[1;33m'
C_CYAN='\033[0;36m'
C_BOLD='\033[1m'

# --- Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

print_banner() {
    echo -e "${C_GREEN}"
    echo "##################################################################################"
    echo "#                                                                                #"
    echo "#      _|_|_|   _|                _|       _|                      _|       _|   #"
    echo "#    _|         _|     _|_|     _|_|_|_|   _|_|_|       _|_|_|   _|_|_|_|   _|   #"
    echo "#      _|_|     _|   _|    _|     _|       _|    _|   _|           _|       _|   #"
    echo "#          _|   _|   _|    _|     _|       _|    _|   _|           _|       _|   #"
    echo "#    _|_|_|     _|     _|_|         _|_|   _|    _|     _|_|_|       _|_|   _|   #"
    echo "#                                                                                #"
    echo "##################################################################################"
    echo -e "${C_RESET}"
}

# --- Main Logic ---
print_banner

echo -e "${C_BOLD}--- Sloth Control Installer ---${C_RESET}"

# --- Download Step ---
echo
echo -e "${C_CYAN}${DOWNLOAD_ICON} Downloading slothctl...${C_RESET}"
echo -e "   ${C_BOLD}From:${C_RESET} $DOWNLOAD_URL"

TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT # Cleanup on exit

DOWNLOAD_CMD=""
if command_exists curl; then
    DOWNLOAD_CMD="curl -sL --progress-bar"
elif command_exists wget; then
    DOWNLOAD_CMD="wget -qO- --show-progress"
else
    echo "Error: Neither curl nor wget are installed."
    exit 1
fi

if ! $DOWNLOAD_CMD "$DOWNLOAD_URL" | tar -xz -C "$TEMP_DIR"; then
    echo "Error: Failed to download or extract the binary. Please check the URL and your network connection."
    exit 1
fi
echo -e "${C_GREEN}   Download complete.${C_RESET}"

# --- Installation Step ---
echo
echo -e "${C_CYAN}${INSTALL_ICON} Installing slothctl...${C_RESET}"

FOUND_BINARY=$(find "$TEMP_DIR" -type f -name "$BINARY_NAME" | head -n 1)

if [ -n "$FOUND_BINARY" ]; then
    echo "   Binary found, moving to $INSTALL_DIR"
    sudo mv "$FOUND_BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Error: Binary '$BINARY_NAME' not found in the downloaded archive."
    exit 1
fi
echo -e "${C_GREEN}   Installation complete.${C_RESET}"

# --- Success Message ---
echo
echo -e "${C_GREEN}--------------------------------------------------${C_RESET}"
echo -e "${C_GREEN}${SUCCESS_ICON} ${C_BOLD}slothctl was installed successfully!${C_RESET}"
echo -e "${C_GREEN}--------------------------------------------------${C_RESET}"
echo
echo -e "${C_BOLD}Next step is to initialize the configuration:${C_RESET}"
echo
echo -e "  ${C_YELLOW}slothctl configure init${C_RESET}"
echo