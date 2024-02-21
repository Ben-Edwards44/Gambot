package board


//https://stackoverflow.com/questions/10067514/correctly-implementing-zobrist-hashing


import "math/rand"


var ZobNums zobristNums


type zobristNums struct {
	PieceVals [64 * 6]uint64
	SideToMove uint64  //enabled if white to move
	CastlingRights [16]uint64
	EpFiles [8]uint64
}


func PrecalculateZobristNums() {
	//This should be called at the start of execution

	var pieceVals [64 * 6]uint64
	for i := 0; i < 64 * 6; i++ {
		num := rand.Uint64()
		pieceVals[i] = num
	}

	sideToMove := rand.Uint64()

	var castlingRights [16]uint64
	for i := 0; i < 16; i++ {
		num := rand.Uint64()
		castlingRights[i] = num
	}

	var epFiles [8]uint64
	for i := 0; i < 8; i++ {
		num := rand.Uint64()
		epFiles[i] = num
	}

	ZobNums = zobristNums{PieceVals: pieceVals, SideToMove: sideToMove, CastlingRights: castlingRights, EpFiles: epFiles}
}


func HashState(state *GameState) uint64 {
	//This should only be called once because it is slow. Use XORs after the initial generation.
	var hash uint64

	//hash board
	for i := 0; i < 64; i++ {
		pieceVal := state.Board[i]
		if pieceVal > 6 {pieceVal -= 6}

		if pieceVal > 0 {
			inx := i * 6 + (pieceVal - 1)
			hash ^= ZobNums.PieceVals[inx]
		}
	}

	hash ^= ZobNums.CastlingRights[state.CastleRights]  //hash castling

	//hash en passant file
	epFile := state.PrevPawnDouble[1]
	if epFile != -1 {
		hash ^= ZobNums.EpFiles[epFile]
	}

	//hash side to move
	if state.WhiteToMove {hash ^= ZobNums.SideToMove}

	return hash
}