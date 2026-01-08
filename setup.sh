#!/bin/bash
# Setup script for the Go API template
# Run this after creating a new repository from this template

set -e

echo "🚀 Go API Template Setup"
echo "========================"
echo ""

# Get service name
read -p "Enter your service name (e.g., my-api-service): " SERVICE_NAME
if [ -z "$SERVICE_NAME" ]; then
    echo "❌ Service name is required"
    exit 1
fi

# Get module path
read -p "Enter your Go module path (e.g., github.com/myorg/my-api): " MODULE_PATH
if [ -z "$MODULE_PATH" ]; then
    echo "❌ Module path is required"
    exit 1
fi

echo ""
echo "📝 Configuration:"
echo "  Service Name: $SERVICE_NAME"
echo "  Module Path: $MODULE_PATH"
echo ""

read -p "Proceed with setup? (y/n): " CONFIRM
if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "Setup cancelled"
    exit 0
fi

echo ""
echo "🔧 Replacing placeholders..."

# Function to replace in files (cross-platform)
replace_in_files() {
    local search="$1"
    local replace="$2"
    
    # Find all relevant files and replace
    find . -type f \( -name "*.go" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.md" -o -name "Makefile" -o -name "Dockerfile" \) \
        -not -path "./.git/*" \
        -not -path "./vendor/*" \
        -exec sed -i.bak "s|$search|$replace|g" {} \;
    
    # Clean up backup files
    find . -name "*.bak" -delete 2>/dev/null || true
}

# Replace placeholders
replace_in_files "{{MODULE_PATH}}" "$MODULE_PATH"
replace_in_files "{{SERVICE_NAME}}" "$SERVICE_NAME"

echo "✅ Placeholders replaced"

# Update go.mod
echo ""
echo "📦 Updating go.mod..."
cat > go.mod << EOF
module $MODULE_PATH

go 1.22

require (
	github.com/gin-contrib/zap v0.2.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.19.0
	github.com/google/uuid v1.6.0
	github.com/knadh/koanf/parsers/yaml v0.1.0
	github.com/knadh/koanf/providers/env v0.1.0
	github.com/knadh/koanf/providers/file v0.1.0
	github.com/knadh/koanf/v2 v2.1.0
	github.com/prometheus/client_golang v1.19.0
	github.com/swaggest/openapi-go v0.2.50
	github.com/swaggest/swgui v1.8.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.49.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0
	go.opentelemetry.io/otel/exporters/prometheus v0.46.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
	go.opentelemetry.io/otel/sdk/metric v1.24.0
	go.uber.org/zap v1.27.0
)
EOF

echo "✅ go.mod created"

# Download dependencies
echo ""
echo "📥 Downloading dependencies..."
go mod tidy

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Review and customize the generated code"
echo "  2. Add your own models in internal/models/"
echo "  3. Implement your services in internal/services/"
echo "  4. Add routes in internal/routes/"
echo "  5. Run 'make run' to start the server"
echo ""
echo "Useful commands:"
echo "  make run          - Run the server locally"
echo "  make test         - Run tests"
echo "  make lint         - Run linter"
echo "  make docker-build - Build Docker image"
echo "  make help         - Show all available commands"