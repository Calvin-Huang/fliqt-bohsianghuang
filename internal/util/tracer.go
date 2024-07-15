package util

import (
	"context"
	"fliqt/config"
	"runtime"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const serviceName = "fliqt_API"

func getCallerStack(stack int) (string, string, int, bool) {
	pc, file, line, ok := runtime.Caller(stack + 1)
	if !ok {
		return "", "", 0, ok
	}

	callStack := runtime.FuncForPC(pc).Name()

	return callStack, file, line, ok
}

func GetFileNameFromCaller() string {
	callStack, file, _, ok := getCallerStack(1)
	if !ok {
		return ""
	}

	rootPackage := strings.Split(callStack, "/")[0]
	relativeFileStartIdx := strings.LastIndex(file, rootPackage)
	if relativeFileStartIdx == -1 {
		// Workaround for containerized environment, we use "src" as the root package
		relativeFileStartIdx = strings.LastIndex(file, "src")
	}
	relativeFilePath := file[relativeFileStartIdx:]

	return relativeFilePath
}

func GetSpanNameFromCaller() string {
	callStack, _, _, ok := getCallerStack(1)
	if !ok {
		return ""
	}

	stacks := strings.Split(callStack, "/")

	return stacks[len(stacks)-1]
}

func getRootPackageFromCaller() string {
	callStack, _, _, ok := getCallerStack(1)
	if !ok {
		return ""
	}

	rootPackage := strings.Split(callStack, "/")[0]

	return rootPackage
}

// Set global http trace provider
func InitTracer(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	// Create the Jaeger exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.TracerEndpoint),
		otlptracegrpc.WithDialOption(
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.6,
					MaxDelay:   15 * time.Second,
				},
				MinConnectTimeout: 0,
			}),
		),
	)
	if err != nil {
		return err
	}
	appTracerProvider := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(exporter)),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("exporter", getRootPackageFromCaller()),
		)),
	)
	otel.SetTracerProvider(appTracerProvider)

	return nil
}
