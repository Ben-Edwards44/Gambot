package search


import "chess-engine/src/engine/moves"


var pieceWeight [6]int = [6]int{10, 30, 30, 50, 100, 90}


func countMaterial(state *moves.GameState) int {
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


func eval(state *moves.GameState) int {
	material := countMaterial(state)

	return material
}