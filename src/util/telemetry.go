package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type SpanPair struct {
	Ctx  context.Context
	Span trace.Span
}

type Telemetry struct {
	Tracer trace.Tracer

	SpanStack []SpanPair
}

var TEL Telemetry

func (t *Telemetry) Init(ctx context.Context, serviceName, deploymentEnvironment string) func(context.Context) error {
	// [1] Init tracer
	{
		exp, err := otlptracehttp.New(ctx)
		if err != nil {
			log.Fatalf("Failed to create an OTLP HTTP exporter: %v", err)
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
				semconv.DeploymentEnvironment(deploymentEnvironment),
			)),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})

		t.Tracer = otel.Tracer(serviceName)
		return tp.Shutdown
	}
}

func (t *Telemetry) Push(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	newCtx, span := t.Tracer.Start(ctx, name, trace.WithAttributes(attrs...))

	t.SpanStack = append(t.SpanStack, SpanPair{Ctx: newCtx, Span: span})
}

func (t *Telemetry) Pop() {
	top := t.SpanStack[len(t.SpanStack)-1]
	top.Span.End()
	t.SpanStack = t.SpanStack[:len(t.SpanStack)-1]
}

func (t *Telemetry) Top() SpanPair {
	return t.SpanStack[len(t.SpanStack)-1]
}

func (t *Telemetry) Ctx() context.Context {
	return t.Top().Ctx
}

func (t *Telemetry) Event(msg string, err error) {
	// Logging
	{
		if err == nil {
			log.Println(msg)
		} else {
			log.Printf("%s: %v\n", msg, err)
		}
	}

	// Tracing
	if len(t.SpanStack) > 0 {
		span := t.SpanStack[len(t.SpanStack)-1].Span

		if err == nil {
			span.AddEvent(msg)
		} else {
			span.AddEvent(msg, trace.WithAttributes(
				attribute.String("error.message", err.Error()),
				attribute.Bool("error", true),
			))
			span.SetStatus(codes.Error, err.Error())
		}
	}
}

func (t *Telemetry) SetAttrib(kv ...attribute.KeyValue) {
	t.Top().Span.SetAttributes(kv...)
}

func (t *Telemetry) SetUser(id uint) {
	t.Top().Span.SetAttributes(attribute.String("user.id", fmt.Sprintf("%d", id)))
}

func (t *Telemetry) Inject(outgoingRequest *http.Request) {
	otel.GetTextMapPropagator().Inject(t.Ctx(), propagation.HeaderCarrier(outgoingRequest.Header))
}
