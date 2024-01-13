package src


import (
	"strconv"
	"chess-engine/src/api"
	"chess-engine/src/engine"
	"chess-engine/src/engine/moves"
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


func unflattenBoard(flattened [64]int) [8][8]int {
	var board [8][8]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			value := flattened[x * 8 + y]
			board[x][y] = value
		}
	}

	return board
}


func engineMove(currentPosition [64]int) {
	newBoard := engine.CalculateMove(currentPosition)
	unflattened := unflattenBoard(newBoard)

	//TODO: actually get the engine move
	moveObj := moves.Move{}

	api.WriteBoardState(unflattened, moveObj)
}


func legalMoves(board [64]int, json map[string]string, prevMove moves.Move) {
	x, err1 := strconv.Atoi(json["piece_x"])
	y, err2 := strconv.Atoi(json["piece_y"])

	if err1 != nil {
		panic(err1)
	} else if err2 != nil {
		panic(err2)
	}

	moves := engine.GetLegalMoves(board, x, y, prevMove)

	api.WriteLegalMoves(moves)
}


func Main() {
	json, parsedBoard, prevMove := api.LoadData()
	action := json["task"]
	flatBoard := flattenBoard(parsedBoard)

	//TODO: not have to precompute at the start of each move (store in a file)
	engine.PrecomputeValues()

	if action == "move_gen" {
		engineMove(flatBoard)
	} else if action == "legal_moves" {
		legalMoves(flatBoard, json, prevMove)
	}
}