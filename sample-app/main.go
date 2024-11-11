package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/joaoestrela/otel-keda-example/generated/counter/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type server struct {
	pb.UnimplementedCounterServiceServer
	counter       int32
	counterMetric metric.Int64UpDownCounter
}

func (s *server) IncreaseCounter(ctx context.Context, req *pb.CounterRequest) (*pb.CounterResponse, error) {
	s.counter += req.Value
	s.counterMetric.Add(ctx, int64(req.Value), metric.WithAttributes(attribute.String("operation", "increase")))
	return &pb.CounterResponse{Result: s.counter}, nil
}

func (s *server) DecreaseCounter(ctx context.Context, req *pb.CounterRequest) (*pb.CounterResponse, error) {
	s.counter -= req.Value
	s.counterMetric.Add(ctx, int64(-req.Value), metric.WithAttributes(attribute.String("operation", "decrease")))
	return &pb.CounterResponse{Result: s.counter}, nil
}

func main() {
	addr := flag.String("server-addr", ":8080", "The server addr")
	flag.Parse()
	ctx := context.Background()
	otelResource, err := iniOtelResource()
	if err != nil {
		log.Fatal(err)
	}

	provider, err := initMeter(otelResource, ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}()

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}

	meter := provider.Meter("counter-server")

	counterMetric, err := meter.Int64UpDownCounter("counter", metric.WithDescription("A counter that can go up and down"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	pb.RegisterCounterServiceServer(s, &server{counterMetric: counterMetric})
	reflection.Register(s)

	log.Printf("Server is running on %s", *addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
