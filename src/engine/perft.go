package engine


import (
	"fmt"
	"strconv"
	"chess-engine/src/engine/moves"
)


func Perft(stateObj moves.GameState, depth int) {
	prevPos := []moves.GameState{stateObj}

	for i := 1; i <= depth; i++ {
		prevPos = getMoves(prevPos)

		d := strconv.Itoa(i)
		total := strconv.Itoa(len(prevPos))

		fmt.Println("Depth " + d + " : " + total)
	}
}


func getMoves(positions []moves.GameState) []moves.GameState {
	var newPositions []moves.GameState

	for _, i := range positions {
		moveList := moves.GenerateAllMoves(i)

		for _, x := range moveList {
			newState := moves.MakeMoveCopy(i, x)
			updated := moves.CreateGameState(newState.Board, newState.WhiteToMove, newState.WhiteKingCastle, newState.WhiteQueenCastle, newState.BlackKingCastle, newState.BlackQueenCastle, newState.PrevPawnDouble)

			newPositions = append(newPositions, updated)
		}
	}

	return newPositions
}