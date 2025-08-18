package service

import (
	"errors"
	"log"

	store "github.com/Maiar0/tictactoe_backend/internal/tictactoe/store"
)

func MakeMove(gameID, move string) (string, error) {
	//prepare move
	turn := move[0]
	position := int(move[1]-'0') - 1
	log.Println("[MakeMove] Turn: ", turn, " Position: ", position)
	gameState, err := store.GetGameState(gameID)
	if err != nil {
		log.Println("[MakeMove] Failed to get game state: ", err)
		return "", err
	}
	//validate move
	if gameState.State[position] != '.' {
		log.Println("[MakeMove] Invalid move: ", move)
		return "", errors.New("invalid move")
	}
	//alter game state
	gameState = alterGameState(gameState, byte(turn), position)
	err = store.UpdateGameState(gameID, gameState)
	if err != nil {
		log.Println("[MakeMove] Failed to update game state: ", err)
		return "", err
	}
	//check if game is over
	if gameWon(gameState.State) {
		if gameState.State[9] == 'x' {
			gameState.Status = gameState.PlayerX
		} else {
			gameState.Status = gameState.PlayerO
		}
		err = store.UpdateGameState(gameID, gameState)
		if err != nil {
			log.Println("[MakeMove] Failed to update game state: ", err)
			return "", err
		}
	} else if gameTied(gameState.State) {
		gameState.Status = "tied"
		err = store.UpdateGameState(gameID, gameState)
		if err != nil {
			log.Println("[MakeMove] Failed to update game state: ", err)
			return "", err
		}
	}
	//ensure we are senidng back truth
	gameState, err = store.GetGameState(gameID)
	if err != nil {
		log.Println("[MakeMove] Failed to get game state: ", err)
		return "", err
	}
	return gameState.State, nil
}

func alterGameState(gameState store.GameState, turn byte, position int) store.GameState {
	log.Println("[alterGameState] Altering game state: ", gameState)
	gameBytes := []byte(gameState.State)
	gameBytes[position] = turn
	var nextTurn byte
	if turn == 'x' {
		nextTurn = 'o'
	} else {
		nextTurn = 'x'
	}
	gameBytes[9] = nextTurn
	gameState.State = string(gameBytes)
	gameState.Status = "active"
	log.Println("[alterGameState] Game state altered: ", gameState)
	return gameState
}

var wins = [8][3]int{
	{0, 1, 2}, // row 1
	{3, 4, 5}, // row 2
	{6, 7, 8}, // row 3
	{0, 3, 6}, // col 1
	{1, 4, 7}, // col 2
	{2, 5, 8}, // col 3
	{0, 4, 8}, // diagonal
	{2, 4, 6}, // diagonal
}

func gameWon(gameState string) bool {
	for _, win := range wins {
		if gameState[win[0]] != '.' && gameState[win[0]] == gameState[win[1]] && gameState[win[0]] == gameState[win[2]] {
			log.Println("[gameWon] Game won by: ", gameState[9])
			return true
		}
	}
	log.Println("[gameWon] Game not won")
	return false
}
func gameTied(gameState string) bool {
	for _, square := range gameState {
		if square == '.' {
			return false
		}
	}
	log.Println("[gameTied] Game tied")
	return true
}
