package post

import (
	"context"
	"time"

	v1 "blog-grpc-microservices/api/protobuf/post/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(logger *log.Logger, conf *config.Config) (v1.PostServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Post.Server.Host+conf.Post.Server.GRPC.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := v1.NewPostServiceClient(conn)
	return client, nil
}
