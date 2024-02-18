package search


import (
	"chess-engine/src/engine/moves"
	"time"
)


const INF int = 9999999
const MATESCORE int = 100000

var searchAbandoned bool

var posSearched int


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
		//we are checkmated :(
		return -MATESCORE  //negative because being checkmated is bad
	} else {
		return 0  //draw
	}
}


func negamax(state *moves.GameState, isWhite bool, depth int, depthInx int, alpha int, beta int, moveChain string, timeLeft time.Duration) (int, moves.Move) {
	if timeLeft < 0 {
		//out of time
		searchAbandoned = true
		return 0, moves.Move{}
	}

	startTime := time.Now()
	
	if depth == 0 {return quiescenceSearch(state, isWhite, alpha, beta, timeLeft), moves.Move{}}

	posSearched++

	moveList, exists := getMoveList(depthInx, moveChain)

	if !exists {
		//moves were not cached, so we need to actually calculate them
		moveList = moves.GenerateAllMoves(state, false)
		orderMoves(state, moveList)

		appendToCache(depthInx, moveChain, moveList)  //Add the move to the cache
	}

	//TODO: draws by repetition and 50 move rule etc.
	if len(moveList) == 0 {return checkWin(state, isWhite), moves.Move{}}  //deal with checkmates and draws

	bestScore := -INF
	allocatedBestMove := false

	var bestMove moves.Move
	var bestInx int

	for inx, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		newChain := moveChain + hashMove(move)
		negScore, _ := negamax(state, !isWhite, depth - 1, depthInx + 1, -beta, -alpha, newChain, timeLeft - elapsed)
		score := -negScore

		moves.UnMakeLastMove(state)

		if score > bestScore || !allocatedBestMove {
			bestMove = move
			bestScore = score
			bestInx = inx
			allocatedBestMove = true
		}

		if score >= beta {break}  //prune

		if score > alpha {alpha = score}
	}

	if bestInx != 0 {
		//change the move we search first
		updateFirstMove(depthInx, moveChain, bestInx)
	}

	return bestScore, bestMove
}


func quiescenceSearch(state *moves.GameState, isWhite bool, alpha int, beta int, timeLeft time.Duration) int {
	//TODO: include checks here
	//IMPORTANT: this will not work with checkmates because eval does not return the checkmate values
	if timeLeft < 0 {
		searchAbandoned = true
		return 0
	}
	
	startTime := time.Now()

	staticEval := eval(state, isWhite)

	if staticEval >= beta {return beta}  //prune
	if staticEval > alpha {alpha = staticEval}

	posSearched++
	
	moveList := moves.GenerateAllMoves(state, true)
	orderMoves(state, moveList)

	for _, i := range moveList {
		moves.MakeMove(state, i)

		elapsed := time.Since(startTime)
		score := -quiescenceSearch(state, !isWhite, -beta, -alpha, timeLeft - elapsed)
	
		moves.UnMakeLastMove(state)

		if score >= beta {return beta}  //prune
		if score > alpha {alpha = score}
	}

	return alpha
}


func GetBestMove(state *moves.GameState) moves.Move {
	searchAbandoned = false
	clearCache()

	startTime := time.Now()
	timeLeft := time.Duration(time.Millisecond * 500)  //NOTE: change back to 500ms for testing

	depth := 1
	var bestMove moves.Move
	var elapsed time.Duration
	for timeLeft > 0 {
		posSearched = 0
		elapsed = time.Since(startTime)

		timeLeft -= elapsed

		score, searchBestMove := negamax(state, state.WhiteToMove, depth, 0, -INF, INF, "", timeLeft)  //NOTE: don't need to -score because this call is from the POV of the engine

		if !searchAbandoned {
			bestMove = searchBestMove
		} else {
			break
		}

		if score == MATESCORE {
			break
		}  //we have found a mate, so don't search any deeper

		depth++
	}

	if (bestMove == moves.Move{}) {
		panic("Best move is an empty move")
	}

	return bestMove
}