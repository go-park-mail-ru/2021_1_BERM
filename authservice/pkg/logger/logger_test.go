package logger_test

import (
	log "authorizationservice/pkg/logger"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestInitLogger(t *testing.T) {
	err := log.InitLogger("stdout")
	require.NoError(t, err)
}

func TestInitLoggerErr(t *testing.T) {
	err := log.InitLogger("kek")
	require.Error(t, err)
}

func TestLoggingRequest(t *testing.T) {
	log.LoggingRequest(uint64(1234), &url.URL{}, "mem")
}
