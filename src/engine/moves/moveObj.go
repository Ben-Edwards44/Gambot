package moves


import "chess-engine/src/engine/board"


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//special moves
	DoublePawnMove bool
	EnPassant bool
	KingCastle bool
	QueenCastle bool
	PromotionValue int
}


func removeFromArray(arr [10][2]int, x int, y int) [10][2]int {
	removed := false
    for i, j := range arr {
		if j[0] == -1 {
			break
		} else if j[0] == x && j[1] == y {
			//remove element
			arr[i][0] = -1
			arr[i][1] = -1
			removed = true
		} else if removed {
			//shift subsequent elements left
			arr[i - 1][0] = j[0]
			arr[i - 1][1] = j[1]
			arr[i][0] = -1
			arr[i][1] = -1
		}
	}

	return arr
}


func appendToArray(arr [10][2]int, x int, y int) [10][2]int {
	for i, j := range arr {
		if j[0] == -1 {
			arr[i][0] = x
			arr[i][1] = y
			break
		}
	}

	return arr
}


func movePos(arr [10][2]int, sX int, sY int, eX int, eY int) [10][2]int {
	for i, x := range arr {
		if x[0] == sX && x[1] == sY {
			//update the piece to the new position
			arr[i][0] = eX
			arr[i][1] = eY
			break
		}
	}

	return arr
}


func updateCapture(state *board.GameState, move Move, ePos int, isWhite bool) {
	//TODO: work with fixed length array
	eVal := state.Board[ePos]
	var enemy [6][10][2]int
	var enemyInx int
	if isWhite {
		enemy = state.BlackPiecePos
		enemyInx = eVal - 7
	} else {
		enemy = state.WhitePiecePos
		enemyInx = eVal - 1
	}

	captX := -1
	captY := -1
	if eVal != 0 {
		captX = move.EndX
		captY = move.EndY
	} else if move.EnPassant {
		captX = move.StartX
		captY = move.EndY

		enemyInx = 0  //the en passant capture must be a pawn
	}

	//if needed, remove position of captured piece
	if captX != -1 {
		enemy[enemyInx] = removeFromArray(enemy[enemyInx], captX, captY)

		if isWhite {
			state.BlackPiecePos = enemy
		} else {
			state.WhitePiecePos = enemy
		}
	}
}


func updatePiecePos(move Move, sPos int, ePos int, sVal int, state *board.GameState) {
	isWhite := sVal < 7

	if move.PromotionValue == 0 {
		//not a promotion
		var friend [6][10][2]int
		var friendInx int
		if isWhite {
			friend = state.WhitePiecePos
			friendInx = sVal - 1
		} else {
			friend = state.BlackPiecePos
			friendInx = sVal - 7
		}

		friend[friendInx] = movePos(friend[friendInx], move.StartX, move.StartY, move.EndX, move.EndY)  //move piece

		if move.KingCastle {
			friend[3] = movePos(friend[3], move.StartX, 7, move.EndX, 5)  //move the rook as well
		} else if move.QueenCastle {
			friend[3] = movePos(friend[3], move.StartX, 0, move.EndX, 3)  //move the rook as well
		}

		if isWhite {
			state.WhitePiecePos = friend
		} else {
			state.BlackPiecePos = friend
		}
	} else {
		//promotion
		if isWhite {
			//remove the pawn
			state.WhitePiecePos[0] = removeFromArray(state.WhitePiecePos[0], move.StartX, move.StartY)

			//add the new piece
			state.WhitePiecePos[move.PromotionValue - 1] = appendToArray(state.WhitePiecePos[move.PromotionValue - 1], move.EndX, move.EndY)
		} else {
			//remove the pawn
			state.BlackPiecePos[0] = removeFromArray(state.BlackPiecePos[0], move.StartX, move.StartY)

			//add the new piece
			state.BlackPiecePos[move.PromotionValue - 7] = appendToArray(state.BlackPiecePos[move.PromotionValue - 7], move.EndX, move.EndY)
		}
	}

	updateCapture(state, move, ePos, isWhite)
} 


func updateBitboards(state *board.GameState) {
	//TODO: make faster??

	kingVal := 11
	kingPos := state.BlackPiecePos[4][0]  //4 not 5 because we are converting from piece value to index
	otherPieces := state.WhitePiecePos
	if state.WhiteToMove {
		kingVal = 5
		kingPos = state.WhitePiecePos[4][0]  //4 not 5 because we are converting from piece value to index
		otherPieces = state.BlackPiecePos
	}

	kingX := kingPos[0]
	kingY := kingPos[1]
	kAttackBlock, pinArray, noKingMove, enPassantPin := GetFilterBitboards(&state.Board, kingX, kingY, kingVal, otherPieces, state.WhiteToMove, state.PrevPawnDouble)

	state.NoKingMoveBitBoard = noKingMove
	state.KingAttackBlocks = kAttackBlock
	state.PinArray = pinArray
	state.EnPassantPin = enPassantPin
}


