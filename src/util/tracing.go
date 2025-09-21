package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer initialies an OLTP tracer. Call this once at program startup.
func InitTracer(ctx context.Context, serviceName, deploymentEnvironment string) func(context.Context) error {
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

	tracer = otel.Tracer(serviceName)
	return tp.Shutdown
}

// NewSpan creates a new span specifically for use inside a Gin handler
// function.
//
// c is the context. name should be a descriptive operation like "create-user"
// or "fetch-all-requests" or "query-db-for-rooms".
//
// Function returns a context (used when submitting HTTP requests yourself, see
// InjectSpan) and a span object (used in AddEvent) which you should close.
func NewSpan(c *gin.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return tracer.Start(c.Request.Context(), name, trace.WithAttributes(attrs...))
}

// SetSpanUser adds user context to the given span.
//
// By default, spans created with NewSpan don't have a user.id set assuming the
// user is not specified or the method is anonymous. You can add further user
// context with this method.
func SetSpanUser(span trace.Span, userId int) {
	span.SetAttributes(attribute.String("user.id", fmt.Sprintf("%d", userId)))
}

// InjectSpan injects a tracere span inside an outgoing http request.
//
// Use this when sending requests to other microservices.
func InjectSpan(tracerCtx context.Context, outgoingRequest *http.Request) {
	otel.GetTextMapPropagator().Inject(tracerCtx, propagation.HeaderCarrier(outgoingRequest.Header))
}

// AddEvent submits an event to the given span.
//
// It works like logging but is not a supstitution for logging.
func AddEvent(span trace.Span, msg string, err error) {
	if err == nil {
		span.AddEvent(msg)
	} else {
		span.AddEvent(msg, trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
	}
}
