#!/bin/bash

set -euo pipefail

# URL to the full installation script, using a specific commit hash for immutability
FULL_INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/chalkan3/slothctl/a8666bfba274ac55d6047683d211862348c58c11/hacks/install_slothctl_full.sh"

echo "Downloading and executing the slothctl installer..."

# Execute the full installer directly
curl -fsSL "${FULL_INSTALL_SCRIPT_URL}" | bash
