package main

import (
	"os"

	"github.com/charmbracelet/log"

	"github.com/rayhanadev/2048/config"
	"github.com/rayhanadev/2048/server"
	"github.com/rayhanadev/2048/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up logging
	log.SetLevel(log.InfoLevel)
	log.Info("SSH 2048 Server starting...")
	log.Info("Configuration loaded",
		"port", cfg.SSHPort,
		"host", cfg.SSHHost,
		"data_dir", cfg.DataDir,
	)

	// Initialize database
	db, err := storage.NewDB(cfg.DataDir)
	if err != nil {
		log.Fatal("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("Database initialized")

	// Create and start SSH server
	srv, err := server.NewServer(cfg, db)
	if err != nil {
		log.Fatal("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Start the server (blocks until shutdown)
	if err := srv.Start(); err != nil {
		log.Fatal("Server error", "error", err)
		os.Exit(1)
	}

	log.Info("Server shutdown complete")
}
