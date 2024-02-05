package search


import "chess-engine/src/engine/moves"


const INF int = 100000


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
			return -INF
		} else {
			//black is checkmated
			return INF
		}
	} else {
		//draw
		return 0
	}
}


func minimax(state *moves.GameState, isWhite bool, depth int, alpha int, beta int) int {
	//NOTE: white is the max player
	if depth == 0 {return eval(state)}

	moveList := moves.GenerateAllMoves(state)

	if len(moveList) == 0 {return checkWin(state, isWhite)}  //terminal node

	bestScore := INF
	if isWhite {
		bestScore = -INF
	} 

	for _, i := range moveList {
		moves.MakeMove(state, i)
		score := minimax(state, !isWhite, depth - 1, alpha, beta)
		moves.UnMakeLastMove(state)

		if isWhite {
			//max player
			if score > bestScore {bestScore = score}
			if score > alpha {alpha = score}
		} else {
			//min player
			if score < bestScore {bestScore = score}
			if score < beta {beta = score}
		}

		if beta <= alpha {break}  //prune position (opponent already has a better position)
	}

	return bestScore
}


func GetBestMove(state *moves.GameState) moves.Move {
	maxDepth := 3  //total moves from current position (so depth=1 means just look at our moves not opponent responses)

	possibleMoves := moves.GenerateAllMoves(state)

	chosen := false
	maxPlayer := state.WhiteToMove

	var bestMove moves.Move
	var bestScore int
	for _, i := range possibleMoves {
		moves.MakeMove(state, i)
		score := minimax(state, state.WhiteToMove, maxDepth - 1, -INF, INF)  //-1 because we are already searching at depth 1
		moves.UnMakeLastMove(state)

		if !chosen || ((maxPlayer && score > bestScore) || (!maxPlayer && score < bestScore)) {
			bestMove = i
			bestScore = score
			chosen = true
		}
	}

	return bestMove
}