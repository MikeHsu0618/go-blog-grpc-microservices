//go:build wireinject
// +build wireinject

package main

import (
	v1 "blog-grpc-microservices/api/protobuf/blog/v1"
	"blog-grpc-microservices/internal/auth"
	"blog-grpc-microservices/internal/blog"
	"blog-grpc-microservices/internal/comment"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/post"
	"blog-grpc-microservices/internal/user"
	"github.com/google/wire"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.BlogServiceServer, error) {
	wire.Build(
		user.NewClient,
		auth.NewClient,
		post.NewClient,
		comment.NewClient,
		blog.NewServer,
	)

	return &blog.Server{}, nil
}
