package moves


//import "fmt"


func getPseudoLegalMoves(state GameState, movesForWhite bool) []Move {
	piecePos := state.BlackPiecePos
	if movesForWhite {
		piecePos = state.WhitePiecePos
	}
	
	var moves []Move
	for _, i := range piecePos {
		//other bit board will default to 0, which is fine since we need to ignore it here
		GetPieceMoves(state, i[0], i[1], &moves)

		//fmt.Println(i)

		//if i[0] == 0 && i[1] == 3 {
		//	fmt.Println(moves)
		//	var a []Move
		//	//fmt.Println(state)
		//	GetPieceMoves(state, 0, 3, &a)
		//	//fmt.Println(a)
		//}
	}

	return moves
}


func getOtherMoveBitBoard(state GameState) uint64 {
	moves := getPseudoLegalMoves(state, !state.WhiteToMove)
	bitBoard := movesToBitBoard(moves)

	return bitBoard
}