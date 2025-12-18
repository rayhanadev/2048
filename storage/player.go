package storage

import (
	"database/sql"
	"errors"
	"time"
)

// Player represents a registered player
type Player struct {
	ID                int64
	PubkeyFingerprint string
	Username          string
	CreatedAt         time.Time
}

// ErrPlayerNotFound is returned when a player doesn't exist
var ErrPlayerNotFound = errors.New("player not found")

// GetPlayerByFingerprint retrieves a player by their SSH public key fingerprint
func (db *DB) GetPlayerByFingerprint(fingerprint string) (*Player, error) {
	player := &Player{}
	err := db.conn.QueryRow(`
		SELECT id, pubkey_fingerprint, username, created_at
		FROM players
		WHERE pubkey_fingerprint = ?
	`, fingerprint).Scan(&player.ID, &player.PubkeyFingerprint, &player.Username, &player.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}

	return player, nil
}

// CreatePlayer creates a new player record
func (db *DB) CreatePlayer(fingerprint, username string) (*Player, error) {
	result, err := db.conn.Exec(`
		INSERT INTO players (pubkey_fingerprint, username)
		VALUES (?, ?)
	`, fingerprint, username)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Player{
		ID:                id,
		PubkeyFingerprint: fingerprint,
		Username:          username,
		CreatedAt:         time.Now(),
	}, nil
}

// UpdateUsername updates a player's username
func (db *DB) UpdateUsername(playerID int64, username string) error {
	_, err := db.conn.Exec(`
		UPDATE players SET username = ? WHERE id = ?
	`, username, playerID)
	return err
}

// GetPlayerBestScore returns the highest score for a player
func (db *DB) GetPlayerBestScore(playerID int64) (int, error) {
	var score sql.NullInt64
	err := db.conn.QueryRow(`
		SELECT MAX(score) FROM scores WHERE player_id = ?
	`, playerID).Scan(&score)

	if err != nil {
		return 0, err
	}

	if !score.Valid {
		return 0, nil
	}

	return int(score.Int64), nil
}
