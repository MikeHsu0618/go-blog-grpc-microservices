//go:build wireinject
// +build wireinject

package main

import (
	v1 "blog-grpc-microservices/api/protobuf/post/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/dbcontext"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/post"
	"github.com/google/wire"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.PostServiceServer, error) {
	wire.Build(
		dbcontext.NewPostDB,
		post.NewRepository,
		post.NewServer,
	)

	return &post.Server{}, nil
}