func updateHash(state *board.GameState, move Move, start int, end int, pieceVal int, captVal int, prevCastle uint8, newCastle uint8, prevEpFile int) {
	newZobHash := state.ZobristHash

	//Update hash with new piece pos
	newZobHash ^= board.ZobNums.PieceVals[start * 6 + pieceVal]  //get rid of piece from start square
	newZobHash ^= board.ZobNums.PieceVals[end * 6 + pieceVal]  //place piece on new square

	if captVal != 0 {
		newZobHash ^= board.ZobNums.PieceVals[end * 6 + captVal]
	}

	//update rook pos as well for castling
	if move.KingCastle {
		rookVal := move.PieceValue - 1

		newZobHash ^= board.ZobNums.PieceVals[(end + 1) * 6 + rookVal]
		newZobHash ^= board.ZobNums.PieceVals[(end - 1) * 6 + rookVal]
	} else if move.QueenCastle {
		rookVal := move.PieceValue - 1

		newZobHash ^= board.ZobNums.PieceVals[(end - 2) * 6 + rookVal]
		newZobHash ^= board.ZobNums.PieceVals[(end + 1) * 6 + rookVal]
	}

	//promotions
	if move.PromotionValue != 0 {
		newZobHash ^= board.ZobNums.PieceVals[start * 6 + pieceVal]  //remove pawn
		newZobHash ^= board.ZobNums.PieceVals[end * 6 + move.PromotionValue]  //add new piece
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


func MakeMove(state *board.GameState, move Move) {
	//updates game state

	state.SetPrevVals()  //so that we can restore later

	start := move.StartX * 8 + move.StartY
	end := move.EndX * 8 + move.EndY
	val := move.PieceValue
	//captVal := state.Board[end]

	updatePiecePos(move, start, end, val, state)

	//move piece
	state.Board[start] = 0
	state.Board[end] = val

	state.WhiteToMove = !state.WhiteToMove  //because we have just made a move

	if move.EnPassant {
		capturePos := move.StartX * 8 + move.EndY
		state.Board[capturePos] = 0
	} else if move.KingCastle {
		rookVal := move.PieceValue - 1

		state.Board[end + 1] = 0
		state.Board[end - 1] = rookVal
	} else if move.QueenCastle {
		rookVal := move.PieceValue - 1

		state.Board[end - 2] = 0
		state.Board[end + 1] = rookVal
	}
	
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
		newCastleRights &= board.InvWkCastle //state.WhiteKingCastle = false
		newCastleRights &= board.InvWqCastle //state.WhiteQueenCastle = false
	} else if move.PieceValue == 11 {
		//black king moving
		newCastleRights &= board.InvBkCastle //state.BlackKingCastle = false
		newCastleRights &= board.InvBqCastle //state.BlackQueenCastle = false
	} else if move.PieceValue == 4 {
		//white rook moving
		if move.StartY == 7 {
			newCastleRights &= board.InvWkCastle //state.WhiteKingCastle = false
		} else if move.StartY == 0 {
			newCastleRights &= board.InvWqCastle //state.WhiteQueenCastle = false
		}
	} else if move.PieceValue == 10 {
		//black rook moving
		if move.StartY == 7 {
			newCastleRights &= board.InvBkCastle //state.BlackKingCastle = false
		} else if move.StartY == 0 {
			newCastleRights &= board.InvBqCastle //state.BlackQueenCastle = false
		}
	}

	state.CastleRights = newCastleRights

	updateBitboards(state)
	//updateHash(state, move, start, end, val, captVal, )
}


func UnMakeLastMove(state *board.GameState) {
	state.RestorePrev()
}


func CreateGameState(b [64]int, whiteMove bool, wkCastle bool, wqCastle bool, bkCastle bool, bqCastle bool, pDouble [2]int) board.GameState {
	//to be called whenever new game state obj is created

	var whitePiecePos [6][10][2]int
	var blackPiecePos [6][10][2]int
	for i := 0; i < 6; i++ {
		for x := 0; x < 10; x++ {
			for y := 0; y < 2; y++ {
				//default values
				whitePiecePos[i][x][y] = -1
				blackPiecePos[i][x][y] = -1
			}
		}
	}

	var inxs [12]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece := b[x * 8 + y]

			if piece != 0 {
				inx := inxs[piece - 1]
				pos := [2]int{x, y}

				if piece < 7 {
					whitePiecePos[piece - 1][inx] = pos
				} else {
					blackPiecePos[piece - 7][inx] = pos
				}

				inxs[piece - 1]++
			}
		}
	}

	var castleRights uint8
	if wkCastle {castleRights |= board.WkCastle}
	if wqCastle {castleRights |= board.WqCastle}
	if bkCastle {castleRights |= board.BkCastle}
	if bkCastle {castleRights |= board.BqCastle}

	state := board.GameState{Board: b, WhiteToMove: whiteMove, CastleRights: castleRights, PrevPawnDouble: pDouble, WhitePiecePos: whitePiecePos, BlackPiecePos: blackPiecePos}

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
	kAttackBlock, PinArray, noKingMove, EnPassantPin := GetFilterBitboards(&state.Board, kingX, kingY, kingVal, otherPieces, whiteMove, pDouble)

	state.NoKingMoveBitBoard = noKingMove
	state.KingAttackBlocks = kAttackBlock
	state.PinArray = PinArray
	state.EnPassantPin = EnPassantPin

	zobHash := board.HashState(&state)
	state.ZobristHash = zobHash

	return state
}