package configs_test

import (
	"github.com/stretchr/testify/require"
	conf "imageservice/configs"
	"testing"
)

func TestNewConfig(t *testing.T) {
	expectConfig := &conf.Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
	config := conf.NewConfig()
	require.Equal(t, config, expectConfig)
}
