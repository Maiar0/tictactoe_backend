package api

import (
	"log"
	"net/http"
)

func Register(mux *http.ServeMux) {
	log.Printf("[Register] tictactoe api endpoints")
	mux.HandleFunc("/api/v1/game/create", newGame)             // POST
	mux.HandleFunc("/api/v1/game/state", getGameState)         // POST
	mux.HandleFunc("/api/v1/game/move", makeMove)              // POST
	mux.HandleFunc("/api/v1/game/choose_player", choosePlayer) // POST
	mux.HandleFunc("/ws", HandleWebSocket)
}
