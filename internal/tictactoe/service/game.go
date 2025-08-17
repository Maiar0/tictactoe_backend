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
