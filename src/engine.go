package src


import (
	"strconv"
	"chess-engine/src/api"
	"chess-engine/src/engine"
)


func engineMove(stateObj engine.GameState) {
	newState := engine.CalculateMove(stateObj)

	api.WriteState(newState)
}


func legalMoves(stateObj engine.GameState, json map[string]string) {
	x, err1 := strconv.Atoi(json["piece_x"])
	y, err2 := strconv.Atoi(json["piece_y"])

	if err1 != nil {
		panic(err1)
	} else if err2 != nil {
		panic(err2)
	}

	moves := engine.GetLegalMoves(stateObj, x, y)

	api.WriteLegalMoves(moves)
}


func Main() {
	json, stateObj := api.LoadGameState()
	action := json["task"]

	//TODO: not have to precompute at the start of each move (store in a file)
	engine.PrecomputeValues()

	if action == "move_gen" {
		engineMove(stateObj)
	} else if action == "legal_moves" {
		legalMoves(stateObj, json)
	}
}