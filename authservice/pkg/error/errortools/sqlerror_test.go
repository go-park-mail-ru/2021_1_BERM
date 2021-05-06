package errortools

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSqlErrorChoice(t *testing.T) {
	err := SqlErrorChoice(errors.New("mem"))
	require.Error(t, err)
}

func TestErrorHandle(t *testing.T) {
	a, b, c := sqlErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}(map[string]interface{}(nil)))
	require.Equal(t, b, 0)
	require.Equal(t, c, false)
}
