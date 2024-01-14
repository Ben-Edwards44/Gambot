package moves


func enPassant(board [64]int, currentX int, currentY int, pieceValue int, prevPawnDouble [2]int) Move {
	otherX := prevPawnDouble[0]
	otherY := prevPawnDouble[1]

	//if opponent double moved pawn and pawns are on same rank
	if otherX == currentX {
		fileDiff := currentY - otherY

		currentWhite := pieceValue < 7
		otherWhite := board[otherX * 8 + otherY] < 7

		//if pawns are next to each other and different colours
		if (fileDiff == 1 || fileDiff == -1) && currentWhite != otherWhite {
			newX := currentX + 1
			if currentWhite {
				newX = currentX - 1
			}

			newY := otherY

			m := Move{StartX: currentX, StartY: currentY, EndX: newX, EndY: newY, PieceValue: pieceValue, EnPassant: true}

			return m
		}
	}

	return Move{}
}