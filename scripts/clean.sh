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
