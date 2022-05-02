// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"blog-grpc-microservices/api/protobuf/user/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/dbcontext"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/user"
)

// Injectors from wire.go:

func InitServer(logger *log.Logger, conf *config.Config) (v1.UserServiceServer, error) {
	db, err := dbcontext.NewUserDB(conf)
	if err != nil {
		return nil, err
	}
	repository := user.NewRepository(logger, db)
	userServiceServer := user.NewServer(logger, repository)
	return userServiceServer, nil
}
