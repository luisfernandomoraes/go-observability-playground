package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/luisfernandomoraes/observability-go/internal/handlers/users"
	"github.com/luisfernandomoraes/observability-go/internal/utils/logger"

	_ "net/http/pprof"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const serviceName = "transform-user-service"

func initTracing(ctx context.Context) {
	// Configurar propagador global
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// initTracerProvider initializes the OpenTelemetry TracerProvider
func initTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, func(), error) {
	// Use otlptracegrpc.NewClient to configure the trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown TracerProvider", "error", err)
		}
	}, nil
}

// initMeterProvider initializes the OpenTelemetry MeterProvider
func initMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, func(), error) {
	// Use otlpmetricgrpc.NewClient to configure the metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider, func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown MeterProvider", "error", err)
		}
	}, nil
}

// Middleware personalizado para métricas
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Adicionar atributos comuns
		attrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		}

		meter := otel.GetMeterProvider().Meter(serviceName)
		counter, _ := meter.Int64Counter("http.requests.total")
		histogram, _ := meter.Float64Histogram("http.request.duration")

		defer func() {
			duration := time.Since(start).Milliseconds()
			counter.Add(r.Context(), 1, metric.WithAttributes(attrs...))
			histogram.Record(r.Context(), float64(duration), metric.WithAttributes(attrs...))
		}()

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Habilitar block profiling
	runtime.SetBlockProfileRate(1)

	// Habilitar mutex profiling
	runtime.SetMutexProfileFraction(1)

	// Initialize logger
	logger.SetupLogger()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize TracerProvider using NewClient
	_, shutdownTracer, err := initTracerProvider(ctx)
	if err != nil {
		slog.Error("Failed to initialize TracerProvider", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer()

	// Initialize MeterProvider using NewClient
	_, shutdownMeter, err := initMeterProvider(ctx)
	if err != nil {
		slog.Error("Failed to initialize MeterProvider", "error", err)
		os.Exit(1)
	}
	defer shutdownMeter()

	initTracing(ctx)

	// Add OpenTelemetry instrumentation to HTTP server
	mux := http.NewServeMux()
	handler := metricsMiddleware(
		otelhttp.NewHandler(
			http.HandlerFunc(users.TransformUsersHandler),
			"POST /users/transform",
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			}),
		),
	)
	mux.Handle("/users/transform", handler)

	// Configurar número máximo de CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Criar mux para debug endpoints
	debugMux := http.NewServeMux()
	// O import _ "net/http/pprof" registra automaticamente todos os handlers do pprof
	// no DefaultServeMux. Aqui estamos copiando esses handlers para nosso debugMux
	debugMux.Handle("/debug/pprof/", http.DefaultServeMux)

	// Iniciar servidor de debug em uma porta separada
	go func() {
		debugAddr := ":6060"
		slog.Info("Starting debug server", "addr", debugAddr)
		if err := http.ListenAndServe(debugAddr, debugMux); err != nil {
			slog.Error("Debug server failed", "error", err)
		}
	}()

	// Start HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Listen for OS signals to shutdown
	<-ctx.Done()
	slog.Info("Shutting down server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
}
