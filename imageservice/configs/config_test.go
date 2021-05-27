package configs

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig(t *testing.T) {
	expectConfig := &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
	config := NewConfig()
	require.Equal(t, config, expectConfig)
}
