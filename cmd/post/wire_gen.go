// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"blog-grpc-microservices/api/protobuf/post/v1"
	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/dbcontext"
	"blog-grpc-microservices/internal/pkg/log"
	"blog-grpc-microservices/internal/post"
)

// Injectors from wire.go:

func InitServer(logger *log.Logger, conf *config.Config) (v1.PostServiceServer, error) {
	db, err := dbcontext.NewPostDB(conf)
	if err != nil {
		return nil, err
	}
	repository := post.NewRepository(logger, db)
	postServiceServer := post.NewServer(logger, repository)
	return postServiceServer, nil
}
