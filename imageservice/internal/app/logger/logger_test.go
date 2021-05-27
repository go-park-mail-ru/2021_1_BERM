package logger_test

import (
	"github.com/stretchr/testify/require"
	"imageservice/internal/app/logger"
	"net/url"
	"testing"
)

func TestInitLogger(t *testing.T) {
	err := logger.InitLogger("stdout")
	require.NoError(t, err)
}

func TestInitLoggerErr(t *testing.T) {
	err := logger.InitLogger("kek")
	require.Error(t, err)
}

func TestLoggingRequest(t *testing.T) {
	url2 := &url.URL{}
	logger.LoggingRequest(228, url2, "POST")
}