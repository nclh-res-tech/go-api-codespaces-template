package observability

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Instrumentation exposes convenience helpers for emitting traces and metrics.
// It uses the global OTel providers set up in telemetry.go.
type Instrumentation struct {
	tracer trace.Tracer
	meter  metric.Meter

	counters   sync.Map // name -> metric.Int64Counter
	histograms sync.Map // name -> metric.Float64Histogram
}

// NewInstrumentation creates a helper for a given instrumentation scope.
// The name should be your service/component name.
func NewInstrumentation(name string) *Instrumentation {
	return &Instrumentation{
		tracer: otel.Tracer(name),
		meter:  otel.Meter(name),
	}
}

// AddEvent records a trace event on the current span, if any.
func (i *Instrumentation) AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// StartSpan starts a span and returns the new context and span.
func (i *Instrumentation) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return i.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// RecordCounter increments (or decrements, with negative value) a named counter.
func (i *Instrumentation) RecordCounter(ctx context.Context, name string, value int64, attrs ...attribute.KeyValue) {
	counter := i.getCounter(name)
	counter.Add(ctx, value, metric.WithAttributes(attrs...))
}

// RecordHistogram records a measurement to a named histogram (float64).
func (i *Instrumentation) RecordHistogram(ctx context.Context, name string, value float64, attrs ...attribute.KeyValue) {
	h := i.getHistogram(name)
	h.Record(ctx, value, metric.WithAttributes(attrs...))
}

func (i *Instrumentation) getCounter(name string) metric.Int64Counter {
	if c, ok := i.counters.Load(name); ok {
		return c.(metric.Int64Counter)
	}
	counter, _ := i.meter.Int64Counter(name)
	i.counters.Store(name, counter)
	return counter
}

func (i *Instrumentation) getHistogram(name string) metric.Float64Histogram {
	if h, ok := i.histograms.Load(name); ok {
		return h.(metric.Float64Histogram)
	}
	hist, _ := i.meter.Float64Histogram(name)
	i.histograms.Store(name, hist)
	return hist
}
