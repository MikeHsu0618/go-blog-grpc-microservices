package jwt

import (
	"testing"

	"blog-grpc-microservices/internal/pkg/config"
	"blog-grpc-microservices/internal/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_Generate(t *testing.T) {
	t.Parallel()
	logger := log.New()
	conf := &config.Config{
		JWT: config.JWT{
			Secret:  "secret",
			Expires: 3600,
		},
	}
	jwtManager := NewManager(logger, conf)
	id := uint64(1)
	tokenStr, err := jwtManager.Generate(id)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)
}

func TestJWTManager_Verify(t *testing.T) {
	t.Parallel()
	logger := log.New()
	conf := &config.Config{
		JWT: config.JWT{
			Secret:  "secret",
			Expires: 3600,
		},
	}
	jwtManager := NewManager(logger, conf)
	id := uint64(2)
	tokenStr, err := jwtManager.Generate(id)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)
	claims, err := jwtManager.Validate(tokenStr)
	require.NoError(t, err)
	require.Equal(t, id, claims.UserId)
}
