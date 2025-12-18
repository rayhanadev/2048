package storage

import (
	"time"
)

// Score represents a game score record
type Score struct {
	ID        int64
	PlayerID  int64
	Score     int
	MaxTile   int
	CreatedAt time.Time
}

// SaveScore saves a game score to the database
func (db *DB) SaveScore(playerID int64, score, maxTile int) error {
	_, err := db.conn.Exec(`
		INSERT INTO scores (player_id, score, max_tile)
		VALUES (?, ?, ?)
	`, playerID, score, maxTile)
	return err
}

// GetPlayerScores returns all scores for a player, ordered by score descending
func (db *DB) GetPlayerScores(playerID int64, limit int) ([]Score, error) {
	rows, err := db.conn.Query(`
		SELECT id, player_id, score, max_tile, created_at
		FROM scores
		WHERE player_id = ?
		ORDER BY score DESC
		LIMIT ?
	`, playerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var s Score
		if err := rows.Scan(&s.ID, &s.PlayerID, &s.Score, &s.MaxTile, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	return scores, rows.Err()
}
