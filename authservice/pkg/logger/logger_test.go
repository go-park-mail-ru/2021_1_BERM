package logger

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger("stdout")
	require.NoError(t, err)
}

func TestInitLoggerErr(t *testing.T) {
	err := InitLogger("kek")
	require.Error(t, err)
}

func TestLoggingRequest(t *testing.T) {
	LoggingRequest(uint64(1234), &url.URL{}, "mem")
}
