package shared

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, serviceName, serviceVersion string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Setup resource.
	res, err := newResource(serviceName, serviceVersion)
	if err != nil {
		handleErr(err)
		return
	}

	// Setup trace provider.
	tracerProvider, err := newTraceProvider(res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Setup meter provider.
	meterProvider, err := newMeterProvider(res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(logger),
	)
	// traceExporter, err := stdouttrace.New(
	// 	stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	batcher := trace.NewBatchSpanProcessor(traceExporter)

	traceProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(batcher),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("zipkin-test"),
		)),
	)
	// otel.SetTracerProvider(tp)

	// traceProvider := trace.NewTracerProvider(
	// 	trace.WithBatcher(traceExporter,
	// 		// Default is 5s. Set to 1s for demonstrative purposes.
	// 		trace.WithBatchTimeout(5*time.Second)),
	// 	trace.WithResource(res),
	// )
	return traceProvider, nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(time.Minute))),
	)
	return meterProvider, nil
}
