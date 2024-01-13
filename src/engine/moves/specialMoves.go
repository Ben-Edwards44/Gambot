package moves


func enPassant(board [64]int, currentX int, currentY int, pieceValue int, prevMove Move) (bool, Move) {
	//if opponent moved pawn 2 and pawns are on same rank
	if prevMove.PawnDoubleMove && prevMove.EndX == currentX {	
		rowDiff := currentY - prevMove.EndY

		if rowDiff == -1 || rowDiff == 1 {
			//account for colour
			newX := currentX + 1
			if pieceValue < 7 {
				newX = currentX - 1
			}

			newY := prevMove.EndY

			good, capture := canMove(board, newX, newY, pieceValue)
			if good && !capture {
				m := Move{StartX: currentX, StartY: currentY, EndX: newX, EndY: newY, PieceValue: pieceValue}
				return true, m
			}
		}
	}

	return false, Move{}
}