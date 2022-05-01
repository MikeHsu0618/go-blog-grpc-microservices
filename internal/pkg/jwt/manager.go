package jwt

import (
	"time"

	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewManager(logger *log.Logger, conf *config.Config) *Manager {
	return &Manager{
		secret:  conf.JWT.Secret,
		expires: conf.JWT.Expires,
		logger:  logger,
	}
}

type Manager struct {
	secret  string
	expires time.Duration
	logger  *log.Logger
}

type UserClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.StandardClaims
}

func (manager *Manager) Generate(id uint64) (string, error) {
	claims := UserClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.expires).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secret))
}

func (manager *Manager) Validate(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}

			return []byte(manager.secret), nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims, nil
}
