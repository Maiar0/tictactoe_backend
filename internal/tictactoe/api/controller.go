package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	tttService "github.com/Maiar0/tictactoe_backend/internal/tictactoe/service"
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
	if r.Method != http.MethodPost {
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
	utils.WriteJSONResponse(w, http.StatusOK, getGameStateResp{GameState: gameState.State})
	log.Println("[getGameState] Game state retrieved successfully: ", gameState)

}

type choosePlayerReq struct {
	PlayerUUID   string `json:"player_uuid"`
	GameID       string `json:"game_id"`
	PlayerChoice string `json:"choice"` // "x" or "o"
}
type choosePlayerResp struct {
	GameState string `json:"game_state"`
}

func choosePlayer(w http.ResponseWriter, r *http.Request) {
	log.Println("[choosePlayer] Request received: ", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not Allowed.")
		return
	}
	var req choosePlayerReq
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if req.PlayerUUID == "" || req.GameID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Player UUID and Game ID Required.")
		return
	}
	log.Println("[choosePlayer] Choosing player for game ID: ", req.GameID)
	gameState, err := tttStore.GetGameState(req.GameID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get game state.")
		return
	}
	//correct input
	if req.PlayerChoice != "x" && req.PlayerChoice != "o" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Player Choice must be 'x' or 'o'.")
		return
	}
	//logic to choose player
	if req.PlayerChoice == "x" && gameState.PlayerX == "" {
		gameState.PlayerX = req.PlayerUUID
	} else if gameState.PlayerO == "" {
		gameState.PlayerO = req.PlayerUUID
	} else {
		utils.WriteJSONError(w, http.StatusForbidden, "Players already chosen. Game is in progress.")
		return
	}
	//update game state
	tttStore.UpdateGameState(req.GameID, gameState)
	//get game state truth
	gameState, err = tttStore.GetGameState(req.GameID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get game state.")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, choosePlayerResp{GameState: gameState.State})
	log.Println("[choosePlayer] Player chosen successfully: ", gameState)
}

type makeMoveReq struct {
	PlayerUUID string `json:"player_uuid"`
	GameID     string `json:"game_id"`
	Move       string `json:"move"` // 2 Character string representing the move char o || x and a number 0-8 (e.g. "x0", "o2", "x8")
}

type makeMoveResp struct {
	GameState string `json:"game_state"`
}

func makeMove(w http.ResponseWriter, r *http.Request) {
	log.Println("[makeMove] Request received: ", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not Allowed.")
		return
	}
	var req makeMoveReq
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if req.PlayerUUID == "" || req.GameID == "" || req.Move == "" || len(req.Move) != 2 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Player UUID, Game ID, and Move Required. Move must be 2 characters.")
		return
	}
	gameState, err := tttStore.GetGameState(req.GameID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get game state.")
		return
	}
	if !strings.EqualFold(string(gameState.State[9]), string(req.Move[0])) {
		utils.WriteJSONError(w, http.StatusForbidden, "It's not your turn.")
		return
	}
	//validate turn
	turn := gameState.State[9]
	if turn == '.' {
		utils.WriteJSONError(w, http.StatusForbidden, "Game is not in progress.")
		return
	}
	var playersTurn string
	if turn == 'x' {
		playersTurn = gameState.PlayerX
	} else {
		playersTurn = gameState.PlayerO
	}
	if playersTurn != req.PlayerUUID {
		utils.WriteJSONError(w, http.StatusForbidden, "It's not your turn.")
		return
	}
	//begin move logic
	log.Println("[makeMove] Making move for player UUID: ", req.PlayerUUID, " and game ID: ", req.GameID, " with move: ", req.Move)
	finalGameState, err := tttService.MakeMove(req.GameID, req.Move)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to make move.")
		return
	}
	//send game state to websocket
	gameStateJSON, err := json.Marshal(map[string]string{"game_state": gameState.State})
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to marshal game state.")
		return
	}
	SendToGame(req.GameID, string(gameStateJSON))
	utils.WriteJSONResponse(w, http.StatusOK, makeMoveResp{GameState: finalGameState})
	log.Println("[makeMove] Move made successfully: ", gameState)
}
