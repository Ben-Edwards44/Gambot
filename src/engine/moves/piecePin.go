package moves


//TODO: do this
//shoot rays out from king and check for enemy piece / skip over 1 friendly
//need to do extra checks for knight/pawn attacks (but not pins)


func straightRay(board [64]int, x int, y int, isWhite bool, attackPiece *[][2]int, pinPiece *[][2]int) {
	for dist := 0; dist < 4; dist++ {
		passedPieces := 0
		passedX := -1
		passedY := -1

		//arrays initialised in moves/pieceMoves
		edgeDist := dists[dist]
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
						passedX = newX
						passedY = newY
					} else {
						break
					}
				} else {
					//enemy piece

					//TODO: add passed to pin or add enemy to attack depending on num passed
					break
				}
			}
		}
	}
}