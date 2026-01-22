#!/usr/bin/env bash
# Dedicated cleanup script for validatord project
# This script cleans all build artifacts, test files, and temporary files

set -euo pipefail

# Get the repository root directory
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🧹 Cleaning validatord project..."
echo "Repository root: $REPO_ROOT"

cd "$REPO_ROOT"

# Clean build artifacts
echo "Cleaning build artifacts..."
rm -f validatord
rm -f coverage.out

# Clean test artifacts
echo "Cleaning test artifacts..."
find . -type f -name "*.test" -delete
find . -type f -name "*.out" -delete

# Clean temporary files
echo "Cleaning temporary files..."
find . -type f -name "*.tmp" -delete
find . -type f -name "*.temp" -delete
find . -type f -name "*.log" -delete

# Clean profiling files
echo "Cleaning profiling files..."
find . -type f -name "*.prof" -delete
find . -type f -name "*.pprof" -delete

# Clean script-specific outputs
echo "Cleaning script outputs..."
rm -rf /tmp/fluffy-check

echo "✅ Cleanup complete!"
