package config_test

import (
	"os"
	"testing"

	"github.com/asankov/gira/cmd/front-end/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	setenv(t, "GIRA_PORT", "4000")
	setenv(t, "GIRA_LOG_LEVEL", "info")
	setenv(t, "GIRA_API_ADDRESS", "localhost:4000")
	setenv(t, "GIRA_SESSION_SECRET", "sec")
	setenv(t, "GIRA_ENFORCE_HTTPS", "false")

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, config.Port, 4000)
	require.Equal(t, config.LogLevel, "info")
	require.Equal(t, config.APIAddress, "localhost:4000")
	require.Equal(t, config.SessionSecret, "sec")
	require.Equal(t, config.EnforceHTTPS, false)
}

func TestRequiredValues(t *testing.T) {
	config, err := config.NewFromEnv()

	require.Error(t, err)
	require.Nil(t, config)
}

func setenv(t *testing.T, key, value string) {
	t.Helper()
	t.Cleanup(func() {
		os.Unsetenv(key)
	})

	err := os.Setenv(key, value)
	require.NoError(t, err)
}
