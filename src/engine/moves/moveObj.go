package moves


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//special moves
	EnPassant bool
	KingCastle bool
	QueenCastle bool
	promotionValue int
}


func updatePiecePos(state *GameState) {
	//TODO: optimise (use map probably) becuase this is really slow

	var whitePiecePos [][2]int 
	var blackPiecePos [][2]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece := state.Board[x * 8 + y]

			if piece != 0 {
				if piece < 7 {
					whitePiecePos = append(whitePiecePos, [2]int{x, y})
				} else {
					blackPiecePos = append(blackPiecePos, [2]int{x, y})
				}
			}
		}
	}

	state.WhitePiecePos = whitePiecePos
	state.BlackPiecePos = blackPiecePos
}


func MakeMoveCopy(state GameState, move Move) GameState {
	//returns a new copy of a game state

	start := move.StartX * 8 + move.StartY
	end := move.EndX * 8 + move.EndY
	val := move.PieceValue

	state.Board[start] = 0
	state.Board[end] = val

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
	} else if move.PieceValue == 9 {
		//black rook moving
		if move.StartY == 7 {
			state.BlackKingCastle = false
		} else if move.StartY == 0 {
			state.BlackQueenCastle = false
		}
	}

	state.WhiteToMove = !state.WhiteToMove  //because we have just made a move
	updatePiecePos(&state)

	return state
}