package errortools_test

import (
	sqlErr "authorizationservice/pkg/error/errortools"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSqlErrorChoice(t *testing.T) {
	err := sqlErr.SqlErrorChoice(errors.New("mem"))
	require.Error(t, err)
}

func TestErrorHandle(t *testing.T) {
	a, b, c := sqlErr.SqlErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}(nil))
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
