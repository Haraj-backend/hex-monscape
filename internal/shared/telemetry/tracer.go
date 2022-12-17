package telemetry

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/validator.v2"
)

type Tracer interface {
	Trace(context context.Context, name string) Spanner
}

type Spanner interface {
	End(options ...trace.SpanEndOption)
}

type OpenTelemetryConfig struct {
	BaseContext context.Context
	Exporter    sdktrace.TracerProvider `validate:"nonnil"`
	ServiceName string                  `validate:"nonzero"`
}

type OpenTelemetryTracer struct {
	baseCtx  context.Context
	Tracer   trace.Tracer
	Exporter sdktrace.TracerProvider
}

func NewOpenTelemetryTracer(config OpenTelemetryConfig) (Tracer, error) {
	err := validator.Validate(config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	svcName := config.ServiceName
	exporter := config.Exporter

	// set tracer provider and propagator properly
	// this is to ensure all instrumentation library could run well
	otel.SetTracerProvider(&exporter)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer := exporter.Tracer(svcName)

	return &OpenTelemetryTracer{
		baseCtx:  config.BaseContext,
		Tracer:   tracer,
		Exporter: config.Exporter,
	}, nil
}

func NewJaegerTracerProvider(url string, svcName string) (*sdktrace.TracerProvider, error) {
	// create the Jaeger exporter
	traceExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, fmt.Errorf("unable to create jaeger exporter: %w", err)
	}

	// initialize tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		)),
	)
	return tp, nil
}

func (t *OpenTelemetryTracer) Trace(ctx context.Context, name string) Spanner {
	_, span := t.Tracer.Start(ctx, name)
	return span
}

// dummy tracer
type InitialTracer struct{}

func NewInitialTracer() (Tracer, error) {
	return &InitialTracer{}, nil
}

type InitialSpanner struct{}

func (s *InitialSpanner) End(options ...trace.SpanEndOption) {
	// do nothing
}

func (t *InitialTracer) Trace(ctx context.Context, name string) Spanner {
	log.Printf("trace name: %s", name)
	return &InitialSpanner{}
}

var tracer *Tracer

func SetTracer(tcr *Tracer) {
	tracer = tcr
}

func GetTracer() Tracer {
	if tracer == nil {
		tracer, _ := NewInitialTracer()
		SetTracer(&tracer)
	}

	return *tracer
}
