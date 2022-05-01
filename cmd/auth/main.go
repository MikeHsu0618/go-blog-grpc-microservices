package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "blog-grpc-microservices/api/protobuf/auth/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
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

	// Start gRPC server
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	if err = metricsServer.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	<-ctx.Done()
	close(ch)
	logger.Info("Graceful Shutdown end")
}
