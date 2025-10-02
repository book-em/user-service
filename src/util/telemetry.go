package utils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	slogmulti "github.com/samber/slog-multi"
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
	// During tests, the tracer is not set up, so we silently ignore tracing.
	// This has to be done manually (i.e. don't call tracer methods, don't touch
	// spans etc.)
	tracerReady bool
	Tracer      trace.Tracer

	SpanStack []SpanPair

	loggerReady bool
	logger      *slog.Logger
}

var TEL Telemetry

func (t *Telemetry) Init(ctx context.Context, serviceName, deploymentEnvironment string) func(context.Context) error {
	// [0] Init logger
	{
		t.logger = slog.New(
			slogmulti.Fanout(
				// TODO: Save to file and then use promtail, use JSON handler there.
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		)

		t.loggerReady = true
	}
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
		t.tracerReady = true
		return tp.Shutdown
	}
}

func (t *Telemetry) Push(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	if t.tracerReady {
		newCtx, span := t.Tracer.Start(ctx, name, trace.WithAttributes(attrs...))
		t.SpanStack = append(t.SpanStack, SpanPair{Ctx: newCtx, Span: span})
	} else {
		newCtx := ctx
		var span trace.Span
		t.SpanStack = append(t.SpanStack, SpanPair{Ctx: newCtx, Span: span})
	}
}

func (t *Telemetry) Pop() {
	top := t.SpanStack[len(t.SpanStack)-1]
	if t.tracerReady {
		top.Span.End()
	}
	t.SpanStack = t.SpanStack[:len(t.SpanStack)-1]
}

func (t *Telemetry) Top() SpanPair {
	return t.SpanStack[len(t.SpanStack)-1]
}

func (t *Telemetry) Ctx() context.Context {
	return t.Top().Ctx
}

// func (t *Telemetry) Event(msg string, err error) {
// 	// Logging
// 	{
// 		if err == nil {
// 			t.logger.Info(msg)
// 		} else {
// 			t.logger.Error(msg, "error", err)
// 		}
// 	}

// 	// Tracing
// 	if t.tracerReady && len(t.SpanStack) > 0 {
// 		span := t.SpanStack[len(t.SpanStack)-1].Span

// 		if err == nil {
// 			span.AddEvent(msg)
// 		} else {
// 			span.AddEvent(msg, trace.WithAttributes(
// 				attribute.String("error.message", err.Error()),
// 				attribute.Bool("error", true),
// 			))
// 			span.SetStatus(codes.Error, err.Error())
// 		}
// 	}
// }

// func (t *Telemetry) Eventf(msg string, err error, a ...any) {
// 	msgFinal := fmt.Sprintf(msg, a...)
// 	t.Event(msgFinal, err)
// }

func (t *Telemetry) SetAttrib(kv ...attribute.KeyValue) {
	if t.tracerReady {
		t.Top().Span.SetAttributes(kv...)
	}
}

func (t *Telemetry) SetUser(id uint) {
	if t.tracerReady {
		t.Top().Span.SetAttributes(attribute.String("user.id", fmt.Sprintf("%d", id)))
	}
}

func (t *Telemetry) Inject(outgoingRequest *http.Request) {
	otel.GetTextMapPropagator().Inject(t.Ctx(), propagation.HeaderCarrier(outgoingRequest.Header))
}

func (t *Telemetry) Info(msg string, attrs ...any) {
	if t.loggerReady {
		t.logger.Info(msg, attrs...)
	}
	if span := t.currentSpan(); span != nil {
		span.AddEvent(msg)
	}
}

func (t *Telemetry) Warn(msg string, attrs ...any) {
	if t.loggerReady {
		t.logger.Warn(msg, attrs...)
	}
	if span := t.currentSpan(); span != nil {
		span.AddEvent(msg)
	}
}

func (t *Telemetry) Debug(msg string, attrs ...any) {
	if t.loggerReady {
		t.logger.Debug(msg, attrs...)
	}
	if span := t.currentSpan(); span != nil {
		span.AddEvent(msg)
	}
}

func (t *Telemetry) Error(msg string, err error, attrs ...any) {
	if t.loggerReady {
		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
		}
		t.logger.Error(msg, attrs...)
	}

	if span := t.currentSpan(); span != nil {
		span.AddEvent(msg, trace.WithAttributes(
			attribute.String("error.message", err.Error()),
			attribute.Bool("error", true),
		))
		span.SetStatus(codes.Error, err.Error())
	}
}

func (t *Telemetry) currentSpan() trace.Span {
	if t.tracerReady && len(t.SpanStack) > 0 {
		return t.SpanStack[len(t.SpanStack)-1].Span
	}
	return nil
}
