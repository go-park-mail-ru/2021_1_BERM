package errortools

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidationErrorHandle(t *testing.T) {
	a, b, c := validationErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}(map[string]interface{}(nil)))
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
