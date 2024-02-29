package evaluation


import "chess-engine/src/engine/board"


var pieceWeight [6]int = [6]int{10, 30, 30, 50, 100, 90}


func countMaterial(state *board.GameState) int {
	white := 0
	black := 0
	for inx := 0; inx < 6; inx++ {
		for _, i := range state.WhitePiecePos[inx] {
			if i[0] == -1 {break}  //because we are using fixed length array
			white += pieceWeight[inx]
		}

		for _, i := range state.BlackPiecePos[inx] {
			if i[0] == -1 {break}  //because we are using fixed length array
			black += pieceWeight[inx]
		}
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