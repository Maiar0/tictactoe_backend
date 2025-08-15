CREATE TABLE IF NOT EXISTS game(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		state TEXT NOT NULL,
		player_one TEXT,
		player_two TEXT,
		last_update INTEGER NOT NULL,
		status TEXT
	)