package moves


type GameState struct {
	Board [64]int

	WhiteToMove bool

	WhiteKingCastle bool
	WhiteQueenCastle bool

	BlackKingCastle bool
	BlackQueenCastle bool

	PrevPawnDouble [2]int

	WhitePiecePos [6][][2]int
	BlackPiecePos [6][][2]int

	NoKingMoveBitBoard uint64
	kingAttackBlocks []uint64

	pinArray [64]uint64
	enPassantPin bool

	//for unmaking moves (end index is the most recent)
	prvBoard [][64]int

	prvWhiteToMove []bool

	prvWhiteKingCastle []bool
	prvWhiteQueenCastle []bool
	prvBlackKingCastle []bool
	prvBlackQueenCastle []bool

	prvPrevPawnDouble [][2]int

	PrvWhitePiecePos [][6][][2]int
	PrvBlackPiecePos [][6][][2]int

	prvNoKingMoveBitBoard []uint64
	prvKingAttackBlocks [][]uint64

	prvPinArray [][64]uint64
	prvEnPassantPin []bool
}


func (state *GameState) SetPrevVals() {
	//slices are passed by reference, so need to be manually copied (even nested slices)

	var cpyWhitePiecePos [6][][2]int
	for i, x := range state.WhitePiecePos {
		dst := make([][2]int, len(x))
		copy(dst, x)

		cpyWhitePiecePos[i] = dst
	}

	var cpyBlackPiecePos [6][][2]int
	for i, x := range state.BlackPiecePos {
		dst := make([][2]int, len(x))
		copy(dst, x)

		cpyBlackPiecePos[i] = dst
	}

	cpyKingAttackBlocks := make([]uint64, len(state.kingAttackBlocks))
	copy(cpyKingAttackBlocks, state.kingAttackBlocks)

	//set prev values
	state.prvBoard = append(state.prvBoard, state.Board)
	state.prvWhiteToMove = append(state.prvWhiteToMove, state.WhiteToMove)
	state.prvWhiteKingCastle = append(state.prvWhiteKingCastle, state.WhiteKingCastle)
	state.prvWhiteQueenCastle = append(state.prvWhiteQueenCastle, state.WhiteQueenCastle)
	state.prvBlackKingCastle = append(state.prvBlackKingCastle, state.BlackKingCastle)
	state.prvBlackQueenCastle = append(state.prvBlackQueenCastle, state.BlackQueenCastle)
	state.prvPrevPawnDouble = append(state.prvPrevPawnDouble, state.PrevPawnDouble)
	state.PrvWhitePiecePos = append(state.PrvWhitePiecePos, cpyWhitePiecePos)
	state.PrvBlackPiecePos = append(state.PrvBlackPiecePos, cpyBlackPiecePos)
	state.prvNoKingMoveBitBoard = append(state.prvNoKingMoveBitBoard, state.NoKingMoveBitBoard)
	state.prvKingAttackBlocks = append(state.prvKingAttackBlocks, cpyKingAttackBlocks)
	state.prvPinArray = append(state.prvPinArray, state.pinArray)
	state.prvEnPassantPin = append(state.prvEnPassantPin, state.enPassantPin)
}


