# {{SERVICE_NAME}} API

A Go microservice built from the go-api-template.

## 🚀 Quick Start

### Using GitHub Codespaces (Recommended)

1. Click "Use this template" → "Create a new repository"
2. Open your new repo in GitHub Codespaces
3. The devcontainer will automatically set up your environment
4. Run the setup script to customize for your project:
   ```bash
   ./setup.sh
   ```
5. Start the server:
   ```bash
   make run
   ```

### Local Development

1. Clone your repository
2. Ensure you have Go 1.22+ installed
3. Run the setup script:
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```
4. Start the server:
   ```bash
   make run
   ```

## 📁 Project Structure

```
.
├── cmd/
│   └── api-server/          # Application entry point
├── common/                   # Shared utilities (don't modify these much)
│   ├── errors/              # Error types
│   ├── httpclient/          # HTTP client utilities
│   ├── httpserver/          # HTTP server, middleware, OpenAPI
│   ├── messagebroker/       # Kafka/message broker stubs
│   ├── observability/       # Telemetry (tracing, metrics)
│   └── utils/               # General utilities
├── config/                   # Configuration files
│   └── config.yaml          # Default configuration
├── internal/                 # Application-specific code
│   ├── config/              # Configuration loading
│   ├── core/                # Core application setup
│   │   ├── app/             # Application lifecycle
│   │   └── httpserver/      # HTTP server setup
│   ├── models/              # Data models (add your models here!)
│   ├── routes/              # Route definitions (add routes here!)
│   ├── services/            # Business logic (add services here!)
│   └── stores/              # Data access layer (add stores here!)
├── .devcontainer/           # GitHub Codespaces configuration
├── .github/                 # GitHub Actions, templates
├── .vscode/                 # VS Code settings
├── Dockerfile               # Container build
├── Makefile                 # Build commands
└── setup.sh                 # Project setup script
```

## 🔧 Development

### Available Commands

```bash
make run           # Run the server locally
make build         # Build the binary
make test          # Run tests
make lint          # Run golangci-lint
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make help          # Show all commands
```

### Adding New Features

#### 1. Add a Model (`internal/models/`)
```go
type MyModel struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

#### 2. Add a Store (`internal/stores/`)
```go
type MyStore struct {
    // database connection, etc.
}

func (s *MyStore) Create(ctx context.Context, m models.MyModel) error {
    // implementation
}
```

#### 3. Add a Service (`internal/services/`)
```go
type MyService struct {
    store *stores.MyStore
}

func (s *MyService) DoSomething(ctx context.Context) error {
    // business logic
}
```

#### 4. Add Routes (`internal/routes/routes.go`)
```go
{
    RouteSpec: httpserver.RouteSpec{
        Method:  http.MethodGet,
        GinPath: "/my-endpoint",
        Summary: "My endpoint description",
        Tags:    []string{"my-feature"},
    },
    Handler: func(c *gin.Context) {
        // handler implementation
    },
},
```

## 📡 API Endpoints

| Method | Path         | Description        |
|--------|--------------|-------------------|
| GET    | /health      | Health check      |
| GET    | /metrics     | Prometheus metrics|
| GET    | /docs        | OpenAPI docs (Swagger UI) |
| GET    | /openapi.json| OpenAPI spec      |
| GET    | /items       | List items        |
| GET    | /items/:id   | Get item by ID    |

## ⚙️ Configuration

Configuration is loaded from:
1. `config/config.yaml` (base config)
2. `config/config-{API_ENVIRONMENT}.yaml` (environment override)
3. Environment variables with `APP_` prefix

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_ENVIRONMENT` | Environment name | development |
| `APP_HTTP__PORT` | HTTP port | 8080 |
| `APP_HTTP__DEBUG` | Debug mode | true |

## 🐳 Docker

```bash
# Build
make docker-build

# Run
docker run -p 8080:8080 -e API_ENVIRONMENT=production {{SERVICE_NAME}}:latest
```

## 📊 Observability

The service includes built-in support for:

- **Metrics**: Prometheus metrics at `/metrics`
- **Tracing**: OpenTelemetry with configurable exporters
- **Logging**: Structured logging with zap

## 📝 License

[Add your license here]