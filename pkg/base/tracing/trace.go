package tracing

import (
	"context"
	"github.com/coreos/rkt/tests/testutils/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"io"
	"os"
)

func NewTracer(serviceName string) func() {
	// tracing
	f, err := os.Create("/tmp/traces.txt")
	if err != nil {
		logger.Fatal(err)
	}
	exp, err := newExporter(f)
	if err != nil {
		logger.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource(serviceName)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatal(err)
		}
		if err := f.Close(); err != nil {
			logger.Fatal(err)
		}
	}
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		//stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource(serviceName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("v1.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}