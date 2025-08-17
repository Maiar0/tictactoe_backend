CREATE TABLE IF NOT EXISTS game(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		state TEXT NOT NULL,
		player_x TEXT,
		player_o TEXT,
		last_update INTEGER NOT NULL,
		status TEXT
	)