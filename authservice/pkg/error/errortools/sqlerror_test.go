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
