package moves


//Not much to see here for now, but I plan to use bitboards a lot more later


func setBitBoard(bitboard *uint64, pos int) {
	//set the inx of a bitboard to a 1
	bitWeight := uint64(1 << uint64(pos))  //2**n
	*bitboard |= bitWeight
}