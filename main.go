package main

import (
	"log"
	"net/http"

	tttApi "github.com/Maiar0/tictactoe_backend/internal/tictactoe/api"
	utils "github.com/Maiar0/tictactoe_backend/internal/utils"
)

func main() {
	log.Println("[Main] Starting TicTacToe backend test...")
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	tttApi.Register(mux)

	// Serve static files test cases
	mux.HandleFunc("/test/create_game", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test/create_game" {
			http.ServeFile(w, r, "test/create_game.html")
		}
	})

	mux.HandleFunc("/test/get_game_state", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test/get_game_state" {
			http.ServeFile(w, r, "test/get_game_state.html")
		}
	})

	loggedMux := utils.LoggingMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", loggedMux))

}
