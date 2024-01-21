package moves


func setBitBoard(bitboard *uint64, pos int) {
	//set the inx of a bitboard to a 1
	bitWeight := uint64(1 << uint64(pos))  //2**n
	*bitboard |= bitWeight
}


func movesToBitBoard(moves []Move) uint64 {
	var bitBoard uint64
	for _, i := range moves {
		pos := i.EndX * 8 + i.EndY
		setBitBoard(&bitBoard, pos)
	}

	return bitBoard
}


func getOtherMoveBitBoard(state GameState) (uint64, []Move) {
	moves := getNoKingMoves(state, !state.WhiteToMove)
	bitBoard := movesToBitBoard(moves)

	return bitBoard, moves
}