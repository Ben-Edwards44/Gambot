package moves


func GetAllMoves(board [64]int, prevPawnDouble [2]int, piecePos [][2]int) []Move {
	var moves []Move
	for _, i := range piecePos {
		GetPieceMoves(board, i[0], i[1], prevPawnDouble, &moves)
	}

	return moves
}