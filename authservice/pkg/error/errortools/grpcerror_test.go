package errortools_test

import (
	errTools "authorizationservice/pkg/error/errortools"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrpcErrorChoice(t *testing.T) {
	a, b, c := errTools.GrpcErrorHandle(errors.New("Err"))
	require.Equal(t, a, nil)
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
