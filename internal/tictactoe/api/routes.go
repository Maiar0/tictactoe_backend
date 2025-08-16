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
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not Allowed.")
		return
	}
	//read request body
	var req newGameReq
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	//logic
	if req.PlayerUUID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Player UUID Required.")
		return
	}
	log.Println("[newGame] Creating new game for player UUID: ", req.PlayerUUID)
	id, err := tttStore.NewGame()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to create game.")
		return
	}
	//write response
	utils.WriteJSONResponse(w, http.StatusCreated, newGameResp{GameID: id})
	log.Println("[newGame] Game created successfully with ID: ", id)
}

type getGameStateReq struct {
	PlayerUUID string `json:"player_uuid"`
	GameID     string `json:"game_id"`
}
type getGameStateResp struct {
	GameState string `json:"game_state"`
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	log.Println("[getGameState] Request received: ", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not Allowed.")
		return
	}
	var req getGameStateReq
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if req.PlayerUUID == "" || req.GameID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Player UUID and Game ID Required.")
		return
	}
	log.Println("[getGameState] Getting game state for player UUID: ", req.PlayerUUID, " and game ID: ", req.GameID)
	gameState, err := tttStore.GetGameState(req.GameID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get game state.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getGameStateResp{GameState: gameState}); err != nil {
		log.Printf("[getGameState] encode error: %v", err)
	}
	log.Println("[getGameState] Game state retrieved successfully: ", gameState)

}
func Register(mux *http.ServeMux) {
	log.Printf("[Register] tictactoe api endpoints")
	mux.HandleFunc("/api/v1/game/create", newGame)     // POST
	mux.HandleFunc("/api/v1/game/state", getGameState) // GET

}
