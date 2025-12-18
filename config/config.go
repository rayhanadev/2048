package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	SSHPort     int
	SSHHost     string
	DataDir     string
	HostKeyPath string
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	cfg := &Config{
		SSHPort:     23234,
		SSHHost:     "0.0.0.0",
		DataDir:     "./data",
		HostKeyPath: ".ssh/2048_host_key",
	}

	if port := os.Getenv("SSH_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.SSHPort = p
		}
	}

	if host := os.Getenv("SSH_HOST"); host != "" {
		cfg.SSHHost = host
	}

	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
	}

	if hostKeyPath := os.Getenv("HOST_KEY_PATH"); hostKeyPath != "" {
		cfg.HostKeyPath = hostKeyPath
	}

	return cfg
}

// EnsureDirectories creates necessary directories if they don't exist
func (c *Config) EnsureDirectories() error {
	// Create data directory
	if err := os.MkdirAll(c.DataDir, 0755); err != nil {
		return err
	}

	// Create host key directory
	hostKeyDir := filepath.Dir(c.HostKeyPath)
	if err := os.MkdirAll(hostKeyDir, 0700); err != nil {
		return err
	}

	return nil
}

// DatabasePath returns the full path to the SQLite database
func (c *Config) DatabasePath() string {
	return filepath.Join(c.DataDir, "2048.db")
}
