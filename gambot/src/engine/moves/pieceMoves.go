package moves


import "gambot/src/engine/board"


//array matching piece values to their appropriate move functions
var moveFunctions [6]func(*board.GameState, int, int, int, *[]*Move, bool) = [6]func(*board.GameState, int, int, int, *[]*Move, bool) {pawnMoves, knightMoves, bishopMoves, rookMoves, kingMoves, queenMoves}

//arrays matching distance index to their x, y multipliers
var xMults [8]int = [8]int{-1, 1, 0, 0, -1, -1, 1, 1}
var yMults [8]int = [8]int{0, 0, -1, 1, -1, 1, -1, 1}


func canMove(board *[64]int, x int, y int, pieceValue int) (bool, bool) {
	inx := x * 8 + y
	sqValue := board[inx]

	if sqValue == 0 {
		return true, false
	} else {
		return (sqValue > 6) != (pieceValue > 6), true
	}
}


func blockKingAttack(x int, y int, attacksOnKing uint64) bool {
	if attacksOnKing == 0 {return true}  //king not in check

	pos := x * 8 + y
	posBB := uint64(1 << pos)

	return attacksOnKing & posBB != 0
}


func checkPin(sX int, sY int, eX int, eY int, PinArray *[64]uint64) bool {
	bitboard := PinArray[sX * 8 + sY]

	if bitboard == 0 {return true}  //piece not pinned

	var posBB uint64
	setBitBoard(&posBB, eX * 8 + eY)

	good := bitboard & posBB

	return good != 0
}


func rookMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		//get from precalculated
		edgeDist := edgeDists[dirInx + dir]

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir]
			newY := y + offset * yMults[dir]

			goodSq, capture := canMove(&state.Board, newX, newY, pieceValue)
			blocking := blockKingAttack(newX, newY, state.Bitboards.AttacksOnKing)
			pin := checkPin(x, y, newX, newY, &state.Bitboards.PinArray)

			if goodSq && blocking && pin {
				m := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
				if !onlyCaptures || capture {*resultSlice = append(*resultSlice, &m)}
			}
			if capture {
				break
			}
		}
	}
}


func bishopMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		//get from precalculated (+4 since we are looking at diagonal)
		edgeDist := edgeDists[dirInx + dir + 4]

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir + 4]
			newY := y + offset * yMults[dir + 4]

			goodSq, capture := canMove(&state.Board, newX, newY, pieceValue)
			blocking := blockKingAttack(newX, newY, state.Bitboards.AttacksOnKing)
			pin := checkPin(x, y, newX, newY, &state.Bitboards.PinArray)

			if goodSq && blocking && pin {
				m := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
				if !onlyCaptures || capture {*resultSlice = append(*resultSlice, &m)}
			}
			if capture {
				break
			}
		}
	}
}


func queenMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	//resultSlice is updated within the functions
	rookMoves(state, x, y, pieceValue, resultSlice, onlyCaptures)
	bishopMoves(state, x, y, pieceValue, resultSlice, onlyCaptures)
}


func kingMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	edgeInx := x * 64 + y * 8

	for dir := 0; dir < 8; dir++ {
		edgeDist := edgeDists[edgeInx + dir]

		if edgeDist > 0 {
			newX := x + xMults[dir]
			newY := y + yMults[dir]
			
			good, capture := canMove(&state.Board, newX, newY, pieceValue)

			if good {
				//ensure not castling into check
				var moveBitBoard uint64 = 0
				setBitBoard(&moveBitBoard, newX * 8 + newY)

				//if not moving to an attacked square
				if moveBitBoard & state.Bitboards.AttackedSquares == 0 {
					m := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
					if !onlyCaptures || capture {*resultSlice = append(*resultSlice, &m)}
				}
			}
		}
	}
}


func knightMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	for xStep := 1; xStep < 3; xStep++ {
		for xMult := -1; xMult < 2; xMult += 2 {
			newX := x + xStep * xMult

			if newX < 0 || newX > 7 {continue}

			//xStep 1 => yStep 2, xStep 2 => yStep 1
			yStep := 3 - xStep
			for yMult := -1; yMult < 2; yMult += 2 {
				newY := y + yStep * yMult

				if newY < 0 || newY > 7 {continue}

				good, capture := canMove(&state.Board, newX, newY, pieceValue)
				blocking := blockKingAttack(newX, newY, state.Bitboards.AttacksOnKing)
				pin := checkPin(x, y, newX, newY, &state.Bitboards.PinArray)

				if good && blocking && pin {
					m := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
					if !onlyCaptures || capture {*resultSlice = append(*resultSlice, &m)}
				}
			}
		}
	}
}


func pawnMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	if x == 0 || x == 7 {return}  //on back rank (although this should never happen)

	isWhite := pieceValue < 7

	xMult := 1
	if isWhite {
		xMult = -1
	}

	if (isWhite && x == 1) || (!isWhite && x == 6) {
		//promotions
		promotion(state, x, y, pieceValue, xMult, resultSlice, onlyCaptures)
		return
	}

	if !onlyCaptures {
		onStart := (isWhite && x == 6) || (!isWhite && x == 1)

		maxStep := 1
		if onStart {
			maxStep = 2
		}
	
		//normal moves - no capture
		for i := 1; i <= maxStep; i++ {
			newX := x + i * xMult
	
			good, capture := canMove(&state.Board, newX, y, pieceValue)
			blocking := blockKingAttack(newX, y, state.Bitboards.AttacksOnKing)
			pin := checkPin(x, y, newX, y, &state.Bitboards.PinArray)
	
			if good && !capture && blocking && pin {
				m := Move{StartX: x, StartY: y, EndX: newX, EndY: y, PieceValue: pieceValue, DoublePawnMove: i == 2}
				*resultSlice = append(*resultSlice, &m)
			} else if !good || capture {
				break  //to prevent double pawn move when there is a piece in front
			}
		}
	}

	//capture moves
	newX := x + xMult
	for i := -1; i < 2; i += 2 {
		newY := y + i

		if 0 <= newY && newY < 8 {
			good, capture := canMove(&state.Board, newX, newY, pieceValue)
			blocking := blockKingAttack(newX, newY, state.Bitboards.AttacksOnKing)
			pin := checkPin(x, y, newX, newY, &state.Bitboards.PinArray)

			if good && capture && blocking && pin {
				m := Move{StartX: x, StartY: y, EndX: newX, EndY: newY, PieceValue: pieceValue}
				*resultSlice = append(*resultSlice, &m)
			}
		}
	}
}


func specialMoves(state *board.GameState, x int, y int, pieceValue int, resultSlice *[]*Move, onlyCaptures bool) {
	//a pointer is used for moves to ensure it is passed by reference
	
	if pieceValue == 1 || pieceValue == 7 {
		//pawn - check for en passant
		enPassant(state, x, y, pieceValue, resultSlice)
	} else if !onlyCaptures && (pieceValue == 5 || pieceValue == 11) {
		//king - check for castle
		castle(state, pieceValue, resultSlice)
	}
} 


func getPieceMoves(state *board.GameState, x int, y int, resultSlice *[]*Move, onlyCaptures bool) {
	pieceValue := state.Board[x * 8 + y]

	if state.DoubleChecked && (pieceValue != 5 && pieceValue != 11) {return}  //if double checked, only the king can be moved

	if pieceValue != 0 {
		var inx int

		//accounts for white/black
		if pieceValue < 7 {
			inx = pieceValue - 1
		} else {
			inx = pieceValue - 7
		}

		moveFunc := moveFunctions[inx]

		//update resultSlice
		moveFunc(state, x, y, pieceValue, resultSlice, onlyCaptures)

		//perform any special moves (en passant, castling etc.). These will be appended to the slice
		specialMoves(state, x, y, pieceValue, resultSlice, onlyCaptures)
	} else {
		panic("trying to find move for empty square")
	}
}


func GenerateAllMoves(state *board.GameState, onlyCaptures bool) []*Move {
	//assumes state has been properly initialised etc.

	piecePos := &board.PieceLists.BlackPieceSquares
	if state.WhiteToMove {
		piecePos = &board.PieceLists.WhitePieceSquares
	}

	moves := make([]*Move, 0, 64)  //need to experiment with how much memory to preallocate (max is 218, but this takes longer to allocate)
	for i := 0; i < len(piecePos); i++ {
		square := piecePos[i]

		if square != -1 {
			x := int(square / 8)
			y := square % 8

			getPieceMoves(state, x, y, &moves, onlyCaptures)
		}
	}

	return moves
}