package search


import (
	"fmt"
	"time"
	"gambot/src/engine/moves"
	"gambot/src/engine/board"
	"gambot/src/engine/evaluation"
)


const inf int = 9999999
const matescore int = 100000
const mateThreshold int = matescore - 1024

var searchAbandoned bool

var posSearched int
var ttLookups int


func inCheck(state *board.GameState, isWhite bool) bool {
	kingPos := board.PieceLists.BlackKingPos
	if isWhite {kingPos = board.PieceLists.WhiteKingPos}

	//set bitboard at king's position
	var kingPosBB uint64
	kingPosBB |= 1 << kingPos

	return (kingPosBB & state.Bitboards.AttackedSquares) != 0
}


func checkWin(state *board.GameState, plyFromRoot int, isWhite bool) int {
	//check who has won, or if it is a draw - assumes that the player has no legal moves
	if inCheck(state, isWhite) {
		//we are checkmated :(
		return -matescore + plyFromRoot //negative because being checkmated is bad, also a larger ply from root is good
	} else {
		return 0  //draw
	}
}


func negamax(state *board.GameState, isWhite bool, depth int, plyFromRoot int, alpha int, beta int, timeLeft time.Duration) (int, *moves.Move) {
	//core minimax/negamax search routine
	if timeLeft < 0 {
		//out of time
		searchAbandoned = true
		return 0, &moves.Move{}
	}

	startTime := time.Now()

	ttSuccess, ttEval := searchTable.lookupEval(state.ZobristHash, depth, plyFromRoot, alpha, beta)

	if ttSuccess {
		//this position is in transposition table. We don't need to search it again
		ttLookups++
		bestMove := searchTable.lookupMove(state.ZobristHash)

		return ttEval, bestMove
	}
	
	if depth == 0 {return quiescenceSearch(state, isWhite, plyFromRoot, alpha, beta, timeLeft), &moves.Move{}}

	posSearched++

	moveList := moves.GenerateAllMoves(state, false)
	hashMove := searchTable.lookupMove(state.ZobristHash)

	orderMoves(state, moveList, hashMove, plyFromRoot)

	//TODO: draws by repetition and 50 move rule etc.
	if len(moveList) == 0 {return checkWin(state, plyFromRoot, isWhite), &moves.Move{}}  //deal with checkmates and draws

	nodeType := allNode
	bestMove := &moves.Move{}

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		negScore, _ := negamax(state, !isWhite, depth - 1, plyFromRoot + 1, -beta, -alpha, timeLeft - elapsed)
		score := -negScore

		moves.UnMakeLastMove(state)

		if searchAbandoned {return 0, bestMove}

		//fail-hard cutoff (prune position)
		if score >= beta {
			searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, beta, cutNode, bestMove)  //we do not actually know if the value of bestMove is best for this position
			addKiller(move, plyFromRoot)

			return beta, bestMove
		}

		//new best move found
		if score > alpha {
			alpha = score
			bestMove = move
			nodeType = pvNode
		}
	}

	searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, alpha, nodeType, bestMove)

	return alpha, bestMove
}


func quiescenceSearch(state *board.GameState, isWhite bool, plyFromRoot int, alpha int, beta int, timeLeft time.Duration) int {
	//search through all capture moves until a quiet position is reached
	if timeLeft < 0 {
		searchAbandoned = true
		return 0
	}

	if inCheck(state, isWhite) {
		//if we are in check, search through all evasive moves. This helps stop obvious blunders and detect mates.
		eval, _ := negamax(state, isWhite, 1, plyFromRoot, alpha, beta, timeLeft)  //TODO: actually do the moveLines
		return eval
	}
	
	startTime := time.Now()

	staticEval := evaluation.Eval(state, isWhite)

	if staticEval >= beta {return beta}  //prune
	if staticEval > alpha {alpha = staticEval}

	posSearched++
	
	moveList := moves.GenerateAllMoves(state, true)
	hashMove := searchTable.lookupMove(state.ZobristHash)

	orderMoves(state, moveList, hashMove, plyFromRoot)

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		score := -quiescenceSearch(state, !isWhite, plyFromRoot + 1, -beta, -alpha, timeLeft - elapsed)
	
		moves.UnMakeLastMove(state)

		if searchAbandoned {return 0}

		if score >= beta {return beta}  //prune (NOTE: no need to add killer move because captures are treated differently anyway)

		if score > alpha {alpha = score}
	}

	return alpha
}


func getAnyMove(state *board.GameState) *moves.Move {
	//when we have not even searched depth 1, get any legal move
	moveList := moves.GenerateAllMoves(state, false)
	return moveList[0]
}


func getPvLine(state *board.GameState, prevPvMove *moves.Move, pvLine *[]*moves.Move, depthSearched int) {
	//use the transposition table to get the pv line
	if len(*pvLine) > depthSearched {return}  //we cannot have a pv line longer than the depth we searched it to

	*pvLine = append(*pvLine, prevPvMove)

	moves.MakeMove(state, prevPvMove)

	pvMove := searchTable.lookupPvMove(state.ZobristHash)
	if pvMove != nil {getPvLine(state, pvMove, pvLine, depthSearched)}

	moves.UnMakeLastMove(state)
}


func uciSearchInfo(depth int, score int, nodes int, ttLookups int, timeMs int64, pvLine []*moves.Move) {
	//send search info in the format required by UCI
	var pvMoves string
	for i, x := range pvLine {
		if i > 0 {pvMoves += " "}
		pvMoves += x.MoveStr()
	}

	fmt.Print("info")

	fmt.Printf(" depth %v", depth)
	fmt.Printf(" score cp %v", score)
	fmt.Printf(" nodes %v", nodes)
	fmt.Printf(" tthits %v", ttLookups)  //NOTE: this is not in UCI, but is useful for debugging
	fmt.Printf(" time %v", timeMs)
	fmt.Printf(" pv %s", pvMoves)
	fmt.Print("\n")
}


func GetBestMove(state *board.GameState, moveTime int) *moves.Move {
	//run iterative deepening to get the (objectively...) best move
	searchAbandoned = false

	startTime := time.Now()
	timeLeft := time.Duration(time.Millisecond * time.Duration(moveTime))

	depth := 1
	searchedDepthOne := false

	pvLine := []*moves.Move{}

	var posScore int
	var bestMove *moves.Move
	var elapsed time.Duration

	for timeLeft > 0 {
		posSearched = 0
		ttLookups = 0

		elapsed = time.Since(startTime)

		timeLeft -= elapsed

		score, searchBestMove := negamax(state, state.WhiteToMove, depth, 0, -inf, inf, timeLeft)  //NOTE: don't need to -score because this call is from the POV of the engine
		
		if searchBestMove.PieceValue != 0 {
			posScore = score
			bestMove = searchBestMove
			searchedDepthOne = true

			pvLine = []*moves.Move{}
			getPvLine(state, searchBestMove, &pvLine, depth)
		}

		uciSearchInfo(depth, score, posSearched, ttLookups, elapsed.Milliseconds(), pvLine)

		if searchAbandoned {break}

		depth++
	}

	if !searchedDepthOne {bestMove = getAnyMove(state)}
	if bestMove.PieceValue == 0 {panic("Best move is an empty move")}  //NOTE: this may be because we are in mate

	//the search may have been cancelled before we could store the result from the final depth in the tt, so let's do it now O-O.
	searchTable.storeEntry(state.ZobristHash, depth, 0, posScore, pvNode, bestMove)

	return bestMove
}