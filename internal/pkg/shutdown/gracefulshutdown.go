package shutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc"
)

func GracefulShutDown(grpcServer *grpc.Server, metricsServer *http.Server, logger *log.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	<-ctx.Done()
	close(ch)
	logger.Info("Graceful Shutdown end")
}
