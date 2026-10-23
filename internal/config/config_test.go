package config_test

import (
	"os"
	"testing"

	"github.com/loadept/website/internal/config"
	"github.com/stretchr/testify/assert"
)

const testEnvKey = "CONFIG_TEST_KEY"

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		prev, had := os.LookupEnv(k)
		os.Unsetenv(k)
		t.Cleanup(func() {
			if had {
				os.Setenv(k, prev)
			} else {
				os.Unsetenv(k)
			}
		})
	}
}

func clearEnv(t *testing.T) {
	t.Helper()
	unsetEnv(t, "ADDR", "APP_URL", "DB_PATH", "GITHUB_URL")
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func validEnv() map[string]string {
	return map[string]string{
		"ADDR":       "127.0.0.1:8080",
		"APP_URL":    "http://localhost:8080",
		"DB_PATH":    "test.db",
		"GITHUB_URL": "https://github.com/example",
	}
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)
	setEnv(t, "DB_PATH", "test.db")
	setEnv(t, "GITHUB_URL", "https://github.com/example")

	cfg, err := config.Load()
	if !assert.NoError(t, err) || !assert.NotNil(t, cfg) {
		return
	}
	assert.Equal(t, "127.0.0.1:8080", cfg.Addr)
	assert.Equal(t, "http://localhost:8080", cfg.AppURL)
	assert.Equal(t, "test.db", cfg.DBPath)
	assert.Equal(t, "https://github.com/example", cfg.GithubURL)
}

func TestLoad_FromEnv(t *testing.T) {
	clearEnv(t)
	setEnv(t, "ADDR", ":9090")
	setEnv(t, "APP_URL", "https://example.com")
	setEnv(t, "DB_PATH", "/var/data/app.db")
	setEnv(t, "GITHUB_URL", "https://github.com/loadept")

	cfg, err := config.Load()
	if !assert.NoError(t, err) || !assert.NotNil(t, cfg) {
		return
	}
	assert.Equal(t, ":9090", cfg.Addr)
	assert.Equal(t, "https://example.com", cfg.AppURL)
	assert.Equal(t, "/var/data/app.db", cfg.DBPath)
	assert.Equal(t, "https://github.com/loadept", cfg.GithubURL)
}

func TestLoad_Errors(t *testing.T) {
	tests := []struct {
		name    string
		set     map[string]string
		drop    []string
		wantErr string
	}{
		{
			name:    "addr_without_port",
			set:     map[string]string{"ADDR": "localhost"},
			wantErr: "env var ADDR is invalid",
		},
		{
			name:    "app_url_without_scheme",
			set:     map[string]string{"APP_URL": "not-a-url"},
			wantErr: "env var APP_URL is invalid",
		},
		{
			name:    "app_url_without_host",
			set:     map[string]string{"APP_URL": "http://"},
			wantErr: "env var APP_URL is invalid",
		},
		{
			name:    "db_path_missing",
			drop:    []string{"DB_PATH"},
			wantErr: "env var DB_PATH is required",
		},
		{
			name:    "db_path_empty",
			set:     map[string]string{"DB_PATH": ""},
			wantErr: "env var DB_PATH is required",
		},
		{
			name:    "github_url_missing",
			drop:    []string{"GITHUB_URL"},
			wantErr: "env var GITHUB_URL is required",
		},
		{
			name:    "github_url_without_scheme",
			set:     map[string]string{"GITHUB_URL": "foo"},
			wantErr: "env var GITHUB_URL is invalid",
		},
		{
			name:    "github_url_without_host",
			set:     map[string]string{"GITHUB_URL": "http://"},
			wantErr: "env var GITHUB_URL is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			env := validEnv()
			for _, k := range tt.drop {
				delete(env, k)
			}
			for k, v := range tt.set {
				env[k] = v
			}
			for k, v := range env {
				setEnv(t, k, v)
			}

			cfg, err := config.Load()
			assert.EqualError(t, err, tt.wantErr)
			assert.Nil(t, cfg)
		})
	}
}

func TestEnv_ReturnsValueWhenSet(t *testing.T) {
	setEnv(t, testEnvKey, "valor")
	assert.Equal(t, "valor", config.Env(testEnvKey))
}

func TestEnv_ReturnsDefaultWhenUnset(t *testing.T) {
	unsetEnv(t, testEnvKey)
	assert.Equal(t, "fallback", config.Env(testEnvKey, "fallback"))
}

func TestEnv_ReturnsLastDefault(t *testing.T) {
	unsetEnv(t, testEnvKey)
	assert.Equal(t, "segundo", config.Env(testEnvKey, "primero", "segundo"))
}

func TestEnv_EmptyWhenUnsetAndNoDefault(t *testing.T) {
	unsetEnv(t, testEnvKey)
	assert.Equal(t, "", config.Env(testEnvKey))
}

func TestEnv_SetEmptyValueIgnoresDefault(t *testing.T) {
	setEnv(t, testEnvKey, "")
	assert.Equal(t, "", config.Env(testEnvKey, "fallback"))
}
