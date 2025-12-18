package storage

import (
	"time"
)

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Rank      int
	Username  string
	Score     int
	MaxTile   int
	CreatedAt time.Time
}

// GetLeaderboard returns the top scores globally
func (db *DB) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	rows, err := db.conn.Query(`
		SELECT 
			p.username,
			s.score,
			s.max_tile,
			s.created_at
		FROM scores s
		JOIN players p ON s.player_id = p.id
		ORDER BY s.score DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.Username, &e.Score, &e.MaxTile, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.Rank = rank
		rank++
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// GetPlayerRank returns the rank of a player's best score
func (db *DB) GetPlayerRank(playerID int64) (int, error) {
	var rank int
	err := db.conn.QueryRow(`
		SELECT COUNT(*) + 1
		FROM scores
		WHERE score > (
			SELECT COALESCE(MAX(score), 0)
			FROM scores
			WHERE player_id = ?
		)
	`, playerID).Scan(&rank)

	return rank, err
}
