package logger

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger("stdout")
	require.NoError(t, err)
}
