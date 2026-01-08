package httpserver

import (
	"{{MODULE_PATH}}/common/httpserver"
)

// New creates an HTTP server with application defaults.
func New(routes []httpserver.Route, opts ...func(*httpserver.Server)) *httpserver.Server {
	return httpserver.New(routes, opts...)
}

// Config holds HTTP server configuration.
type Config struct {
	Port  string `koanf:"port"`
	Debug bool   `koanf:"debug"`
}
