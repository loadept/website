package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/loadept/loadept.com/internal/config"
	"github.com/stretchr/testify/assert"
)

const yamlContent = `config:
  debug_mode: true
  app:
    name: example
    static_files: web/static
  http:
    addr: :8080
    read_timeout_seconds: 5
    write_timeout_seconds: 10
    idle_timeout_seconds: 120
  https:
    addr: :4433
    cert_file: cert.pem
    key_file: key.pem
    read_timeout_seconds: 5
    write_timeout_seconds: 10
    idle_timeout_seconds: 120
  database:
    migrations_path: ./migrations
    db_path: db.sqlite3
    pool_size: 3
    busy_timeout: 5000
`

func TestLoadConfig_File(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "config.yaml")
	fatalIfErr(t, os.WriteFile(file, []byte(yamlContent), 0o600))

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"exist file", file, false},
		{"non exist file", filepath.Join(tmpDir, "nonexist.yaml"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.file)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, os.ErrNotExist)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &config.Config{}, cfg)
				assert.Equal(t, cfg.App.Name, "example")
			}
		})
	}
}

func TestLoadConfig_Extensions(t *testing.T) {
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	ymlFile := filepath.Join(tmpDir, "config.yml")
	fatalIfErr(t, os.WriteFile(yamlFile, []byte(yamlContent), 0o600))
	fatalIfErr(t, os.WriteFile(ymlFile, []byte(yamlContent), 0o600))

	txtFile := filepath.Join(tmpDir, "config.txt")
	fatalIfErr(t, os.WriteFile(txtFile, []byte("txt content"), 0o600))

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"valid yaml", yamlFile, false},
		{"valid yml", ymlFile, false},
		{"invalid txt", txtFile, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.file)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, config.ErrFileNoYaml)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &config.Config{}, cfg)
				assert.Equal(t, cfg.App.Name, "example")
			}
		})
	}
}

func fatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
