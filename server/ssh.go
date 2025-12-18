package server

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	gossh "golang.org/x/crypto/ssh"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rayhanadev/2048/config"
	"github.com/rayhanadev/2048/storage"
	"github.com/rayhanadev/2048/ui"
)

// Server represents the SSH server
type Server struct {
	config *config.Config
	db     *storage.DB
	server *ssh.Server
}

// NewServer creates a new SSH server
func NewServer(cfg *config.Config, db *storage.DB) (*Server, error) {
	s := &Server{
		config: cfg,
		db:     db,
	}

	// Ensure host key exists
	if err := s.ensureHostKey(); err != nil {
		return nil, fmt.Errorf("failed to ensure host key: %w", err)
	}

	// Create Wish server
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.SSHHost, fmt.Sprintf("%d", cfg.SSHPort))),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithPublicKeyAuth(s.publicKeyHandler),
		wish.WithMiddleware(
			bubbletea.Middleware(s.teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH server: %w", err)
	}

	s.server = server
	return s, nil
}

// publicKeyHandler handles public key authentication
// We accept all public keys (public access) but extract the fingerprint for player identification
func (s *Server) publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	// Always accept - public access
	// The key is stored in context for later use
	return true
}

// teaHandler creates a Bubbletea program for each SSH session
func (s *Server) teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Get terminal size
	pty, _, ok := sess.Pty()
	if !ok {
		log.Warn("No PTY requested, using default size")
	}

	// Extract public key fingerprint
	fingerprint := s.getFingerprint(sess)

	// Check if player exists
	var player *storage.Player
	if s.db != nil {
		p, err := s.db.GetPlayerByFingerprint(fingerprint)
		if err == nil {
			player = p
		}
	}

	// Create the model
	var initialState ui.AppState
	if player == nil {
		initialState = ui.StateUsernameEntry
	} else {
		initialState = ui.StatePlaying
	}
	model := ui.NewModel(s.db, fingerprint, player, initialState)

	// Set initial terminal size
	if ok {
		model = model.WithSize(pty.Window.Width, pty.Window.Height)
	}

	return model, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

// getFingerprint extracts the SSH public key fingerprint from the session
func (s *Server) getFingerprint(sess ssh.Session) string {
	key := sess.PublicKey()
	if key == nil {
		// Fallback for connections without public key (shouldn't happen with pubkey auth)
		return fmt.Sprintf("anon-%s-%d", sess.RemoteAddr().String(), time.Now().UnixNano())
	}

	// Generate SHA256 fingerprint
	hash := sha256.Sum256(key.Marshal())
	fingerprint := base64.RawStdEncoding.EncodeToString(hash[:])
	return fmt.Sprintf("SHA256:%s", fingerprint)
}

// ensureHostKey ensures an SSH host key exists, generating one if needed
func (s *Server) ensureHostKey() error {
	keyPath := s.config.HostKeyPath

	// Create directory if needed
	keyDir := filepath.Dir(keyPath)
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Check if key exists
	if _, err := os.Stat(keyPath); err == nil {
		log.Info("Using existing host key", "path", keyPath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check host key: %w", err)
	}

	// Generate new ED25519 key
	log.Info("Generating new host key", "path", keyPath)

	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Marshal the private key in OpenSSH format
	privKeyBlock, err := gossh.MarshalPrivateKey(privKey, "")
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	privKeyPEM := pem.EncodeToMemory(privKeyBlock)

	// Write to file with restrictive permissions
	if err := os.WriteFile(keyPath, privKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write host key: %w", err)
	}

	log.Info("Generated new host key", "path", keyPath)
	return nil
}

// Start starts the SSH server
func (s *Server) Start() error {
	addr := net.JoinHostPort(s.config.SSHHost, fmt.Sprintf("%d", s.config.SSHPort))
	log.Info("Starting SSH server", "address", addr)

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Server error", "error", err)
		}
	}()

	log.Info("SSH 2048 server is running", "address", addr)
	log.Info("Connect with: ssh localhost -p " + fmt.Sprintf("%d", s.config.SSHPort))

	<-done

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
