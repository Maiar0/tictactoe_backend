package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type GameStore struct {
	db *sql.DB
}

func NewGameStore(db *sql.DB) *GameStore {
	return &GameStore{db: db}
}

// CreateGameState inserts a new row representing a game state
func (g *GameStore) CreateGameState(state, playerOne, playerTwo, status string) (sql.Result, error) {
	result, err := g.db.Exec(`
        INSERT INTO game (state, player_x, player_o, last_update, status)
        VALUES (?, ?, ?, ?, ?)
    `, state, playerOne, playerTwo, time.Now().Unix(), status)
	if err != nil {
		log.Println("[CreateGameState] Failed to create game state: ", err)
		return nil, err
	}
	log.Println("[CreateGameState] Game state created: ", result)
	return result, err
}

// ReadGameState queries rows by a field and value or all if no value
func (g *GameStore) ReadGameState(field string, values ...any) (*sql.Rows, error) {
	if len(values) == 0 {
		// No value provided → read all rows
		return g.db.Query("SELECT * FROM game")
	}
	// Value provided → filter by field
	return g.db.Query(fmt.Sprintf("SELECT * FROM game WHERE %s = ?", field), values[0])
}

// UpdateGameState updates columns in rows that match field=value
func (g *GameStore) UpdateGameState(field string, value any, updates map[string]any) (sql.Result, error) {
	setClause := ""
	args := []any{}
	for k, v := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = ?", k)
		args = append(args, v)
	}
	args = append(args, value)
	query := fmt.Sprintf("UPDATE game SET %s WHERE %s = ?", setClause, field)
	return g.db.Exec(query, args...)
}

// DeleteGameState deletes rows matching field=value
func (g *GameStore) DeleteGameState(field string, value any) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM game WHERE %s = ?", field)
	return g.db.Exec(query, value)
}
