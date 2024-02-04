package search


import "chess-engine/src/engine/moves"


func checkWin(state *moves.GameState, isWhite bool) int {
	//check who has won, or if it is a draw - assumes that the player has no legal moves

	kingPos := state.BlackPiecePos[4][0]
	if isWhite {kingPos = state.WhitePiecePos[4][0]}

	pos := kingPos[0] * 8 + kingPos[1]

	//set bitboard at king's position
	var kingPosBB uint64
	kingPosBB |= 1 << pos

	inCheck := (kingPosBB & state.NoKingMoveBitBoard) != 0

	if inCheck {
		if isWhite {
			//white is checkmated
			return -10000
		} else {
			//black is checkmated
			return 10000
		}
	} else {
		//draw
		return 0
	}
}


func minimax(state *moves.GameState, isWhite bool, depth int) (int, moves.Move) {
	//NOTE: white is the max player
	if depth == 0 {return eval(state), moves.Move{}}

	moveList := moves.GenerateAllMoves(state)

	if len(moveList) == 0 {return checkWin(state, isWhite), moves.Move{}}  //terminal node

	var scores []int
	for _, i := range moveList {
		moves.MakeMove(state, i)
		score, _ := minimax(state, !isWhite, depth - 1)
		moves.UnMakeLastMove(state)

		scores = append(scores, score)
	}

	bestScore := scores[0]
	bestMove := moveList[0]

	for i, score := range scores[1:] {
		if isWhite {
			//max player
			if score > bestScore {
				bestScore = score
				bestMove = moveList[i + 1]  //+1 because we exclude the first move from the loop
			}
		} else {
			//min player
			if score < bestScore {
				bestScore = score
				bestMove = moveList[i + 1]  //+1 because we exclude the first move from the loop
			}
		}
	}

	return bestScore, bestMove
}


func GetBestMove(state *moves.GameState) moves.Move {
	maxDepth := 3
	_, move := minimax(state, state.WhiteToMove, maxDepth)

	return move
}