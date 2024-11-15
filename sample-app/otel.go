package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func iniOtelResource() (*resource.Resource, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("build.compiler", Compiler),
			attribute.String("build.go_version", GoVersion),
			attribute.String("build.platform", Platform),
			semconv.ServiceName("otel-keda-example"),
		),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func initMeterProvider(resource *resource.Resource, ctx context.Context) (*sdkmetric.MeterProvider, error) {
	stdOutExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}
	grpcReader, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(stdOutExporter, sdkmetric.WithInterval(10*time.Second))),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(grpcReader, sdkmetric.WithInterval(10*time.Second))),
	)

	otel.SetMeterProvider(provider)
	return provider, nil
}
