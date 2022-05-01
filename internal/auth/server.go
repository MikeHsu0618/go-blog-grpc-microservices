package auth

import (
	"context"

	"blog-grpc-microservices/api/protobuf/auth/v1"
	"blog-grpc-microservices/internal/pkg/jwt"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer(logger *log.Logger, jwtManager *jwt.Manager) v1.AuthServiceServer {
	return &Server{
		logger:     logger,
		jwtManager: jwtManager,
	}
}

type Server struct {
	v1.UnimplementedAuthServiceServer
	logger     *log.Logger
	jwtManager *jwt.Manager
}

func (s Server) GenerateToken(ctx context.Context, req *v1.GenerateTokenRequest) (*v1.GenerateTokenResponse, error) {
	token, err := s.jwtManager.Generate(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &v1.GenerateTokenResponse{
		Token: token,
	}, nil
}

func (s Server) ValidateToken(ctx context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	claims, err := s.jwtManager.Validate(req.GetToken())
	if claims.UserId == 0 || err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")

	}
	return &v1.ValidateTokenResponse{
		Valid: true,
	}, nil
}

func (s Server) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	claims, err := s.jwtManager.Validate(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	token, err := s.jwtManager.Generate(claims.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &v1.RefreshTokenResponse{
		Token: token,
	}, nil
}
