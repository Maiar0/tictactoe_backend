package store

import( 
	"database/sql"
	"time"
)

func ensureSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS game(
		id TEXT PRIMARY KEY,
		state TEXT NOT NULL,
		player_one TEXT,
		player_two TEXT,
		last_update INTEGER NOT NULL,
		status TEXT
	);`)
	return err
}