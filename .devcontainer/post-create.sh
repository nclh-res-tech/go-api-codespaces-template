#!/bin/bash

echo "🚀 Setting up Go API development environment..."

# Check if setup has been run
if grep -q "{{MODULE_PATH}}" go.mod 2>/dev/null; then
    echo ""
    echo "⚠️  Template not yet configured!"
    echo ""
    echo "The GitHub Action may still be running to configure your project."
    echo "Wait a moment and run: git pull"
    echo ""
    echo "Or run the setup script manually:"
    echo ""
    echo "    ./setup.sh"
    echo ""
    echo "Skipping Go module download until setup is complete."
    echo ""
else
    # Download Go dependencies
    echo "📦 Downloading Go modules..."
    export GOTOOLCHAIN=local
    go mod tidy -compat=1.22 || echo "Warning: go mod tidy had issues, continuing..."
    go mod download || echo "Warning: go mod download had issues, continuing..."

    # Install development tools
    echo "🔧 Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || true
    go install github.com/go-delve/delve/cmd/dlv@latest || true

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

# Always exit successfully - don't fail container creation
exit 0