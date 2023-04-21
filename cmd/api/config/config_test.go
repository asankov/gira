package config_test

import (
	"os"
	"testing"

	"github.com/asankov/gira/cmd/api/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	setenv(t, "GIRA_PORT", "4000")
	setenv(t, "GIRA_SECRET", "sec")
	setenv(t, "GIRA_USE_SSL", "false")
	setenv(t, "GIRA_LOG_LEVEL", "info")
	setenv(t, "GIRA_DB_HOST", "localhost")
	setenv(t, "GIRA_DB_PORT", "5432")
	setenv(t, "GIRA_DB_PASSWORD", "pass")
	setenv(t, "GIRA_DB_USER", "user")
	setenv(t, "GIRA_DB_NAME", "name")

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.NotNil(t, config)
	require.Equal(t, config.Port, 4000)
	require.Equal(t, config.Secret, "sec")
	require.Equal(t, config.UseSSL, false)
	require.Equal(t, config.LogLevel, "info")
	require.Equal(t, config.DB.Host, "localhost")
	require.Equal(t, config.DB.Port, 5432)
	require.Equal(t, config.DB.Password, "pass")
	require.Equal(t, config.DB.User, "user")
	require.Equal(t, config.DB.Name, "name")
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
