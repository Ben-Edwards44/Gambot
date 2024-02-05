package src


import (
	"os"
	"strconv"
	"runtime/pprof"
	"chess-engine/src/api"
	"chess-engine/src/engine"
	"chess-engine/src/engine/moves"
)


func engineMove(stateObj *moves.GameState) {
	newState := engine.CalculateMove(stateObj)

	api.WriteState(newState)
}


func legalMoves(stateObj *moves.GameState, json map[string]string) {
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


func perft(stateObj *moves.GameState, json map[string]string) {
	depth, err := strconv.Atoi(json["perft_depth"])

	if err != nil {panic(err)}

	file, err := os.Create("profile.prof")

	if err != nil {panic(err)}

	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	//use: go tool pprof -http=:8080 profile.prof

	engine.Perft(stateObj, depth)
}


func Main() {
	//TODO: not have to precompute at the start of each move (store in a file)
	engine.PrecomputeValues()
	
	json, stateObj := api.LoadGameState()
	action := json["task"]

	if action == "move_gen" {
		engineMove(&stateObj)
	} else if action == "legal_moves" {
		legalMoves(&stateObj, json)
	} else if action == "perft" {
		perft(&stateObj, json)
	}
}