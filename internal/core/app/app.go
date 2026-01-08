package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

// Listener represents something that can listen (e.g., HTTP server).
type Listener interface {
	Listen() error
	Shutdown(ctx context.Context) error
}

// App manages the application lifecycle.
type App struct {
	logger      *zap.Logger
	shutdownFns []func()
	mu          sync.Mutex
}

// New creates a new App instance.
func New(logger *zap.Logger) *App {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &App{logger: logger}
}

// OnShutdown registers a function to be called during shutdown.
func (a *App) OnShutdown(fn func()) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shutdownFns = append(a.shutdownFns, fn)
}

// Start runs the application with the provided startup function.
func Start(startFn func(ctx context.Context, app *App) ([]Listener, error)) {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := New(logger)

	listeners, err := startFn(ctx, app)
	if err != nil {
		logger.Fatal("failed to start application", zap.Error(err))
	}

	// Start all listeners
	var wg sync.WaitGroup
	for _, l := range listeners {
		wg.Add(1)
		go func(listener Listener) {
			defer wg.Done()
			if err := listener.Listen(); err != nil {
				logger.Error("listener error", zap.Error(err))
			}
		}(l)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down...")
	cancel()

	// Shutdown all listeners
	shutdownCtx := context.Background()
	for _, l := range listeners {
		if err := l.Shutdown(shutdownCtx); err != nil {
			logger.Error("listener shutdown error", zap.Error(err))
		}
	}

	// Run shutdown functions
	for _, fn := range app.shutdownFns {
		fn()
	}

	wg.Wait()
	logger.Info("shutdown complete")
}
