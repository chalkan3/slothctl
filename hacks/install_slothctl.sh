#!/bin/bash

set -euo pipefail

FULL_INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/chalkan3/slothctl/master/hacks/install_slothctl_full.sh"

echo "Downloading and executing the full slothctl installer..."

# Use a cache-busting parameter to ensure the latest version is fetched
curl -fsSL "${FULL_INSTALL_SCRIPT_URL}?$(date +%s)" | bash
