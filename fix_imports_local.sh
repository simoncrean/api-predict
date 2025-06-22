#!/bin/bash

# Fix Go import errors for local development
# This script fixes the module name and import paths

set -e

echo "🔧 Fixing Go import paths for local development..."

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: main.go not found. Please run this script from the project root."
    exit 1
fi

# Backup existing go.mod if it exists
if [ -f "go.mod" ]; then
    cp go.mod go.mod.backup
    echo "📦 Backed up existing go.mod"
fi

# Remove old module files
rm -f go.mod go.sum

# Initialize Go module with local name
echo "📝 Initializing Go module with local name..."
go mod init depin-compatibility-api

# Install required dependencies
echo "📦 Installing dependencies..."
go get github.com/gin-gonic/gin@latest
go get golang.org/x/time@latest

# Tidy up
go mod tidy

# Test compilation
echo "🧪 Testing compilation..."
if go build -o /tmp/depin-api-test main.go; then
    rm -f /tmp/depin-api-test
    echo "✅ Import errors fixed! The application compiles successfully."
    echo ""
    echo "🎉 You can now run:"
    echo "  go run main.go"
    echo "  # or"
    echo "  make run"
else
    echo "❌ Still having compilation issues. Check the error messages above."
    
    # Restore backup if available
    if [ -f "go.mod.backup" ]; then
        echo "📦 Restoring backup..."
        mv go.mod.backup go.mod
    fi
    exit 1
fi

# Clean up backup
rm -f go.mod.backup
echo "✅ Setup complete!"
