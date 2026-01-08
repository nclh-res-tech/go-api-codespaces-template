package httpclient

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Config configures the outbound HTTP client.
type Config struct {
	BaseURL   string            `koanf:"base_url"`
	Timeout   time.Duration     `koanf:"timeout"`
	Headers   map[string]string `koanf:"headers"`
	UserAgent string            `koanf:"user_agent"`
}

// New constructs a Resty client with sensible defaults and tracing instrumentation.
func New(cfg Config) *resty.Client {
	// Create HTTP transport with OpenTelemetry instrumentation
	transport := otelhttp.NewTransport(http.DefaultTransport)

	client := resty.NewWithClient(&http.Client{
		Transport: transport,
	})

	if cfg.BaseURL != "" {
		client.SetBaseURL(cfg.BaseURL)
	}
	if cfg.Timeout > 0 {
		client.SetTimeout(cfg.Timeout)
	}
	if cfg.UserAgent != "" {
		client.SetHeader("User-Agent", cfg.UserAgent)
	}
	for k, v := range cfg.Headers {
		client.SetHeader(k, v)
	}

	return client
}
