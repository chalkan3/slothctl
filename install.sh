#!/bin/bash

set -euo pipefail

# URL to the full installation script, using a specific commit hash for immutability
FULL_INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/chalkan3/slothctl/50da7ae4250bce55ebbfd7ce8a81c73abe620a26/hacks/install_slothctl_full.sh"

echo "Downloading and executing the slothctl installer..."

# Execute the full installer directly
curl -fsSL "${FULL_INSTALL_SCRIPT_URL}" | bash
