package store

import (
	"crypto/rand"
	"log"
	"math/big"

	sqlite "github.com/Maiar0/tictactoe_backend/internal/store"
)

const (
	baseDir      = "Storage/games/tictactoe"
	schemaPath   = "internal/tictactoe/store/schema.sql"
	initialState = ".........X" //9 blank squares last value is whose turn it is WIP
)

type GameState struct {
	ID         int16  `db:"id"`
	State      string `db:"state"`
	PlayerX    string `db:"player_one"`
	PlayerO    string `db:"player_two"`
	LastUpdate int64  `db:"last_update"`
	Status     string `db:"status"`
}

func newGameID() string { //TODO:: huh
	const bank = "abcdefghijklmnopqrstuvwxyz0123456789"
	const n = 9
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(len(bank)))) //TODO:: Error handling
		b[i] = bank[k.Int64()]
	}
	return string(b)
}

func NewGame() (string, error) {
	log.Println("[NewGame] Starting new game creation: ", baseDir, " : ", schemaPath, " : ", initialState)
	id := newGameID()
	log.Println("[NewGame] Generated game ID", id)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(id, schemaPath)
	if err != nil {
		log.Println("[NewGame] Failed to open DB: ", err)
		return "", err
	}

	log.Println("[NewGame] DB Opened succesfully: ", id)
	gameStore := NewGameStore(db)
	res, err := gameStore.CreateGameState(initialState, "", "", "active")
	defer db.Close()
	if err != nil {
		log.Println("[NewGame] Failed to create game state: ", err)
		return "", err
	}
	insertID, _ := res.LastInsertId() // return PK
	log.Println("[NewGame] Game inserted succesfully. Insert ID: ", insertID)
	return id, nil
}

func GetGameState(gameID string) (GameState, error) {
	var gameState GameState
	log.Println("[GetGameState] Getting game state for game ID: ", gameID)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(gameID, schemaPath)
	if err != nil {
		log.Println("[GetGameState] Failed to open DB: ", err)
		return gameState, err
	}
	log.Println("[GetGameState] DB Opened succesfully: ", gameID)
	gameStore := NewGameStore(db)
	rows, err := gameStore.ReadGameState("id")
	if err != nil {
		log.Println("[GetGameState] Failed to read game state: ", err)
		return gameState, err
	}
	var highestID int16
	for rows.Next() {
		var id int16
		var state GameState
		if err := rows.Scan(&id, &state.State, &state.PlayerX, &state.PlayerO, &state.LastUpdate, &state.Status); err != nil {
			log.Println("[GetGameState] Failed to scan row: ", err)
			return gameState, err
		}
		if id > highestID {
			highestID = id
			gameState = state
		}
	}
	log.Println("[GetGameState] ID: ", highestID, " State: ", gameState.State)
	defer db.Close()
	return gameState, nil
}

func UpdateGameState(gameID string, gameState GameState) error {
	log.Println("[UpdateGameState] Updating game state for game ID: ", gameID)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(gameID, schemaPath)
	if err != nil {
		log.Println("[UpdateGameState] Failed to open DB: ", err)
		return err
	}
	log.Println("[UpdateGameState] DB Opened succesfully: ", gameID)
	gameStore := NewGameStore(db)
	_, err = gameStore.CreateGameState(gameState.State, gameState.PlayerX, gameState.PlayerO, gameState.Status)
	if err != nil {
		log.Println("[UpdateGameState] Failed to update game state: ", err)
		return err
	}
	log.Println("[UpdateGameState] Game state updated succesfully: ", gameState)
	return nil
}
