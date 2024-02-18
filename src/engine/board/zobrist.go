package board


import "math/rand"


type zobristNums struct {
	pieceVals [64 * 6]uint64
	sideToMove uint64
	castlingRights [4]uint64
	epFiles [8]uint64
}


func generateZobristNums() zobristNums {
	//This should be called at the start of execution

	var pieceVals [64 * 6]uint64
	for i := 0; i < 64 * 6; i++ {
		num := rand.Uint64()
		pieceVals[i] = num
	}

	sideToMove := rand.Uint64()

	var castlingRights [4]uint64
	for i := 0; i < 4; i++ {
		num := rand.Uint64()
		castlingRights[i] = num
	}

	var epFiles [8]uint64
	for i := 0; i < 8; i++ {
		num := rand.Uint64()
		epFiles[i] = num
	}

	obj := zobristNums{pieceVals: pieceVals, sideToMove: sideToMove, castlingRights: castlingRights, epFiles: epFiles}

	return obj
}


func hashState(state *GameState, nums zobristNums) uint64 {
	//This should only be called once because it is slow. Use XORs after the initial generation.
	var hash uint64

	//hash board
	for i := 0; i < 64; i++ {
		pieceVal := state.Board[i]
		if pieceVal > 6 {pieceVal -= 6}

		if pieceVal > 0 {
			inx := i * 6 + (pieceVal - 1)
			hash ^= nums.pieceVals[inx]
		}
	}

	//hash castling
	if state.WhiteKingCastle {hash ^= nums.castlingRights[0]}
	if state.WhiteQueenCastle {hash ^= nums.castlingRights[1]}
	if state.BlackKingCastle {hash ^= nums.castlingRights[2]}
	if state.BlackQueenCastle {hash ^= nums.castlingRights[3]}

	//hash en passant file
	epFile := state.PrevPawnDouble[1]
	if epFile != -1 {
		hash ^= nums.epFiles[epFile]
	}

	return hash
}