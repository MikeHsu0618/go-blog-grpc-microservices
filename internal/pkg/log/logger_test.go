package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	logger := New()
	require.NotNil(t, logger)
}
