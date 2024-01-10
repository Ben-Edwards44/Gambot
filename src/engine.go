package src


import (
	"chess-engine/src/api"
	"chess-engine/src/engine"
)


func EngineMove() {
	currentPosition := api.LoadBoardState()
	newBoard := engine.CalculateMove(currentPosition)
	api.WriteBoardState(newBoard)
}


func Test() {
	api.Test()
}