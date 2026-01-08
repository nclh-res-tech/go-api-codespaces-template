package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"{{MODULE_PATH}}/common/errors"
)

const (
	// ErrServer is the error returned when the server stops due to an error.
	ErrServer = errors.Error("listen stopped with error")
)

const (
	readHeaderTimeout = 60 * time.Second
)

// Config represents the configuration of the http listener.
type Config struct {
	Port  string `yaml:"port" koanf:"port" validate:"required"`
	Debug bool   `yaml:"debug" koanf:"debug"`
}

// Server represents a http server that listens on a port.
type Server struct {
	server *http.Server
	port   string
}

// NewListener instantiates a new instance of Server.
func NewListener(handler http.Handler, cfg Config) (*Server, error) {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return With(baseContext, From(baseContext))
			},
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		port: cfg.Port,
	}, nil
}

// Listen starts the server and listens on the configured port.
func (s *Server) Listen(ctx context.Context) error {
	From(ctx).Info(fmt.Sprintf("http server starting on port: %s", s.port))

	err := s.server.ListenAndServe()
	if err != nil {
		return ErrServer.Wrap(err)
	}

	From(ctx).Info("http server stopped")

	return nil
}
