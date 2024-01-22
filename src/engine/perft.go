package engine


import (
	"fmt"
	"strconv"
	"chess-engine/src/engine/moves"
)


var moveCounts map[int]int
var donePos []map[[64]int]bool


func Perft(stateObj moves.GameState, depth int) {
	moveCounts = make(map[int]int)

	for i := 0; i <= depth; i++ {
		m := make(map[[64]int]bool)
		donePos = append(donePos, m)
	}

	recordMoves(stateObj, 1, depth)

	for k, v := range moveCounts {
		str := "Depth " + strconv.Itoa(k) + " : " + strconv.Itoa(v)
		
		fmt.Println(str)
	}
}


func recordMoves(state moves.GameState, depth int, maxDepth int) {
	if depth > maxDepth {return}

	doneMap := donePos[depth]
	moveList := moves.GenerateAllMoves(state)

	_, exists := moveCounts[depth]
	if !exists {moveCounts[depth] = 0}

	for _, i := range moveList {
		newState := moves.MakeMoveCopy(state, i)

		_, doneMove := doneMap[newState.Board]

		if !doneMove {
			moveCounts[depth] += 1
			doneMap[newState.Board] = true
		}

		recordMoves(newState, depth + 1, maxDepth)
	}
}