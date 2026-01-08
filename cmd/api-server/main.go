package main

import (
	"context"
	"fmt"
	"net/http"

	"{{MODULE_PATH}}/common/httpserver"
	"{{MODULE_PATH}}/common/observability"
	"{{MODULE_PATH}}/internal/config"
	"{{MODULE_PATH}}/internal/core/app"
	"{{MODULE_PATH}}/internal/routes"
	"{{MODULE_PATH}}/internal/services"
	"{{MODULE_PATH}}/internal/stores"

	"go.uber.org/zap"
)

func main() {
	app.Start(run)
}

func run(ctx context.Context, application *app.App) ([]app.Listener, error) {
	// Load configuration
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Set up logging
	logger, _ := zap.NewProduction()
	if cfg.HTTP.Debug {
		logger, _ = zap.NewDevelopment()
	}

	// Set up telemetry
	providers, err := observability.Setup(ctx, cfg.Telemetry)
	if err != nil {
		logger.Warn("failed to set up telemetry", zap.Error(err))
	}
	if providers != nil {
		application.OnShutdown(func() {
			_ = providers.Shutdown(ctx)
		})
	}

	// Initialize stores
	itemStore := stores.NewItemStore()

	// Initialize services
	itemService := services.NewItemService(itemStore)

	// Build routes
	routeList := routes.BuildRoutes(routes.RouteDependencies{
		ItemService: itemService,
	}, logger)

	// Set up OpenAPI documentation
	openAPI := httpserver.NewOpenAPI(logger, "{{SERVICE_NAME}} API", "1.0.0", "API documentation")

	// Create HTTP server
	var serverOpts []func(*httpserver.Server)
	serverOpts = append(serverOpts,
		httpserver.WithLogger(logger),
		httpserver.WithServiceName(cfg.Telemetry.ServiceName),
		httpserver.WithGinDebug(cfg.HTTP.Debug),
		httpserver.WithOpenAPI(openAPI),
	)
	if providers != nil && providers.MetricsHandler != nil {
		serverOpts = append(serverOpts, httpserver.WithMetricsHandler(providers.MetricsHandler))
	}

	server := httpserver.New(routeList, serverOpts...)

	// Create listener
	listener := &httpListener{
		addr:   ":" + cfg.HTTP.Port,
		server: server,
		logger: logger,
	}

	logger.Info("starting server", zap.String("port", cfg.HTTP.Port))
	return []app.Listener{listener}, nil
}

type httpListener struct {
	addr    string
	server  *httpserver.Server
	logger  *zap.Logger
	httpSrv *http.Server
}

func (l *httpListener) Listen() error {
	router := l.server.Router()
	l.httpSrv = &http.Server{
		Addr:    l.addr,
		Handler: router,
	}
	l.logger.Info("HTTP server listening", zap.String("addr", l.addr))
	return l.httpSrv.ListenAndServe()
}

func (l *httpListener) Shutdown(ctx context.Context) error {
	if l.httpSrv != nil {
		return l.httpSrv.Shutdown(ctx)
	}
	return nil
}
