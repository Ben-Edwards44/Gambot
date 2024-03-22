package engine


import (
	"os"
	"runtime/pprof"
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/search"
)


func CalculateMove(stateObj *board.GameState, moveTime int) *moves.Move {
	//NOTE: UCI will handle updating board
	
	//start profiling (go tool pprof -http=:8080 profile.prof)
	file, err := os.Create("profile.prof")
	if err != nil {panic(err)}

	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
	
	move := search.GetBestMove(stateObj, moveTime)
	
	return move
}