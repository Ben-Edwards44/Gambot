package engine


import (
	"fmt"
	"time"
	"strconv"
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
)


var fileNames [8]string = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}


func getMoveStr(move moves.Move) string {
	startRank := strconv.Itoa(8 - move.StartX)
	startFile := fileNames[move.StartY]
	endRank := strconv.Itoa(8 - move.EndX)
	endFile := fileNames[move.EndY]

	return startFile + startRank + endFile + endRank
}


func bulkCount(position *board.GameState, depth int) int {	
	moveList := moves.GenerateAllMoves(position, false)

	if depth == 1 {return len(moveList)}

	total := 0
	for _, i := range moveList {
		moves.MakeMove(position, i)
		total += bulkCount(position, depth - 1)
		moves.UnMakeLastMove(position)
	}

	return total
}


func dividePerft(stateObj *board.GameState, maxDepth int) int {
	initMoves := moves.GenerateAllMoves(stateObj, false)

	total := 0
	for _, i := range initMoves {
		str := getMoveStr(i)

		moves.MakeMove(stateObj, i)

		current := 1
		if maxDepth > 1 {
			current = bulkCount(stateObj, maxDepth - 1)
		} 

		total += current

		moves.UnMakeLastMove(stateObj)

		fmt.Print(str + ": ")
		fmt.Println(current)
	}

	return total
}


func Perft(stateObj *board.GameState, maxDepth int) {
	start := time.Now()
	
	nodes := dividePerft(stateObj, maxDepth)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Time taken: ")
	fmt.Println(elapsed)

	fmt.Print("Nodes searched: ")
	fmt.Println(nodes)
}