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
#!/bin/bash
# Cleanup script for the scripts directory
# Removes temporary files, logs, and build artifacts

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🧹 Cleaning scripts directory: $SCRIPT_DIR"

# Remove log files
if find "$SCRIPT_DIR" -maxdepth 1 -name "*.log" -type f | grep -q .; then
    echo "  Removing log files..."
    find "$SCRIPT_DIR" -maxdepth 1 -name "*.log" -type f -delete
else
    echo "  No log files found"
fi

# Remove temporary files
if find "$SCRIPT_DIR" -maxdepth 1 \( -name "*.tmp" -o -name "*.temp" \) -type f | grep -q .; then
    echo "  Removing temporary files..."
    find "$SCRIPT_DIR" -maxdepth 1 \( -name "*.tmp" -o -name "*.temp" \) -type f -delete
else
    echo "  No temporary files found"
fi

# Remove editor backup files
if find "$SCRIPT_DIR" -maxdepth 1 -name "*~" -type f | grep -q .; then
    echo "  Removing editor backup files..."
    find "$SCRIPT_DIR" -maxdepth 1 -name "*~" -type f -delete
else
    echo "  No backup files found"
fi

# Remove node_modules if present (JavaScript dependencies)
if [ -d "$SCRIPT_DIR/node_modules" ]; then
    echo "  Removing node_modules..."
    rm -rf "$SCRIPT_DIR/node_modules"
else
    echo "  No node_modules directory found"
fi

# Remove package-lock.json if present
if [ -f "$SCRIPT_DIR/package-lock.json" ]; then
    echo "  Removing package-lock.json..."
    rm -f "$SCRIPT_DIR/package-lock.json"
else
    echo "  No package-lock.json found"
fi

echo "✅ Scripts directory cleaned successfully!"
