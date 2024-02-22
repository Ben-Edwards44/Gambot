package search


import (
	"fmt"
	"time"
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/evaluation"
)


const INF int = 9999999
const MATESCORE int = 100000

var searchAbandoned bool

var bestMoves map[[64]int]moves.Move

var posSearched int


func checkWin(state *board.GameState, isWhite bool) int {
	//TODO: include depths in checkmates (for transposition table)
	//check who has won, or if it is a draw - assumes that the player has no legal moves

	kingPos := state.BlackPiecePos[4][0]
	if isWhite {kingPos = state.WhitePiecePos[4][0]}

	pos := kingPos[0] * 8 + kingPos[1]

	//set bitboard at king's position
	var kingPosBB uint64
	kingPosBB |= 1 << pos

	inCheck := (kingPosBB & state.NoKingMoveBitBoard) != 0

	if inCheck {
		//we are checkmated :(
		return -MATESCORE  //negative because being checkmated is bad
	} else {
		return 0  //draw
	}
}


func negamax(state *board.GameState, isWhite bool, depth int, alpha int, beta int, timeLeft time.Duration) (int, moves.Move) {
	if timeLeft < 0 {
		//out of time
		searchAbandoned = true
		return 0, moves.Move{}
	}

	startTime := time.Now()

	//ttSuccess, ttEval := evaluation.LookupEval(state.ZobristHash, depth, alpha, beta)

	//if ttSuccess {
	//	//this position is in transposition table. We don't need to search it again
	//	bestMove := evaluation.LookupMove(state.ZobristHash)  //TODO: only do this when searching root node
	//	return ttEval, bestMove
	//}
	
	if depth == 0 {return quiescenceSearch(state, isWhite, alpha, beta, timeLeft), moves.Move{}}

	posSearched++

	moveList := moves.GenerateAllMoves(state, false)
	orderMoves(state, moveList, bestMoves[state.Board])

	//TODO: draws by repetition and 50 move rule etc.
	if len(moveList) == 0 {return checkWin(state, isWhite), moves.Move{}}  //deal with checkmates and draws

	bestScore := -INF
	allocatedBestMove := false

	var bestMove moves.Move

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		negScore, _ := negamax(state, !isWhite, depth - 1, -beta, -alpha, timeLeft - elapsed)
		score := -negScore

		moves.UnMakeLastMove(state)

		if score > bestScore || !allocatedBestMove {
			bestMove = move
			bestScore = score
			allocatedBestMove = true
		}

		//fail-hard cutoff (prune position)
		if score >= beta {
			//evaluation.StoreEntry(state.ZobristHash, depth, beta, )

			return beta, bestMove
		}  

		if score > alpha {alpha = score}
	}

	//update the best moves (because we will be searching at a greater depth)
	bestMoves[state.Board] = bestMove

	return bestScore, bestMove
}


func quiescenceSearch(state *board.GameState, isWhite bool, alpha int, beta int, timeLeft time.Duration) int {
	//TODO: include checks here
	//IMPORTANT: this will not work with checkmates because eval does not return the checkmate values
	if timeLeft < 0 {
		searchAbandoned = true
		return 0
	}
	
	startTime := time.Now()

	staticEval := evaluation.Eval(state, isWhite)

	if staticEval >= beta {return beta}  //prune
	if staticEval > alpha {alpha = staticEval}

	posSearched++
	
	moveList := moves.GenerateAllMoves(state, true)
	orderMoves(state, moveList, moves.Move{})//bestMoves[state.Board])

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		score := -quiescenceSearch(state, !isWhite, -beta, -alpha, timeLeft - elapsed)
	
		moves.UnMakeLastMove(state)

		if score >= beta {return beta}  //prune
		if score > alpha {alpha = score}
	}

	return alpha
}


func GetBestMove(state *board.GameState) moves.Move {
	searchAbandoned = false
	bestMoves = make(map[[64]int]moves.Move)

	startTime := time.Now()
	timeLeft := time.Duration(time.Millisecond * 500)  //NOTE: change to 500ms for testing

	depth := 1

	searchedDepthOne := false

	var bestMove moves.Move
	var elapsed time.Duration
	for timeLeft > 0 {
		posSearched = 0
		elapsed = time.Since(startTime)

		timeLeft -= elapsed

		score, searchBestMove := negamax(state, state.WhiteToMove, depth, -INF, INF, timeLeft)  //NOTE: don't need to -score because this call is from the POV of the engine

		fmt.Print("Depth: ")
		fmt.Print(depth)
		fmt.Print(", Searched: ")
		fmt.Print(posSearched)
		fmt.Print(", Elapsed: ")
		fmt.Println(elapsed)  //time taken must be last for speed test

		if !searchAbandoned {
			bestMove = searchBestMove
			searchedDepthOne = true
		} else {
			break
		}

		if score == MATESCORE {break}  //we have found a mate, so don't search any deeper

		depth++
	}

	if !searchedDepthOne {
		//could not even search to depth one, just play the first available move

		moveList := moves.GenerateAllMoves(state, false)
		orderMoves(state, moveList, moves.Move{})  //order so that a reasonable looking move is played

		bestMove = moveList[0]
	}

	if (bestMove == moves.Move{}) {
		panic("Best move is an empty move")
	}

	return bestMove
}