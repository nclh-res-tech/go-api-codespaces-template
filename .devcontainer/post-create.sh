#!/bin/bash
set -e

echo "🚀 Setting up Go API development environment..."

# Download Go dependencies
echo "📦 Downloading Go modules..."
go mod download

# Install development tools
echo "🔧 Installing development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/go-delve/delve/cmd/dlv@latest

# Check if setup has been run
if grep -q "{{MODULE_PATH}}" go.mod 2>/dev/null; then
    echo ""
    echo "⚠️  Template not yet configured!"
    echo ""
    echo "Please run the setup script to configure your service:"
    echo ""
    echo "    ./setup.sh"
    echo ""
else
    echo ""
    echo "✅ Development environment ready!"
    echo ""
    echo "Quick commands:"
    echo "  make run     - Start the API server"
    echo "  make test    - Run tests"
    echo "  make build   - Build the binary"
    echo "  make lint    - Run linter"
    echo ""
fi