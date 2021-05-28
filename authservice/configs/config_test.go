пшpackage configs_test

import (
	conf "authorizationservice/configs"
	"github.com/stretchr/testify/require"
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
