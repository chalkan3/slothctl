#!/bin/bash

set -euo pipefail

# --- Configuration ---
# Set this to 'release' for a full release, or 'snapshot' for a local build.
# Use 'release' only when you have a Git tag pushed and GITHUB_TOKEN set.
BUILD_MODE="release"

# --- Functions ---
function usage() {
    echo "Usage: $0 [release|snapshot]"
    echo "  release: Performs a full goreleaser release (requires GITHUB_TOKEN and Git tag)."
    echo "  snapshot: Performs a local goreleaser build (snapshot) for testing."
    exit 1
}

# --- Main Logic ---
if [ "$#" -gt 0 ]; then
    BUILD_MODE="$1"
fi

if [ "$BUILD_MODE" == "release" ]; then
    echo "--- Performing full Goreleaser Release ---"
    if [ -z "${GITHUB_TOKEN}" ]; then
        echo "Error: GITHUB_TOKEN environment variable is not set."
        echo "Please set it before running a full release (e.g., export GITHUB_TOKEN=\"your_token\")."
        exit 1
    fi
    echo "Running: goreleaser release --clean"
    goreleaser release --clean
    echo "Goreleaser release completed. Check your GitHub repository releases page."
elif [ "$BUILD_MODE" == "snapshot" ]; then
    echo "--- Performing Goreleaser Snapshot Build (local only) ---"
    echo "Running: goreleaser build --clean --snapshot"
    goreleaser build --clean --snapshot
    echo "Goreleaser snapshot build completed. Binaries are in the 'dist/' directory."
else
    usage
fi
