package search


import "chess-engine/src/engine/moves"


var pieceWeight [6]int = [6]int{10, 30, 30, 50, 100, 90}


func countMaterial(state *moves.GameState) int {
	white := 0
	black := 0
	for inx := 0; inx < 6; inx++ {
		for i := 0; i < len(state.WhitePiecePos[inx]); i++ {
			white += pieceWeight[inx]
		}

		for i := 0; i < len(state.BlackPiecePos[inx]); i++ {
			black += pieceWeight[inx]
		}
	}

	return white - black
}


func eval(state *moves.GameState) int {
	material := countMaterial(state)

	return material
}