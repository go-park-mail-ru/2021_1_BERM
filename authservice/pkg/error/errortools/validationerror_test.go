package errortools_test

import (
	errTools "authorizationservice/pkg/error/errortools"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidationErrorHandle(t *testing.T) {
	a, b, c := errTools.ValidationErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}(nil))
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
