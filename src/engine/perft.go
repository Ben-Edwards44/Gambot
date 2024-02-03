package engine


import (
	"fmt"
	"time"
	"chess-engine/src/engine/moves"
)


func Perft(stateObj moves.GameState, maxDepth int) {
	//test(&stateObj, maxDepth)
	//return
	start := time.Now()
	
	nodes := bulkCount(&stateObj, maxDepth)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Nodes searched: ")
	fmt.Println(nodes)

	fmt.Print("Time taken: ")
	fmt.Println(elapsed)
}


func test(state *moves.GameState, depth int) {
	for i := 0; i < depth; i++ {
		m := moves.GenerateAllMoves(*state)
		moves.MakeMove(state, m[0])

		fmt.Println(i)
		fmt.Println(state.PrvWhitePiecePos)
		fmt.Println(state.PrvBlackPiecePos)
		fmt.Print("\n\n")
	}
	for i := 0; i < depth; i++ {
		moves.UnMakeLastMove(state)

		fmt.Println(i)
		fmt.Println(state.WhitePiecePos)
		fmt.Println(state.BlackPiecePos)
		fmt.Print("\n\n")
	}
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