package interceptor

import (
	"context"
	"strings"

	"blog-grpc-microservices/internal/pkg/jwt"
	"blog-grpc-microservices/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

var (
	expectedScheme  = "bearer"
	headerAuthorize = "authorization"
	ContextKeyID    = contextKey("ID")
)

func NewAuthInterceptor(logger *log.Logger, jwtManager *jwt.Manager, authMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{
		logger:      logger,
		jwtManager:  jwtManager,
		authMethods: authMethods,
	}
}

type AuthInterceptor struct {
	logger      *log.Logger
	jwtManager  *jwt.Manager
	authMethods map[string]bool
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		i.logger.Info("--> unary interceptor: ", info.FullMethod)
		claims, err := i.Authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if claims != nil {
			ctx = context.WithValue(ctx, ContextKeyID, claims.UserId)
		}
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) Authorize(ctx context.Context, method string) (*jwt.UserClaims, error) {
	b, ok := i.authMethods[method]
	if !ok || !b {
		return nil, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	values := md[headerAuthorize]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	splits := strings.SplitN(values[0], " ", 2)

	if len(splits) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}

	if !strings.EqualFold(splits[0], expectedScheme) {
		return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}

	token := splits[1]
	claims, err := i.jwtManager.Validate(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is invalid")
	}

	return claims, nil
}
