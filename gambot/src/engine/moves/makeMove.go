package moves


import "gambot/src/engine/board"


func updateBitboards(state *board.GameState) {
	//TODO: make faster??

	kingVal := 11
	kingPos := board.PieceLists.BlackKingPos
	otherPieces := &board.PieceLists.WhitePieceSquares
	if state.WhiteToMove {
		kingVal = 5
		kingPos = board.PieceLists.WhiteKingPos
		otherPieces = &board.PieceLists.BlackPieceSquares
	}

	bitboards := board.Bitboard{}

	enPassantPin, doubleChecked := GetFilterBitboards(&state.Board, kingPos, kingVal, otherPieces, state.WhiteToMove, state.PrevPawnDouble, &bitboards)

	state.DoubleChecked = doubleChecked
	state.EnPassantPin = enPassantPin
	state.Bitboards = &bitboards
}


func updateHash(state *board.GameState, move *Move, start int, end int, pieceVal int, captVal int, prevCastle uint8, newCastle uint8, prevEpFile int) {
	//NOTE: this function should not really check the state.Board because, by this time, it will have been updated
	
	//adjust for white/black and for the fact that these values will be indices
	if pieceVal > 6 {
		pieceVal -= 7
	} else {
		pieceVal--
	}
	
	newZobHash := state.ZobristHash

	//Update hash with new piece pos
	newZobHash ^= board.ZobNums.PieceVals[start * 6 + pieceVal]  //get rid of piece from start square
	newZobHash ^= board.ZobNums.PieceVals[end * 6 + pieceVal]  //place piece on new square

	if captVal != 0 {
		//convert captVal to index
		if captVal > 6 {
			captVal -= 7
		} else {
			captVal--
		}

		newZobHash ^= board.ZobNums.PieceVals[end * 6 + captVal]
	}

	//update rook pos as well for castling
	if move.KingCastle {
		rookVal := pieceVal - 1

		newZobHash ^= board.ZobNums.PieceVals[(end + 1) * 6 + rookVal]
		newZobHash ^= board.ZobNums.PieceVals[(end - 1) * 6 + rookVal]
	} else if move.QueenCastle {
		rookVal := pieceVal - 1

		newZobHash ^= board.ZobNums.PieceVals[(end - 2) * 6 + rookVal]
		newZobHash ^= board.ZobNums.PieceVals[(end + 1) * 6 + rookVal]
	}

	//promotions
	if move.PromotionValue != 0 {
		pVal := move.PromotionValue
		if pVal > 6 {
			pVal -= 7
		} else {
			pVal--
		}

		newZobHash ^= board.ZobNums.PieceVals[end * 6 + pieceVal]  //remove pawn (which we moved to end space earlier)
		newZobHash ^= board.ZobNums.PieceVals[end * 6 + pVal]  //add new piece
	}

	//deal with ep passant capture
	if move.EnPassant {
		catpPos := move.StartX * 8 + move.EndY
		newZobHash ^= board.ZobNums.PieceVals[catpPos * 6]  //don't need to add pieceValue because it will be a pawn (so index 0)
	}

	//add en passant file (if needed)
	if move.DoublePawnMove {
		newZobHash ^= board.ZobNums.EpFiles[move.EndY]
	}

	//update castle rights. NOTE: the first 4 bits of the uint8 act as flags from white king/queen and black king/queen castling
	if prevCastle != newCastle {
		newZobHash ^= board.ZobNums.CastlingRights[prevCastle]
		newZobHash ^= board.ZobNums.CastlingRights[newCastle]
	}

	//get rid of the en passant target square from the previous move (if needed)
	if prevEpFile != -1 {
		newZobHash ^= board.ZobNums.EpFiles[prevEpFile]
	}

	newZobHash ^= board.ZobNums.SideToMove

	state.ZobristHash = newZobHash
}


func MakeMove(state *board.GameState, move *Move) {
	//updates game state

	state.SetPrevVals()  //so that we can restore later

	start := move.StartX * 8 + move.StartY
	end := move.EndX * 8 + move.EndY
	val := move.PieceValue
	captVal := state.Board[end]

	board.MovePiecePosition(start, end, move.PieceValue)

	//move piece
	state.Board[start] = 0
	state.Board[end] = val

	state.WhiteToMove = !state.WhiteToMove  //because we have just made a move

	if move.EnPassant {
		capturePos := move.StartX * 8 + move.EndY
		state.Board[capturePos] = 0

		board.EnPassant(capturePos, move.PieceValue > 6)
	} else if move.KingCastle {
		rookVal := move.PieceValue - 1

		state.Board[end + 1] = 0
		state.Board[end - 1] = rookVal

		board.Castle(end + 1, end - 1, rookVal)  //move the rook
	} else if move.QueenCastle {
		rookVal := move.PieceValue - 1

		state.Board[end - 2] = 0
		state.Board[end + 1] = rookVal

		board.Castle(end - 2, end + 1, rookVal)  //move the rook
	}
	
	oldEpFile := state.PrevPawnDouble[1]
	if move.DoublePawnMove {
		state.PrevPawnDouble = [2]int{move.EndX, move.EndY}
	} else {
		state.PrevPawnDouble = [2]int{-1, -1}
	}

	if move.PromotionValue != 0 {
		state.Board[end] = move.PromotionValue
	}

	newCastleRights := state.CastleRights
	if move.PieceValue == 5 {
		//white king moving
		newCastleRights &= board.InvWkCastle
		newCastleRights &= board.InvWqCastle
	} else if move.PieceValue == 11 {
		//black king moving
		newCastleRights &= board.InvBkCastle
		newCastleRights &= board.InvBqCastle
	} else if move.PieceValue == 4 {
		//white rook moving
		if move.StartY == 7 {
			newCastleRights &= board.InvWkCastle
		} else if move.StartY == 0 {
			newCastleRights &= board.InvWqCastle
		}
	} else if move.PieceValue == 10 {
		//black rook moving
		if move.StartY == 7 {
			newCastleRights &= board.InvBkCastle
		} else if move.StartY == 0 {
			newCastleRights &= board.InvBqCastle
		}
	}

	oldCastleRights := state.CastleRights
	state.CastleRights = newCastleRights

	updateBitboards(state)
	updateHash(state, move, start, end, val, captVal, oldCastleRights, newCastleRights, oldEpFile)
}


func UnMakeLastMove(state *board.GameState) {
	state.RestorePrev()
	board.UnMoveLastPiece()
}


func CreateGameState(b [64]int, whiteMove bool, castleRights uint8, pDouble [2]int) board.GameState {
	//to be called whenever new game state obj is created
	state := board.GameState{Board: b, WhiteToMove: whiteMove, CastleRights: castleRights, PrevPawnDouble: pDouble}
	bitboards := board.Bitboard{}

	board.InitPieceLists(&state)
	
	kingVal := 11
	kingPos := board.PieceLists.BlackKingPos
	otherPieces := &board.PieceLists.WhitePieceSquares
	if whiteMove {
		kingVal = 5
		kingPos = board.PieceLists.WhiteKingPos
		otherPieces = &board.PieceLists.BlackPieceSquares
	}

	enPassantPin, doubleChecked := GetFilterBitboards(&state.Board, kingPos, kingVal, otherPieces, whiteMove, pDouble, &bitboards)

	state.DoubleChecked = doubleChecked
	state.EnPassantPin = enPassantPin
	state.Bitboards = &bitboards

	zobHash := board.HashState(&state)
	state.ZobristHash = zobHash

	return state
}