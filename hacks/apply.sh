#!/bin/bash

set -euo pipefail

echo "Building slothctl binary..."
go build -o slothctl ./cmd/slothctl

echo "Initializing Incus (non-interactively) with preseed file..."
cat /home/igor/.chalkan3/slothctl/hacks/incus-preseed.yaml | sudo incus admin init --preseed || true # || true to ignore if already initialized

echo "Creating Arch Linux container 'slothctl-arch-test'..."
sudo incus launch images:archlinux/current slothctl-arch-test

echo "Waiting for container to come online..."
sleep 10 # Give it some time to boot up

echo "Pushing slothctl binary to container..."
sudo incus file push slothctl slothctl-arch-test/usr/local/bin/slothctl --mode=0755

echo "Running slothctl configure init --mode apply inside the container..."
sudo incus exec slothctl-arch-test -- /usr/local/bin/slothctl configure init --mode apply

echo "Apply script finished."
