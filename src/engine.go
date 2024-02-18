package src

import (
	"chess-engine/src/api"
	"chess-engine/src/engine"
	"chess-engine/src/engine/board"
	"os"
	"runtime/pprof"
	"strconv"	
)


func engineMove(stateObj *board.GameState) {
	file, err := os.Create("profile.prof")

	if err != nil {panic(err)}

	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	newState := engine.CalculateMove(stateObj)

	api.WriteState(newState)
}


func legalMoves(stateObj *board.GameState, json map[string]string) {
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


func perft(stateObj *board.GameState, json map[string]string) {
	depth, err := strconv.Atoi(json["perft_depth"])
	test := json["perft_test"] == "true"

	if err != nil {panic(err)}

	//use: go tool pprof -http=:8080 profile.prof

	engine.Perft(stateObj, depth, test)
}


func checkWin(stateObj *board.GameState) {
	writeValue := engine.CheckWin(stateObj)

	api.WriteCheckWin(writeValue)
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
	} else if action == "check_win" {
		checkWin(&stateObj)
	} else {
		panic("Invalid task from API")
	}
}