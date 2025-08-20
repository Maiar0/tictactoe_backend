package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Track connections by Player UUID
var playerConnections = make(map[string]*websocket.Conn)

// Track players in each game
var gameGroups = make(map[string][]string)

// addClient adds a client to the playerConnections map
func addClient(conn *websocket.Conn, playerUUID string) {

	playerConnections[playerUUID] = conn
	log.Printf("Player %s connected/reconnected", playerUUID)
}

// removeClient removes a client from the playerConnections map
func removeClient(playerUUID string) {
	delete(playerConnections, playerUUID)
	for gameID, players := range gameGroups {
		for i, player := range players {
			if player == playerUUID {
				gameGroups[gameID] = append(players[:i], players[i+1:]...)
				break
			}
		}
	}
}

// addPlayerToGame adds a player to a game
func addPlayerToGame(playerUUID, gameID string) {
	if gameGroups[gameID] == nil {
		gameGroups[gameID] = []string{}
	}
	gameGroups[gameID] = append(gameGroups[gameID], playerUUID)
	log.Printf("Player %s added to game %s", playerUUID, gameID)
}

func SendToPlayer(playerUUID string, message any) {
	log.Printf("[SendToPlayer] Sending to player %s: %+v", playerUUID, message)
	if conn, exists := playerConnections[playerUUID]; exists {
		jsonMsg, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, jsonMsg)
	}
}

func SendToGame(gameID string, message any) {
	log.Printf("[SendToGame] Sending to game: %+v", message)
	if players, exists := gameGroups[gameID]; exists {
		jsonMsg, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			return
		}
		for _, playerUUID := range players {
			if conn, exists := playerConnections[playerUUID]; exists {
				conn.WriteMessage(websocket.TextMessage, jsonMsg)
			}
		}
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("[HandleWebSocket] Request received: ", r.Method, r.URL.Path)
	// Upgrade HTTP to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	log.Printf("Remote address: %s", ws.RemoteAddr()) // Client IP:port
	log.Printf("Subprotocol: %s", ws.Subprotocol())   // If specified
	//handle heartbeat death
	defer func() {
		// Find which player this connection belongs to
		for playerUUID, conn := range playerConnections { //TODO:: less effecient than creatign another map
			if conn == ws {
				log.Printf("Connection died, cleaning up player: %s", playerUUID)
				removeClient(playerUUID)
				break
			}
		}
		ws.Close()

	}()

	// Handle WebSocket connection
	for {
		// Read message
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Check if it's a heartbeat message and ignores it
		if strings.Contains(string(message), "heartbeat") {
			continue
		}

		log.Printf("Received: %s", string(message))

		// Send response
		handleWebSocketMessage(ws, string(message))
	}
}

type WebSocketMessage struct {
	PlayerUUID string `json:"playerId"`
	GameID     string `json:"gameId"`
	Message    string `json:"message"`
}

func handleWebSocketMessage(ws *websocket.Conn, message string) {
	var msg WebSocketMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("Failed to unmarshal WebSocket message: %v", err)
		ws.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
		return
	}
	log.Printf("[WebSocket]Received: %+v", msg)
	switch msg.Message {
	case "heartbeat":
		SendToPlayer(msg.PlayerUUID, "heartbeat")
	case "register":
		addClient(ws, msg.PlayerUUID)
		addPlayerToGame(msg.PlayerUUID, msg.GameID)
		SendToPlayer(msg.PlayerUUID, "registered")
	case "join_game":
		addPlayerToGame(msg.PlayerUUID, msg.GameID)
		SendToPlayer(msg.PlayerUUID, "joined_game")
	case "leave_game":
		removeClient(msg.PlayerUUID)
		SendToPlayer(msg.PlayerUUID, "left_game")
	case "get_game_state":
		SendToPlayer(msg.PlayerUUID, "game_state")
	default:
		log.Printf("Unknown message: %s", msg.Message)
		SendToPlayer(msg.PlayerUUID, "unknown_message")
	}
}
