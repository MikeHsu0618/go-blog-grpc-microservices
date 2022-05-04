package comment

import (
	"context"
	"time"

	v1 "blog-grpc-microservices/api/protobuf/comment/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(logger *log.Logger, conf *config.Config) (v1.CommentServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.Comment.Server.Host+conf.Comment.Server.GRPC.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := v1.NewCommentServiceClient(conn)
	return client, nil
}
