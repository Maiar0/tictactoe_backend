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
	mux.HandleFunc("/test/together", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test/together" {
			http.ServeFile(w, r, "test/together.html")
		}
	})

	loggedMux := utils.LoggingMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", loggedMux))

}
