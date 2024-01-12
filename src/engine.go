package src


import (
	"chess-engine/src/api"
	"chess-engine/src/engine"
)


func engineMove(currentPosition [8][8]int) {
	newBoard := engine.CalculateMove(currentPosition)
	api.WriteBoardState(newBoard)
}


func Main() {
	json, parsedBoard := api.LoadData()
	action := json["task"]

	if action == "move_gen" {
		engineMove(parsedBoard)
	} else if action == "legal_moves" {
		//TODO: actually generate legal moves
		moves := [][2]int{{0, 0}, {1, 1}, {2, 2}}
		api.WriteLegalMoves(moves)
	}
}