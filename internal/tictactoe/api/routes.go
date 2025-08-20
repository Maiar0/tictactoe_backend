package api

import (
	"log"
	"net/http"
)

func Register(mux *http.ServeMux) {
	log.Printf("[Register] tictactoe api endpoints")
	mux.HandleFunc("/api/v1/tictactoe/create", newGame)             // POST
	mux.HandleFunc("/api/v1/tictactoe/state", getGameState)         // POST
	mux.HandleFunc("/api/v1/tictactoe/move", makeMove)              // POST
	mux.HandleFunc("/api/v1/tictactoe/choose_player", choosePlayer) // POST
	mux.HandleFunc("/ws", HandleWebSocket)
}
