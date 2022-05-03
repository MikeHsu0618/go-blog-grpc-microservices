package main

import (
	"flag"
	"net"
	"net/http"

	v1 "blog-grpc-microservices/api/protobuf/auth/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/pkg/shutdown"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var flagConfig = flag.String("config", "./configs/config.yaml", "path to config file")

func main() {
	flag.Parse()
	logger := log.New()

	conf, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal(err)
	}

	// grpc server
	authServer, err := InitServer(logger, conf)
	if err != nil {
		logger.Fatal(err)
	}
	// grpc health check
	healthServer := health.NewServer()

	// grpc middleware
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_validator.UnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(),
	)))

	v1.RegisterAuthServiceServer(grpcServer, authServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", conf.Auth.Server.GRPC.Port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	// Start grpc server
	logger.Infof("gPRC Listening on port %s", conf.Auth.Server.GRPC.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatal(err)
		}
	}()

	// Start Metrics server
	logger.Infof("Metrics Listening on port %s", conf.Auth.Server.Metrics.Port)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    conf.Auth.Server.Metrics.Port,
		Handler: metricsMux,
	}
	go func() {
		if err = metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	shutdown.GracefulShutDown(grpcServer, metricsServer, logger)
}
