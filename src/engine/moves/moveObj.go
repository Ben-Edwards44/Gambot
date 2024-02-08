package moves


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//special moves
	doublePawnMove bool
	EnPassant bool
	KingCastle bool
	QueenCastle bool
	PromotionValue int
}


func removeFromArray(arr [8][2]int, x int, y int) [8][2]int {
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
			arr[i - 1][0] = -1
			arr[i - 1][1] = -1
		}
	}

	return arr
}


func appendToArray(arr [8][2]int, x int, y int) [8][2]int {
	for i, j := range arr {
		if j[0] == -1 {
			arr[i][0] = x
			arr[i][1] = y
			break
		}
	}

	return arr
}


func movePos(arr [8][2]int, sX int, sY int, eX int, eY int) [8][2]int {
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


func updateCapture(state *GameState, move Move, ePos int, isWhite bool) {
	//TODO: work with fixed length array
	eVal := state.Board[ePos]
	var enemy [6][8][2]int
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


func updatePiecePos(move Move, sPos int, ePos int, sVal int, state *GameState) {
	isWhite := sVal < 7

	if move.PromotionValue == 0 {
		//not a promotion
		var friend [6][8][2]int
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


func updateBitboards(state *GameState) {
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
	kAttackBlock, pinArray, noKingMove, enPassantPin := getFilterBitboards(&state.Board, kingX, kingY, kingVal, otherPieces, state.WhiteToMove, state.PrevPawnDouble)

	state.NoKingMoveBitBoard = noKingMove
	state.kingAttackBlocks = kAttackBlock
	state.pinArray = pinArray
	state.enPassantPin = enPassantPin
}


func MakeMove(state *GameState, move Move) {
	//updates game state

	state.SetPrevVals()  //so that we can restore later

	start := move.StartX * 8 + move.StartY
	end := move.EndX * 8 + move.EndY
	val := move.PieceValue

	updatePiecePos(move, start, end, val, state)

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
	
	if move.doublePawnMove {
		state.PrevPawnDouble = [2]int{move.EndX, move.EndY}
	} else {
		state.PrevPawnDouble = [2]int{-1, -1}
	}

	if move.PromotionValue != 0 {
		state.Board[end] = move.PromotionValue
	}

	if move.PieceValue == 5 {
		//white king moving
		state.WhiteKingCastle = false
		state.WhiteQueenCastle = false
	} else if move.PieceValue == 11 {
		//black king moving
		state.BlackKingCastle = false
		state.BlackQueenCastle = false
	} else if move.PieceValue == 4 {
		//white rook moving
		if move.StartY == 7 {
			state.WhiteKingCastle = false
		} else if move.StartY == 0 {
			state.WhiteQueenCastle = false
		}
	} else if move.PieceValue == 10 {
		//black rook moving
		if move.StartY == 7 {
			state.BlackKingCastle = false
		} else if move.StartY == 0 {
			state.BlackQueenCastle = false
		}
	}

	updateBitboards(state)
}


func UnMakeLastMove(state *GameState) {
	state.RestorePrev()
}