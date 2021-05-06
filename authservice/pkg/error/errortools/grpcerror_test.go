package errortools

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrpcErrorChoice(t *testing.T) {
	a, b, c := grpcErrorHandle(errors.New("Err"))
	require.Equal(t, a, nil)
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
