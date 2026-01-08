package observability

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Config represents telemetry configuration for traces and metrics.
type Config struct {
	ServiceName string        `koanf:"service_name" validate:"required"`
	Traces      []TraceConfig `koanf:"traces"`
	Metrics     MetricConfig  `koanf:"metrics"`
}

// TraceConfig allows selecting the trace exporter.
type TraceConfig struct {
	Exporter string `koanf:"exporter"` // stdout, otlphttp, otlpgrpc, none
	Endpoint string `koanf:"endpoint"`
	Insecure bool   `koanf:"insecure"`
}

// MetricConfig allows selecting the metrics exporter.
type MetricConfig struct {
	Exporter string `koanf:"exporter"` // prometheus, none
}

// Providers bundles telemetry providers and handlers.
type Providers struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	MetricsHandler http.Handler
	Shutdown       func(context.Context) error
}

// Setup configures global OpenTelemetry providers.
func Setup(ctx context.Context, cfg Config) (*Providers, error) {
	base := resource.Default()
	res, err := resource.Merge(
		base,
		resource.NewWithAttributes(
			base.SchemaURL(),
			attribute.String("service.name", cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tb, err := buildTracerProvider(ctx, cfg, res)
	if err != nil {
		return nil, err
	}

	mp, metricsHandler, err := buildMeterProvider(cfg, res)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tb)
	otel.SetMeterProvider(mp)

	shutdown := func(ctx context.Context) error {
		if err := mp.Shutdown(ctx); err != nil {
			return err
		}
		return tb.Shutdown(ctx)
	}

	return &Providers{
		TracerProvider: tb,
		MeterProvider:  mp,
		MetricsHandler: metricsHandler,
		Shutdown:       shutdown,
	}, nil
}

func buildTracerProvider(ctx context.Context, cfg Config, res *resource.Resource) (*trace.TracerProvider, error) {
	var processors []trace.SpanProcessor

	for _, tc := range cfg.Traces {
		if tc.Exporter == "" || tc.Exporter == "none" {
			continue
		}

		exp, err := buildSpanExporter(ctx, tc)
		if err != nil {
			return nil, err
		}
		if exp != nil {
			processors = append(processors, trace.NewBatchSpanProcessor(exp))
		}
	}

	opts := []trace.TracerProviderOption{trace.WithResource(res)}
	for _, p := range processors {
		opts = append(opts, trace.WithSpanProcessor(p))
	}

	return trace.NewTracerProvider(opts...), nil
}

func buildSpanExporter(ctx context.Context, tc TraceConfig) (trace.SpanExporter, error) {
	switch tc.Exporter {
	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	case "otlphttp":
		opts := []otlptracehttp.Option{}
		if tc.Endpoint != "" {
			opts = append(opts, otlptracehttp.WithEndpoint(tc.Endpoint))
		}
		if tc.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, opts...)
	case "otlpgrpc":
		opts := []otlptracegrpc.Option{}
		if tc.Endpoint != "" {
			opts = append(opts, otlptracegrpc.WithEndpoint(tc.Endpoint))
		}
		if tc.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		return otlptracegrpc.New(ctx, opts...)
	default:
		return nil, errors.New("unsupported trace exporter")
	}
}

func buildMeterProvider(cfg Config, res *resource.Resource) (*metric.MeterProvider, http.Handler, error) {
	switch cfg.Metrics.Exporter {
	case "", "prometheus":
		exp, err := prometheus.New(prometheus.WithoutUnits())
		if err != nil {
			return nil, nil, err
		}
		mp := metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(exp),
		)
		return mp, promhttp.Handler(), nil
	case "none":
		mp := metric.NewMeterProvider(metric.WithResource(res))
		return mp, nil, nil
	default:
		return nil, nil, errors.New("unsupported metrics exporter")
	}
}
