package search


import (
	"fmt"
	"time"
	"chess-engine/src/engine/moves"
)


const INF int = 100000


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


func minimax(state *moves.GameState, isWhite bool, depth int, alpha int, beta int) (int, moves.Move) {
	//NOTE: white is the max player
	if depth == 0 {return quiescenceSearch(state, isWhite, alpha, beta),  moves.Move{}}  //TODO: do this (figure out whether we should be white/black)

	moveList := moves.GenerateAllMoves(state, false)
	orderMoves(state, moveList)

	if len(moveList) == 0 {return checkWin(state, isWhite, depth), moves.Move{}}  //terminal node

	allocatedBestMove := false

	var bestScore int
	var bestMove moves.Move
	for _, i := range moveList {
		moves.MakeMove(state, i)
		score, _ := minimax(state, !isWhite, depth - 1, alpha, beta)
		moves.UnMakeLastMove(state)

		if isWhite {
			//max player
			if score > bestScore || !allocatedBestMove {
				bestScore = score
				bestMove = i
				allocatedBestMove = true
			}

			if score > alpha {alpha = score}
		} else {
			//min player
			if score < bestScore || !allocatedBestMove {
				bestScore = score
				bestMove = i
				allocatedBestMove = true
			}

			if score < beta {beta = score}
		}

		if beta <= alpha {break}  //prune position (opponent already has a better position)
	}

	return bestScore, bestMove
}


func quiescenceSearch(state *moves.GameState, isWhite bool, alpha int, beta int) int {
	//this does not work :(
	staticEval := eval(state)

	return staticEval

	if isWhite {
		//max player
		if staticEval > alpha {alpha = staticEval}
	} else {
		//min player
		if staticEval < beta {beta = staticEval}
	}

	if staticEval >= beta {return beta}  //prune

	captMoves := moves.GenerateAllMoves(state, true)

	if len(captMoves) == 0 {return staticEval}  //quiet position reached - return the evalutation

	bestScore := INF + 1  //+1 so that engine will still move if in forced mate
	if isWhite {
		bestScore = -INF - 1  //-1 so that engine will still move if in forced mate
	} 

	for _, i := range captMoves {
		moves.MakeMove(state, i)
		score := quiescenceSearch(state, !isWhite, -alpha, -beta)
		moves.UnMakeLastMove(state)

		if isWhite {
			//max player
			if score > bestScore {
				bestScore = score
			}

			if score > alpha {alpha = score}
		} else {
			//min player
			if score < bestScore {
				bestScore = score
			}

			if score < beta {beta = score}
		}

		if beta <= alpha {break}  //prune position (opponent already has a better position)
	}

	return bestScore
}


func GetBestMove(state *moves.GameState) moves.Move {
	start := time.Now()

	maxDepth := 5  //total moves from current position (so depth=1 means just look at our moves not opponent responses)

	_, bestMove := minimax(state, state.WhiteToMove, maxDepth, -INF, INF)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Time elapsed: ")
	fmt.Println(elapsed)

	return bestMove
}