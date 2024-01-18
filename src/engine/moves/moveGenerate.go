package moves


func getPseudoLegalMoves(state GameState, movesForWhite bool) []Move {
	piecePos := state.BlackPiecePos
	if movesForWhite {
		piecePos = state.WhitePiecePos
	}
	
	var moves []Move
	for _, i := range piecePos {		
		GetPieceMoves(state, i[0], i[1], &moves)
	}

	return moves
}