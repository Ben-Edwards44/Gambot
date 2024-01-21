package moves


func getAtttackPin(board [64]int, x int, y int, isWhite bool) ([][2]int, [][2]int) {
	//shoot rays out from king and check pieces
	
	var attackPieces [][2]int
	var pinPieces [][2]int

	dirInx := x * 64 + y * 8

	for dist := 0; dist < 8; dist++ {
		passedPieces := 0
		passedPos := [2]int{-1, -1}

		//arrays initialised in moves/pieceMoves
		edgeDist := dists[dirInx + dist]
		xStep := xMults[dist]
		yStep := yMults[dist]

		for i := 1; i <= edgeDist; i++ {
			newX := x + xStep * i
			newY := y + yStep * i

			val := board[newX * 8 + newY]
			if val != 0 {
				if val < 7 == isWhite {
					//friendly piece

					if passedPieces == 0 {
						passedPieces += 1
						passedPos = [2]int{newX, newY}
					} else {
						break
					}
				} else {
					//enemy piece

					actVal := val
					if val > 7 {
						actVal -= 6
					}

					bishop := dist > 3 && actVal == 3  //on diagonal and bishop
					rook := dist < 4 && actVal == 4  //on straight and rook

					if actVal == 6 || bishop || rook {
						if passedPieces == 0 {
							attackPieces = append(attackPieces, [2]int{newX, newY})
						} else {
							pinPieces = append(pinPieces, passedPos)
						}
					}

					break
				}
			}
		}
	}

	return attackPieces, pinPieces
}


func checkKnightAttacks(board [64]int, x int, y int, isWhite bool) [][2]int {
	var attackPieces [][2]int

	knightValue := 2  //white knight
	if isWhite {
		knightValue = 8  //black knight (swapped because we are looking for enemy knight)
	}

	for xStep := 1; xStep < 3; xStep++ {
		for xMult := -1; xMult < 2; xMult += 2 {
			newX := x + xStep * xMult

			if newX < 0 || newX > 7 {continue}

			//xStep 1 => yStep 2, xStep 2 => yStep 1
			yStep := 3 - xStep
			for yMult := -1; yMult < 2; yMult += 2 {
				newY := y + yStep * yMult

				if newY < 0 || newY > 7 {continue}

				val := board[newX * 8 + newY]

				if val == knightValue {
					attackPieces = append(attackPieces, [2]int{newX, newY})
				}
			}
		}
	}

	return attackPieces
}


func checkPawnAttacks(board [64]int, x int, y int, isWhite bool) [][2]int {
	var attackPieces [][2]int
	
	newX := x + 1
	pawnValue := 1  //white pawn
	if isWhite {
		newX = x - 1
		pawnValue = 7  //black pawn
	}

	if newX < 0 || newX > 7 {return attackPieces}

	for yStep := -1; yStep < 2; yStep += 2 {
		newY := y + yStep

		if newY < 0 || newY > 7 {continue}

		val := board[newX * 8 + newY]
		if val == pawnValue {
			attackPieces = append(attackPieces, [2]int{newX, newY})
		}
	}

	return attackPieces
}


func getStep(start int, end int) int {
	step := 0
	if end > start {
		step = 1
	} else if end < start {
		step = -1
	}

	return step
}


func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}


func fillBitboard(sX int, sY int, eX int, eY int) uint64 {
	//returns a bitboard with bits on line (sX, sY) to (eX, eY) set (start / end points inclusive)

	xStep := getStep(sX, sY)
	yStep := getStep(sY, eY)

	//get manhattan distance
	dist := abs(sX - eX)
	if dist == 0 {
		dist = abs(sY - eY)
	}

	var bitBoard uint64
	for i := 0; i <= dist; i++ {
		x := sX + xStep * i
		y := sY + yStep * i
		
		setBitBoard(&bitBoard, x * 8 + y)
	}

	return bitBoard
}


func addCaptureMoves(pos [][2]int, bitboards *[]uint64) {
	for _, i := range pos {
		var bitboard uint64
		setBitBoard(&bitboard, i[0] * 8 + i[1])

		*bitboards = append(*bitboards, bitboard)
	}
}


func legalFilterBitboards(board [64]int, kingX int, kingY int, isWhite bool) []uint64 {
	//TODO: pins
	var attackBB []uint64
	//var pinBB []uint64

	attackPos, _ := getAtttackPin(board, kingX, kingY, isWhite)

	for _, i := range attackPos {
		bitboard := fillBitboard(kingX, kingY, i[0], i[1])
		attackBB = append(attackBB, bitboard)
	}

	knight := checkKnightAttacks(board, kingX, kingY, isWhite)
	pawn := checkPawnAttacks(board, kingX, kingY, isWhite)

	addCaptureMoves(knight, &attackBB)
	addCaptureMoves(pawn, &attackBB)

	return attackBB
}