package config

import (
	"errors"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type ConfigYaml struct {
	Config Config
}

type Config struct {
	DebugMode bool     `yaml:"debug_mode"`
	App       App      `yaml:"app"`
	HTTP      HTTP     `yaml:"http"`
	Database  Database `yaml:"database"`
}

type App struct {
	Name        string `yaml:"name"`
	StaticFiles string `yaml:"static_files"`
}

type HTTP struct {
	Addr                string `yaml:"addr"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
	IdleTimeoutSeconds  int    `yaml:"idle_timeout_seconds"`
}

type Database struct {
	MigrationsPath string `yaml:"migrations_path"`
	DBPath         string `yaml:"db_path"`
	PoolSize       int    `yaml:"pool_size"`
	BusyTimeout    int    `yaml:"busy_timeout"`
}

var ErrFileNoYaml = errors.New("file formar not allowed")

func Load(name string) (*Config, error) {
	cleanedName := filepath.Clean(name)
	if ext := filepath.Ext(cleanedName); ext != ".yaml" && ext != ".yml" {
		return nil, ErrFileNoYaml
	}
	if !filepath.IsAbs(cleanedName) {
		var err error
		cleanedName, err = filepath.Abs(cleanedName)
		if err != nil {
			return nil, err
		}
	}
	baseDir := filepath.Dir(cleanedName)

	var config ConfigYaml
	file, err := os.ReadFile(cleanedName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if !filepath.IsAbs(config.Config.App.StaticFiles) {
		config.Config.App.StaticFiles = filepath.Join(baseDir, config.Config.App.StaticFiles)
	}
	if !filepath.IsAbs(config.Config.Database.MigrationsPath) {
		config.Config.Database.MigrationsPath = filepath.Join(baseDir, config.Config.Database.MigrationsPath)
	}
	if !filepath.IsAbs(config.Config.Database.DBPath) {
		config.Config.Database.DBPath = filepath.Join(baseDir, config.Config.Database.DBPath)
	}
	return &config.Config, nil
}