func (state *GameState) RestorePrev() {
	//restore the previous values. Note that the slices are only shallow copies
	
	//restore value
	state.Board = state.prvBoard[len(state.prvBoard) - 1]
	state.WhiteToMove = state.prvWhiteToMove[len(state.prvWhiteToMove) - 1]
	state.WhiteKingCastle = state.prvWhiteKingCastle[len(state.prvWhiteKingCastle) - 1]
	state.WhiteQueenCastle = state.prvWhiteQueenCastle[len(state.prvWhiteQueenCastle) - 1]
	state.BlackKingCastle = state.prvBlackKingCastle[len(state.prvBlackKingCastle) - 1]
	state.BlackQueenCastle = state.prvBlackQueenCastle[len(state.prvBlackQueenCastle) - 1]
	state.PrevPawnDouble = state.prvPrevPawnDouble[len(state.prvPrevPawnDouble) - 1]
	state.WhitePiecePos = state.PrvWhitePiecePos[len(state.PrvWhitePiecePos) - 1]//cpyWhitePiecePos
	state.BlackPiecePos = state.PrvBlackPiecePos[len(state.PrvBlackPiecePos) - 1]//cpyBlackPiecePos
	state.NoKingMoveBitBoard = state.prvNoKingMoveBitBoard[len(state.prvNoKingMoveBitBoard) - 1]
	state.kingAttackBlocks = state.prvKingAttackBlocks[len(state.prvKingAttackBlocks) - 1]
	state.pinArray = state.prvPinArray[len(state.prvPinArray) - 1]
	state.enPassantPin = state.prvEnPassantPin[len(state.prvEnPassantPin) - 1]

	//pop end of slice
	state.prvBoard = state.prvBoard[:len(state.prvBoard) - 1]
	state.prvWhiteToMove = state.prvWhiteToMove[:len(state.prvWhiteToMove) - 1]
	state.prvWhiteKingCastle = state.prvWhiteKingCastle[:len(state.prvWhiteKingCastle) - 1]
	state.prvWhiteQueenCastle = state.prvWhiteQueenCastle[:len(state.prvWhiteQueenCastle) - 1]
	state.prvBlackKingCastle = state.prvBlackKingCastle[:len(state.prvBlackKingCastle) - 1]
	state.prvBlackQueenCastle = state.prvBlackQueenCastle[:len(state.prvBlackQueenCastle) - 1]
	state.prvPrevPawnDouble = state.prvPrevPawnDouble[:len(state.prvPrevPawnDouble) - 1]
	state.PrvWhitePiecePos = state.PrvWhitePiecePos[:len(state.PrvWhitePiecePos) - 1]
	state.PrvBlackPiecePos = state.PrvBlackPiecePos[:len(state.PrvBlackPiecePos) - 1]
	state.prvNoKingMoveBitBoard = state.prvNoKingMoveBitBoard[:len(state.prvNoKingMoveBitBoard) - 1]
	state.prvKingAttackBlocks = state.prvKingAttackBlocks[:len(state.prvKingAttackBlocks) - 1]
	state.prvPinArray = state.prvPinArray[:len(state.prvPinArray) - 1]
	state.prvEnPassantPin = state.prvEnPassantPin[:len(state.prvEnPassantPin) - 1]
}


func CreateGameState(b [64]int, whiteMove bool, wkCastle bool, wqCastle bool, bkCastle bool, bqCastle bool, pDouble [2]int) GameState {
	//to be called whenever new game state obj is created

	var whitePiecePos [6][][2]int
	var blackPiecePos [6][][2]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece := b[x * 8 + y]

			if piece != 0 {
				pos := [2]int{x, y}

				if piece < 7 {
					whitePiecePos[piece - 1] = append(whitePiecePos[piece - 1], pos)
				} else {
					blackPiecePos[piece - 7] = append(blackPiecePos[piece - 7], pos)
				}
			}
		}
	}

	state := GameState{Board: b, WhiteToMove: whiteMove, WhiteKingCastle: wkCastle, WhiteQueenCastle: wqCastle, BlackKingCastle: bkCastle, BlackQueenCastle: bqCastle, PrevPawnDouble: pDouble, WhitePiecePos: whitePiecePos, BlackPiecePos: blackPiecePos}

	kingVal := 11
	kingPos := blackPiecePos[4][0]  //4 not 5 because we convert from piece value to index
	otherPieces := whitePiecePos
	if whiteMove {
		kingVal = 5
		kingPos = whitePiecePos[4][0]  //4 not 5 because we convert from piece value to index
		otherPieces = blackPiecePos
	}

	kingX := kingPos[0]
	kingY := kingPos[1]
	kAttackBlock, pinArray, noKingMove, enPassantPin := getFilterBitboards(&state.Board, kingX, kingY, kingVal, otherPieces, whiteMove, pDouble)

	state.NoKingMoveBitBoard = noKingMove
	state.kingAttackBlocks = kAttackBlock
	state.pinArray = pinArray
	state.enPassantPin = enPassantPin

	return state
}