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
	log.Println("[CreateGame] Starting new game creation: ", baseDir, " : ", schemaPath, " : ", initialState)
	id := newGameID()
	log.Println("[CreateGame] Generated game ID", id)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(id, schemaPath)
	if err != nil {
		log.Println("[CreateGame] Failed to open DB: ", err)
		return "", err
	}

	log.Println("[CreateGame] DB Opened succesfully: ", id)
	gameStore := NewGameStore(db)
	res, err := gameStore.CreateGameState(initialState, "", "", "active")
	defer db.Close()
	if err != nil {
		return "", err
	}
	insertID, _ := res.LastInsertId() // return PK
	log.Println("[CreateGame] Game inserted succesfully. Insert ID: ", insertID)
	return id, nil
}

func GetGameState(gameID string) (string, error) {
	log.Println("[GetGameState] Getting game state for game ID: ", gameID)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(gameID, schemaPath)
	if err != nil {
		log.Println("[GetGameState] Failed to open DB: ", err)
		return "", err
	}
	log.Println("[GetGameState] DB Opened succesfully: ", gameID)
	gameStore := NewGameStore(db)
	rows, err := gameStore.ReadGameState("id")
	if err != nil {
		log.Println("[GetGameState] Failed to read game state: ", err)
		return "", err
	}
	var highestID int16
	var state string
	for rows.Next() {
		var id int16
		if err := rows.Scan(&id, &state); err != nil {
			log.Println("[GetGameState] Failed to scan row: ", err)
			return "", err
		}
		if id > highestID {
			highestID = id
		}
		log.Println("[GetGameState] ID: ", id, " State: ", state)
	}
	defer db.Close()
	return state, nil
}
