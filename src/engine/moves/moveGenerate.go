package moves


func addPawnCaptures(x int, y int, pieceValue int, movesForWhite bool, moves *[]Move) {
	newX := x + 1
	if movesForWhite {
		newX = x - 1
	}

	if newX < 0 || newX > 7 {return}

	for yStep := -1; yStep < 2; yStep += 2 {
		newY := y + yStep
		if 0 <= newY && newY < 8 {
			move := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
			*moves = append(*moves, move)
		}
	}
}


func getNoKingMoves(state GameState, movesForWhite bool) []Move {
	//find where king cannot move (because it would move into check)

	piecePos := state.BlackPiecePos
	if movesForWhite {
		piecePos = state.WhitePiecePos
	}
	
	var moves []Move
	for _, i := range piecePos {	
		pieceValue := state.Board[i[0] * 8 + i[1]]

		if pieceValue != 1 && pieceValue != 7 {
			//not a pawn so all moves are potential captures
			GetPieceMoves(state, i[0], i[1], &moves)
		} else {
			//we only want to include pawn captures
			addPawnCaptures(i[0], i[1], pieceValue, movesForWhite, &moves)
		}
	}

	return moves
}


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