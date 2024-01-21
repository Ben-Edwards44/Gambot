package moves


func enPassant(state GameState, currentX int, currentY int, pieceValue int) Move {
	otherX := state.PrevPawnDouble[0]
	otherY := state.PrevPawnDouble[1]

	//if opponent double moved pawn and pawns are on same rank
	if otherX == currentX {
		fileDiff := currentY - otherY

		currentWhite := pieceValue < 7
		otherWhite := state.Board[otherX * 8 + otherY] < 7

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


func castle(state GameState, pieceValue int, resultSlice *[]Move) {
	//black values
	rookValue := 10
	kingValue := 11
	x := 0
	canCastleKing := state.BlackKingCastle
	canCastleQueen := state.BlackQueenCastle
	if pieceValue < 7 {
		//white values
		rookValue = 4
		kingValue = 5
		x = 7
		canCastleKing = state.WhiteKingCastle
		canCastleQueen = state.WhiteQueenCastle
	}

	rankInx := x * 8
	kingPos := rankInx + 4
	kRookPos := rankInx + 7
	qRookPos := rankInx
	
	if canCastleKing {
		//the king/rook have not been explicitly moved (could have been captured though)
		if state.Board[kingPos] == kingValue && state.Board[kRookPos] == rookValue {
			var badBitBoard uint64
			var pieceInWay bool
			for i := kingPos; i <= kRookPos - 1; i++ {
				setBitBoard(&badBitBoard, i)

				//we are looping over the king square to ensure the king is not in check
				if state.Board[i] != 0 && state.Board[i] != kingValue {
					pieceInWay = true
					break
				}
			}

			//bitwise AND the bitboards to ensure no crossover
			if !pieceInWay && (badBitBoard & state.otherMoveBitBoard == 0) {
				//not going into check
				m := Move{StartX: x, StartY: 4, EndX: x, EndY: 6, PieceValue: kingPos, KingCastle: true}
				*resultSlice = append(*resultSlice, m)
			}
		}
	}
	if canCastleQueen {
		if state.Board[kingPos] == kingValue && state.Board[qRookPos] == rookValue {
			var badBitBoard uint64
			var pieceInWay bool
			for i := qRookPos + 2; i <= kingPos; i++ {
				setBitBoard(&badBitBoard, i)

				//we are looping over the king square to ensure the king is not in check
				if state.Board[i] != 0 && state.Board[i] != kingValue {
					pieceInWay = true
					break
				}
			}

			//extra check for queenside only (other positions will have been checked in above loop)
			if state.Board[qRookPos + 1] != 0 {
				pieceInWay = true
			}

			//bitwise AND the bitboards to ensure no crossover
			if !pieceInWay && (badBitBoard & state.otherMoveBitBoard == 0) {
				m := Move{StartX: x, StartY: 4, EndX: x, EndY: 2, PieceValue: kingPos, QueenCastle: true}
				*resultSlice = append(*resultSlice, m)
			}
		}
	}
}