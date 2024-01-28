package moves


type GameState struct {
	Board [64]int

	WhiteToMove bool

	WhiteKingCastle bool
	WhiteQueenCastle bool

	BlackKingCastle bool
	BlackQueenCastle bool

	PrevPawnDouble [2]int

	WhitePiecePos [][2]int
	BlackPiecePos [][2]int

	noKingMoveBitBoard uint64
	kingAttackBlocks []uint64

	pinArray [64]uint64
	enPassantPin bool
}


func CreateGameState(b [64]int, whiteMove bool, wkCastle bool, wqCastle bool, bkCastle bool, bqCastle bool, pDouble [2]int) GameState {
	//to be called whenever new game state obj is created

	var whitePiecePos [][2]int 
	var blackPiecePos [][2]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece := b[x * 8 + y]

			if piece != 0 {
				if piece < 7 {
					whitePiecePos = append(whitePiecePos, [2]int{x, y})
				} else {
					blackPiecePos = append(blackPiecePos, [2]int{x, y})
				}
			}
		}
	}

	state := GameState{Board: b, WhiteToMove: whiteMove, WhiteKingCastle: wkCastle, WhiteQueenCastle: wqCastle, BlackKingCastle: bkCastle, BlackQueenCastle: bqCastle, PrevPawnDouble: pDouble, WhitePiecePos: whitePiecePos, BlackPiecePos: blackPiecePos}

	kingVal := 11
	otherPieces := whitePiecePos
	if whiteMove {
		kingVal = 5
		otherPieces = blackPiecePos
	}

	kingX, kingY := getPiecePos(state, kingVal)
	kAttackBlock, pinArray, noKingMove, enPassantPin := getFilterBitboards(state.Board, kingX, kingY, kingVal, otherPieces, whiteMove, pDouble)

	state.noKingMoveBitBoard = noKingMove
	state.kingAttackBlocks = kAttackBlock
	state.pinArray = pinArray
	state.enPassantPin = enPassantPin

	return state
}


func getPiecePos(state GameState, pieceValue int) (int, int) {
	possible := state.WhitePiecePos
	if !state.WhiteToMove {
		possible = state.BlackPiecePos
	}

	for _, i := range possible {
		x := i[0]
		y := i[1]

		if state.Board[x * 8 + y] == pieceValue {
			return x, y
		}
	}

	return -1, -1
}