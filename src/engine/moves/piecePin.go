package moves


var attackFunctions [6]func([64]int, int, int, int, int, *uint64, *[]uint64, *[64]uint64) = [6]func([64]int, int, int, int, int, *uint64, *[]uint64, *[64]uint64){pawnAttacks, knightAttacks, bishopAttacks, rookAttacks, kingAttacks, queenAttacks}


func rookAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		returnAfterStraight := false
		passedPiece := false
		passedInx := -1
		
		edgeDist := dists[dirInx + dir]  //get from precalculated

		var currentStraight uint64
		setBitBoard(&currentStraight, x * 8 + y)  //because we can always take the rook if it is attacking

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir]
			newY := y + offset * yMults[dir]
			pos := newX * 8 + newY

			goodSq, capture := canMove(board, newX, newY, pieceValue)

			if !passedPiece {setBitBoard(noKingBB, pos)}

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
						pinBB[passedInx] = currentStraight
						return
					} else {
						//this rook is directly attacking king (cannot return because we need to finish updating noKingArray)
						*attackBB = append(*attackBB, currentStraight)

						returnAfterStraight = true
					}

					setBitBoard(noKingBB, pos)
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


func bishopAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	dirInx := x * 64 + y * 8

	for dir := 0; dir < 4; dir++ {
		returnAfterStraight := false
		passedPiece := false
		passedInx := -1

		edgeDist := dists[dirInx + dir + 4]  //get from precalculated

		var currentStraight uint64
		setBitBoard(&currentStraight, x * 8 + y)  //because we can always take the rook if it is attacking

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir + 4]
			newY := y + offset * yMults[dir + 4]
			pos := newX * 8 + newY

			goodSq, capture := canMove(board, newX, newY, pieceValue)

			if !passedPiece {setBitBoard(noKingBB, pos)}

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
						pinBB[passedInx] = currentStraight
						return
					} else {
						//this rook is directly attacking king (cannot return because we need to finish updating noKingArray)
						*attackBB = append(*attackBB, currentStraight)

						returnAfterStraight = true
					}

					setBitBoard(noKingBB, pos)
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


func queenAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	rookAttacks(board, x, y, pieceValue, kingValue, noKingBB, attackBB, pinBB)
	bishopAttacks(board, x, y, pieceValue, kingValue, noKingBB, attackBB, pinBB)
}


func kingAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	edgeInx := x * 64 + y * 8

	for dir := 0; dir < 8; dir++ {
		edgeDist := dists[edgeInx + dir]

		if edgeDist > 0 {
			newX := x + xMults[dir]
			newY := y + yMults[dir]
			pos := newX * 8 + newY

			setBitBoard(noKingBB, pos)
		}
	}
}


func knightAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	var posBB uint64
	setBitBoard(&posBB, x * 8 + y)
	
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

				setBitBoard(noKingBB, pos)
				if board[pos] == kingValue {*attackBB = append(*attackBB, posBB)}
			}
		}
	}
}


func pawnAttacks(board [64]int, x int, y int, pieceValue int, kingValue int, noKingBB *uint64, attackBB *[]uint64, pinBB *[64]uint64) {
	//TODO: do this

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

			setBitBoard(noKingBB, pos)

			if capture && board[pos] == kingValue {
				var posBB uint64
				setBitBoard(&posBB, x * 8 + y)

				*attackBB = append(*attackBB, posBB)
			}
		}
	}
}


func enPassantPin(board [64]int, kingX int, kingY int, isWhite bool, prevPawnDouble [2]int) bool {
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
	edgeDist := dists[distInx + 1]
	if prevPawnDouble[1] < kingY {
		yStep = -1
		edgeDist = dists[distInx]
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


func getFilterBitboards(board [64]int, kingX int, kingY int, kingValue int, otherPiecePos [][2]int, isWhite bool, prevPawnDouble [2]int) ([]uint64, [64]uint64, uint64, bool) {
	//return the attack, pin, and no king move bitboards, as well as whether en passant is pinned

	var attackBB []uint64
	var pinBB [64]uint64
	var noKingBB uint64

	for _, i := range otherPiecePos {
		pieceValue := board[i[0] * 8 + i[1]]

		inx := pieceValue - 1
		if isWhite {inx = pieceValue - 7}  //counter-intuitive way around because we are looking at the enemy moves

		atkFunc := attackFunctions[inx]
		atkFunc(board, i[0], i[1], pieceValue, kingValue, &noKingBB, &attackBB, &pinBB)
	}

	enPassPin := enPassantPin(board, kingX, kingY, isWhite, prevPawnDouble)

	return attackBB, pinBB, noKingBB, enPassPin
}