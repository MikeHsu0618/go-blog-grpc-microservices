//go:build wireinject
// +build wireinject

package main

import (
	v1 "blog-grpc-microservices/api/protobuf/user/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/dbcontext"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/user"
	"github.com/google/wire"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.UserServiceServer, error) {
	wire.Build(
		dbcontext.NewUserDB,
		user.NewRepository,
		user.NewServer,
	)
	return &user.Server{}, nil
}
