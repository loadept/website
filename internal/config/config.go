// Package config It allows you to create application configurations from a YAML file.
// NOTE: This implementation is in case my application grows and this type of configuration
// is needed. For now, what I tried to do was just unnecessary over-engineering :)
package config

import (
	"errors"
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr      string
	AppURL    string
	DBPath    string
	GithubURL string
}

func Load() (*Config, error) {
	_ = godotenv.Overload()
	cfg := new(Config)

	addr := Env("ADDR", "127.0.0.1:8080")
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return nil, errors.New("env var ADDR is invalid")
	}
	cfg.Addr = addr

	appURL := Env("APP_URL", "http://localhost:8080")
	au, err := url.Parse(appURL)
	if err != nil || au.Scheme == "" || au.Host == "" {
		return nil, errors.New("env var APP_URL is invalid")
	}
	cfg.AppURL = appURL

	dbPath := Env("DB_PATH")
	if dbPath == "" {
		return nil, errors.New("env var DB_PATH is required")
	}
	cfg.DBPath = dbPath

	githubURL := Env("GITHUB_URL")
	if githubURL == "" {
		return nil, errors.New("env var GITHUB_URL is required")
	}
	gu, err := url.Parse(githubURL)
	if err != nil || gu.Scheme == "" || gu.Host == "" {
		return nil, errors.New("env var GITHUB_URL is invalid")
	}
	cfg.GithubURL = githubURL

	return cfg, nil
}

func Env(key string, defaultValue ...string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	if len(defaultValue) > 0 {
		return defaultValue[len(defaultValue)-1]
	}
	return ""
}
