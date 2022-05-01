package dbcontext

import (
	"testing"

	"blog-grpc-microservices/internal/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	path := config.GetPath()
	conf, err := config.Load(path)
	require.NoError(t, err)

	userDB, err := NewUserDB(conf)
	require.NoError(t, err)
	require.NotNil(t, userDB)

	postDB, err := NewPostDB(conf)
	require.NoError(t, err)
	require.NotNil(t, postDB)

	commentDB, err := NewCommentDB(conf)
	require.NoError(t, err)
	require.NotNil(t, commentDB)

}
