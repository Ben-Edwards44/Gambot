package engine


import (
	"fmt"
	"time"
	"chess-engine/src/engine/moves"
)


func Perft(stateObj moves.GameState, maxDepth int) {
	start := time.Now()
	
	nodes := bulkCount(&stateObj, maxDepth)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Nodes searched: ")
	fmt.Println(nodes)

	fmt.Print("Time taken: ")
	fmt.Println(elapsed)
}


func bulkCount(position *moves.GameState, depth int) int {	
	moveList := moves.GenerateAllMoves(*position)

	if depth == 1 {return len(moveList)}

	total := 0
	for _, i := range moveList {
		moves.MakeMove(position, i)
		total += bulkCount(position, depth - 1)
		moves.UnMakeLastMove(position)
	}

	return total
}