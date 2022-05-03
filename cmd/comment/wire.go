//go:build wireinject
// +build wireinject

package main

import (
	v1 "blog-grpc-microservices/api/protobuf/comment/v1"
	"blog-grpc-microservices/internal/comment"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/dbcontext"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/google/wire"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.CommentServiceServer, error) {
	wire.Build(
		dbcontext.NewCommentDB,
		comment.NewRepository,
		comment.NewServer,
	)
	return &comment.Server{}, nil
}
