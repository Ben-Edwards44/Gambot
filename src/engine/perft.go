package engine


import (
	"fmt"
	"time"
	"strconv"
	"chess-engine/src/engine/moves"
)


var fileNames [8]string = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}


func Perft(stateObj moves.GameState, maxDepth int) {
	start := time.Now()
	initMoves := moves.GenerateAllMoves(stateObj)

	var totals []int
	var movesFromStart []map[string]int
	for i := 0; i < maxDepth; i++ {
		totals = append(totals, 0)
		movesFromStart = append(movesFromStart, make(map[string]int))
	}

	for _, move := range initMoves {
		startRank := strconv.Itoa(8 - move.StartX)
		startFile := fileNames[move.StartY]
		endRank := strconv.Itoa(8 - move.EndX)
		endFile := fileNames[move.EndY]

		key := startFile + startRank + endFile + endRank

		//add the depth 1 move
		totals[0]++
		addToMap(key, 1, movesFromStart[0])

		new := makeMove(stateObj, move)
		prevPos := []moves.GameState{new}

		for depth := 1; depth < maxDepth; depth++ {
			prevPos = getMoves(prevPos)

			totals[depth] += len(prevPos)
			addToMap(key, len(prevPos), movesFromStart[depth])
		}
	}
	
	for i, x := range movesFromStart {
		fmt.Println("Depth " + strconv.Itoa(i + 1) + ":")
		
		for k, v := range x {
			fmt.Println(k + ": " + strconv.Itoa(v))
		}

		fmt.Println()
	}

	fmt.Println(totals)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Time taken: ")
	fmt.Println(elapsed)
}


func addToMap(key string, value int, m map[string]int) {
	_, exists := m[key]

	if exists {
		m[key] += value
	} else {
		m[key] = value
	}
}


func makeMove(currentState moves.GameState, move moves.Move) moves.GameState {
	newState := moves.MakeMoveCopy(currentState, move)
	updated := moves.CreateGameState(newState.Board, newState.WhiteToMove, newState.WhiteKingCastle, newState.WhiteQueenCastle, newState.BlackKingCastle, newState.BlackQueenCastle, newState.PrevPawnDouble)

	return updated
}


func getMoves(positions []moves.GameState) []moves.GameState {
	var newPositions []moves.GameState

	for _, i := range positions {
		moveList := moves.GenerateAllMoves(i)

		for _, x := range moveList {
			new := makeMove(i, x)
			newPositions = append(newPositions, new)
		}
	}

	return newPositions
}