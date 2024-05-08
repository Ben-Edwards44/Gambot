package moves


import "gambot/src/engine/board"


var numChecks int
var attackFunctions [6]func(*[64]int, int, int, int, int, *board.Bitboard) = [6]func(*[64]int, int, int, int, int, *board.Bitboard){pawnAttacks, knightAttacks, bishopAttacks, rookAttacks, kingAttacks, queenAttacks}


func setBitBoard(bb *uint64, pos int) {
	//set the inx of a bitboard to a 1
	bitWeight := uint64(1 << uint64(pos))  //2**n
	*bb |= bitWeight
}


func rookAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		returnAfterStraight := false
		passedPiece := false
		passedInx := -1
		
		edgeDist := edgeDists[dirInx + dir]  //get from precalculated

		var currentStraight uint64
		setBitBoard(&currentStraight, x * 8 + y)  //because we can always take the rook if it is attacking

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir]
			newY := y + offset * yMults[dir]
			pos := newX * 8 + newY

			goodSq, capture := canMove(board, newX, newY, pieceValue)

			if !passedPiece {setBitBoard(&bitboard.AttackedSquares, pos)}

			if !goodSq {
				break  //met friendly piece
			} else if !capture {
				//blank square
				setBitBoard(&currentStraight, pos)
			} else {
				//capture

				if board[pos] == kingValue {
					if passedPiece {
						//the passed piece is pinned
						bitboard.PinArray[passedInx] = currentStraight
						return
					} else {
						//this rook is directly attacking king (cannot return because we need to finish updating noKingArray)
						bitboard.AttacksOnKing |= currentStraight
						numChecks++

						returnAfterStraight = offset > 1  //we do not return after straight if a rook/queen is right next to king because other king moves must be blocked
					}

					setBitBoard(&bitboard.AttackedSquares, pos)				
				} else {
					if passedPiece {
						break
					} else {
						passedPiece = true
						passedInx = pos
					}
				}
			}
		}

		if returnAfterStraight {return}
	}
}


func bishopAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		returnAfterStraight := false
		passedPiece := false
		passedInx := -1

		edgeDist := edgeDists[dirInx + dir + 4]  //get from precalculated

		var currentStraight uint64
		setBitBoard(&currentStraight, x * 8 + y)  //because we can always take the rook if it is attacking

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir + 4]
			newY := y + offset * yMults[dir + 4]
			pos := newX * 8 + newY

			goodSq, capture := canMove(board, newX, newY, pieceValue)

			if !passedPiece {setBitBoard(&bitboard.AttackedSquares, pos)}

			if !goodSq {
				break  //met friendly piece
			} else if !capture {
				//blank square
				setBitBoard(&currentStraight, pos)
			} else {
				//capture

				if board[pos] == kingValue {
					if passedPiece {
						//the passed piece is pinned
						bitboard.PinArray[passedInx] = currentStraight
						return
					} else {
						//this rook is directly attacking king (cannot return because we need to finish updating noKingArray)
						bitboard.AttacksOnKing |= currentStraight
						numChecks++

						returnAfterStraight = true
					}

					setBitBoard(&bitboard.AttackedSquares, pos)
				} else {
					if passedPiece {
						break
					} else {
						passedPiece = true
						passedInx = pos
					}
				}
			}
		}

		if returnAfterStraight {return}
	}
}


func queenAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {
	rookAttacks(board, x, y, pieceValue, kingValue, bitboard)
	bishopAttacks(board, x, y, pieceValue, kingValue, bitboard)
}


func kingAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {
	edgeInx := x * 64 + y * 8

	for dir := 0; dir < 8; dir++ {
		edgeDist := edgeDists[edgeInx + dir]

		if edgeDist > 0 {
			newX := x + xMults[dir]
			newY := y + yMults[dir]
			pos := newX * 8 + newY

			setBitBoard(&bitboard.AttackedSquares, pos)
		}
	}
}


func knightAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {	
	for xStep := 1; xStep < 3; xStep++ {
		for xMult := -1; xMult < 2; xMult += 2 {
			newX := x + xStep * xMult

			if newX < 0 || newX > 7 {continue}

			//xStep 1 => yStep 2, xStep 2 => yStep 1
			yStep := 3 - xStep
			for yMult := -1; yMult < 2; yMult += 2 {
				newY := y + yStep * yMult
				pos := newX * 8 + newY

				if newY < 0 || newY > 7 {continue}

				setBitBoard(&bitboard.AttackedSquares, pos)

				if board[pos] == kingValue {
					posBB := uint64(1 << (x * 8 + y))

					bitboard.AttacksOnKing |= posBB
					numChecks++
				}
			}
		}
	}
}


