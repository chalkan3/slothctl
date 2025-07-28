#!/bin/bash

set -euo pipefail

echo "Stopping and deleting container 'slothctl-arch-test'..."
sudo incus stop slothctl-arch-test || true
sudo incus delete slothctl-arch-test || true

echo "Destroy script finished."
