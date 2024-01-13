package src


import (
	"strconv"
	"chess-engine/src/api"
	"chess-engine/src/engine"
)


func flattenBoard(board [8][8]int) [64]int {
	var flattened [64]int

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			flattened[x * 8 + y] = board[x][y]
		}
	}

	return flattened
}


func engineMove(currentPosition [8][8]int) {
	newBoard := engine.CalculateMove(currentPosition)
	api.WriteBoardState(newBoard)
}


func legalMoves(board [64]int, json map[string]string) {
	x, err1 := strconv.Atoi(json["piece_x"])
	y, err2 := strconv.Atoi(json["piece_y"])

	if err1 != nil {
		panic(err1)
	} else if err2 != nil {
		panic(err2)
	}

	moves := engine.GetLegalMoves(board, x, y)

	api.WriteLegalMoves(moves)
}


func Main() {
	json, parsedBoard := api.LoadData()
	action := json["task"]
	flatBoard := flattenBoard(parsedBoard)

	//TODO: not have to precompute at the start of each move (store in a file)
	engine.PrecomputeValues()

	if action == "move_gen" {
		//TODO: use flattened board
		engineMove(parsedBoard)
	} else if action == "legal_moves" {
		//TODO: actually generate legal moves
		legalMoves(flatBoard, json)
	}
}