func pawnAttacks(board *[64]int, x int, y int, pieceValue int, kingValue int, bitboard *board.Bitboard) {
	if x == 0 || x == 7 {return}  //on back rank so cannot put king in check

	xMult := 1
	if pieceValue < 7 {
		xMult = -1
	}
	
	//capture moves
	newX := x + xMult
	for i := -1; i < 2; i += 2 {
		newY := y + i

		if 0 <= newY && newY < 8 {
			pos := newX * 8 + newY
			_, capture := canMove(board, newX, newY, pieceValue)

			setBitBoard(&bitboard.AttackedSquares, pos)
			setBitBoard(&bitboard.PawnAttacks, pos)

			if capture && board[pos] == kingValue {
				//pawn is checking king
				posBB := uint64(1 << (x * 8 + y))

				bitboard.AttacksOnKing |= posBB
				numChecks++
			}
		}
	}
}


func enPassantPin(board *[64]int, kingX int, kingY int, isWhite bool, prevPawnDouble [2]int) bool {
	if prevPawnDouble[0] != kingX || prevPawnDouble[0] == -1 {return false}  //en passant not on right rank (also exits if no en passant)

	friendPawn := 7
	enemyPawn := 1
	enemyRook := 4
	enemyQueen := 6
	if isWhite {
		enemyRook = 10
		enemyQueen = 12
		friendPawn, enemyPawn = enemyPawn, friendPawn
	}

	var friendPawnAdj bool
	for yStep := -1; yStep < 2; yStep += 2 {
		y := prevPawnDouble[1] + yStep
		if 0 <= y && y < 8 && board[prevPawnDouble[0] * 8 + y] == friendPawn {friendPawnAdj = true}
	}

	if !friendPawnAdj {return false}  //no adjacent pawn to do en passant
	
	distInx := kingX * 64 + kingY * 8

	yStep := 1
	edgeDist := edgeDists[distInx + 3]
	if prevPawnDouble[1] < kingY {
		yStep = -1
		edgeDist = edgeDists[distInx + 2]
	}

	passedFriend := false
	passedEnemy := false
	for i := 1; i <= edgeDist; i++ {
		y := kingY + i * yStep

		val := board[kingX * 8 + y]
		if val == friendPawn {
			if passedFriend {return false}  //2 friendly pawns to block pin
			
			passedFriend = true
		} else if val == enemyPawn {
			if passedEnemy {return false}  //2 enemy pawns to block pin
			
			passedEnemy = true
		} else if val == enemyQueen || val == enemyRook {
			return passedEnemy && passedFriend  //if we have met both pawns and no other piece, en passant is pinned
		} else if val != 0 {
			return false  //some other piece to block pin
		}
	}

	return false
}


func resetCurrent(bitboard *board.Bitboard) {
	//reset the bitboards that will be overwritten (but not the other ones)
	bitboard.AttackedSquares = 0
	bitboard.AttacksOnKing = 0
	bitboard.PawnAttacks = 0
	bitboard.PinArray = [64]uint64{}
}


func GetFilterBitboards(board *[64]int, kingPos int, kingValue int, otherPiecePos *[16]int, isWhite bool, prevPawnDouble [2]int, bitboard *board.Bitboard) (bool, bool) {
	//update the bitboard struct with new bitboards and return whether double checked or en passant pinned
	resetCurrent(bitboard)
	
	numChecks = 0

	for i := 0; i < len(otherPiecePos); i++ {
		square := otherPiecePos[i]

		if square != -1 {
			pieceVal := board[square]

			atkInx := (pieceVal - 1) % 6  //white and black pieces use same attack function
			atkFunc := attackFunctions[atkInx]

			x := int(square / 8)
			y := square % 8

			atkFunc(board, x, y, pieceVal, kingValue, bitboard)
		}
	}

	kingX := int(kingPos / 8)
	kingY := kingPos % 8

	enPassPin := enPassantPin(board, kingX, kingY, isWhite, prevPawnDouble)

	doubleChecked := numChecks > 1

	return enPassPin, doubleChecked
}