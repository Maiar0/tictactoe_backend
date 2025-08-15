package store

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"

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

func CreateGame() (string, error) {
	log.Println("[CreateGame] Starting new game creation: ", baseDir, " : ", schemaPath, " : ", initialState)
	id := newGameID()
	log.Println("[CreateGame] Generated game ID", id)
	st := sqlite.New(baseDir)
	db, err := st.OpenFor(id, schemaPath)
	if err != nil {
		log.Println("[CreateGame] Failed to open DB: ", err)
		return "", err
	}
	defer db.Close()
	log.Println("[CreateGame] DB Opened succesfully: ", id)
	res, err := db.Exec(`INSERT INTO game(state,player_one,player_two,last_update,status)
	                  VALUES (?,?,?,?,?)`, initialState, "", "", time.Now().Unix(), "active")
	if err != nil {
		return "", err
	}
	insertID, _ := res.LastInsertId() // return PK
	log.Println("[CreateGame] Game inserted succesfully. Insert ID: ", insertID)
	return id, nil
}
