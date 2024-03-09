package evaluation


import "chess-engine/src/engine/board"


var pieceWeight [6]int = [6]int{10, 30, 30, 50, 100, 90}


func countMaterial(state *board.GameState) int {
	white := 0
	black := 0

	wPieces := &board.PieceLists.WhitePieceSquares
	bPieces := &board.PieceLists.BlackPieceSquares
	
	for i := 0; i < len(wPieces); i++ {
		wPos := wPieces[i]
		bPos := bPieces[i]

		if wPos != -1 {white += pieceWeight[state.Board[wPos] - 1]}
		if bPos != -1 {black += pieceWeight[state.Board[bPos] - 7]}
	}

	return white - black
}


func Eval(state *board.GameState, whiteToMove bool) int {
	//NOTE: with negamax, the eval should always be in the perspective of the current player (so times by -1 for black)
	material := countMaterial(state)

	perspective := -1
	if whiteToMove {perspective = 1}

	//eval = (good for white - good for black) * perspective

	return material * perspective
}