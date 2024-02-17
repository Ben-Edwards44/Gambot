package search


import (
	"time"
	"chess-engine/src/engine/moves"
)


const INF int = 100000

var searchAbandoned bool


func getMoveOrder(state *moves.GameState, move moves.Move) int {
	score := 0

	currentPiece := move.PieceValue
	promotion := move.PromotionValue
	captPiece := state.Board[move.EndX * 8 + move.EndY] - 6  //-6 because it is opposite colour to current
	if currentPiece > 6 {
		currentPiece -= 6
		promotion -= 6
		captPiece += 6
	}

	if captPiece > 0 {score += captPiece - currentPiece}  //capturing high value pieces with low value ones is good

	score += promotion  //promotions are good (if not promotion, this will just add 0 to score)

	var posBB uint64
	posBB |= 1 << (move.StartX * 8 + move.StartY)

	if posBB & state.NoKingMoveBitBoard != 0 {score -= currentPiece}  //moving to an attacked square is not good

	return score
}


func orderMoves(state *moves.GameState, moveList []moves.Move) {
	//slices are passed by reference, so no need to return

	var moveScores []int
	for _, i := range moveList {
		moveScores = append(moveScores, getMoveOrder(state, i))
	}

	quickSort(moveList, moveScores, 0, len(moveList) - 1)
}


func checkWin(state *moves.GameState, isWhite bool, depth int) int {
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
			return -INF - depth
		} else {
			//black is checkmated
			return INF + depth
		}
	} else {
		//draw
		return 0
	}
}


func negamax(state *moves.GameState, isWhite bool, depth int, alpha int, beta int, timeLeft time.Duration) (int, moves.Move) {
	if timeLeft < 0 {
		//out of time
		searchAbandoned = true
		return 0, moves.Move{}
	}

	startTime := time.Now()
	
	if depth == 0 {return quiescenceSearch(state, isWhite, alpha, beta), moves.Move{}}

	moveList := moves.GenerateAllMoves(state, false)
	orderMoves(state, moveList)

	bestScore := -INF
	allocatedBestMove := false
	var bestMove moves.Move

	for _, i := range moveList {
		moves.MakeMove(state, i)

		elapsed := time.Since(startTime)
		negScore, _ := negamax(state, !isWhite, depth - 1, -beta, -alpha, timeLeft - elapsed)
		score := -negScore

		moves.UnMakeLastMove(state)

		if score > bestScore || !allocatedBestMove {
			bestMove = i
			bestScore = score
			allocatedBestMove = true
		}

		if score >= beta {return beta, moves.Move{}}

		if score > alpha {alpha = score}
	}

	return bestScore, bestMove
}



func quiescenceSearch(state *moves.GameState, isWhite bool, alpha int, beta int) int {
	//this does not work :(
	staticEval := eval(state, isWhite)

	return staticEval

	//IMPORTANT: this will not work with checkmates becauseeval does ont return the checkmate values

	if isWhite {
		//max player
		if staticEval > alpha {
			alpha = staticEval
		}
	} else {
		//min player
		if staticEval < beta {
			beta = staticEval
		}
	}

	if staticEval >= beta {
		return beta
	} //prune

	captMoves := moves.GenerateAllMoves(state, true)

	if len(captMoves) == 0 {return staticEval}  //quiet position reached - return the evalutation

	bestScore := INF + 1  //+1 so that engine will still move if in forced mate
	if isWhite {bestScore = -INF - 1}  //-1 so that engine will still move if in forced mate

	for _, i := range captMoves {
		moves.MakeMove(state, i)
		score := quiescenceSearch(state, !isWhite, alpha, beta)
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
	searchAbandoned = false

	startTime := time.Now()
	timeLeft := time.Duration(time.Millisecond * 1000)

	depth := 1
	var bestMove moves.Move
	var elapsed time.Duration
	for timeLeft > 0 {
		elapsed = time.Since(startTime)
		timeLeft -= elapsed

		_, searchBestMove := negamax(state, state.WhiteToMove, depth, -INF, INF, timeLeft)
		
		if !searchAbandoned {
			bestMove = searchBestMove
		} else {
			break
		}

		depth++
	}

	return bestMove
}