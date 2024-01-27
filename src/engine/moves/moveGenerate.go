package moves


func GenerateAllMoves(state GameState) []Move {
	//assumes state has been properly initialised etc.

	piecePos := state.BlackPiecePos
	if state.WhiteToMove {
		piecePos = state.WhitePiecePos
	}
	
	var moves []Move
	for _, i := range piecePos {	
		GetPieceMoves(state, i[0], i[1], &moves)
	}

	return moves
}