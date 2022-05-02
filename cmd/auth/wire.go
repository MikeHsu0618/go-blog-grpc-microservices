//go:build wireinject
// +build wireinject

package main

import (
	"blog-grpc-microservices/api/protobuf/auth/v1"
	"blog-grpc-microservices/internal/auth"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/jwt"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/google/wire"
)

func InitServer(logger *log.Logger, conf *config.Config) (v1.AuthServiceServer, error) {
	wire.Build(
		jwt.NewManager,
		auth.NewServer,
	)
	return &auth.Server{}, nil
}
