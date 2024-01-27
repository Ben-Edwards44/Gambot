package moves


func setBitBoard(bitboard *uint64, pos int) {
	//set the inx of a bitboard to a 1
	bitWeight := uint64(1 << uint64(pos))  //2**n
	*bitboard |= bitWeight
}