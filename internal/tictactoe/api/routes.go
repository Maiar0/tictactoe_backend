package api

import (
	"encoding/json"
	"log"
	"net/http"

	tttStore "github.com/Maiar0/tictactoe_backend/internal/tictactoe/store"
	utils "github.com/Maiar0/tictactoe_backend/internal/utils"
)

type newGameReq struct {
	PlayerUUID string `json:"player_uuid"`
}
type newGameResp struct {
	GameID string `json:"game_id"`
}

func newGame(w http.ResponseWriter, r *http.Request) {
	log.Println("[newGame] Request received: ", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not Allowed.")
		return
	}
	http.MaxBytesReader(w, r.Body, 1<<20)
	var in newGameReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if in.PlayerUUID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Player UUID Required.")
		return
	}
	log.Println("[newGame] Creating new game for player UUID: ", in.PlayerUUID)
	id, err := tttStore.NewGame()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create game.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newGameResp{GameID: id}); err != nil {
		log.Printf("[newGame] encode error: %v", err)
	}
	log.Println("[newGame] Game created successfully with ID: ", id)
}

func Register(mux *http.ServeMux) {
	log.Printf("[Register] tictactoe api endpoints")
	mux.HandleFunc("/api/v1/game/create", newGame) // POST

}
