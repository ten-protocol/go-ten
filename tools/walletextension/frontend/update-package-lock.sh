#!/bin/bash

# Script to regenerate package-lock.json for Docker compatibility
# This is needed because we use pnpm workspace but Docker needs npm package-lock.json

set -e  # Exit on any error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$SCRIPT_DIR"
TEMP_DIR="/tmp/npm-lock-update-$$"

echo "🔄 Regenerating package-lock.json for Docker compatibility..."

# Create temporary directory
mkdir -p "$TEMP_DIR"

# Copy package.json to temp directory
cp "$FRONTEND_DIR/package.json" "$TEMP_DIR/"

# Change to temp directory and generate clean package-lock.json
cd "$TEMP_DIR"
echo "📦 Generating clean npm package-lock.json..."
npm install --package-lock-only

# Copy the generated package-lock.json back to project
cp package-lock.json "$FRONTEND_DIR/"

# Clean up temp directory
cd /
rm -rf "$TEMP_DIR"

echo "✅ package-lock.json updated successfully!"
echo "📋 The following dependencies are now included:"
grep '".*":' "$FRONTEND_DIR/package.json" | head -10
echo "   ... and more"
echo ""
echo "🐳 You can now run Docker build successfully."
