package search


import (
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/evaluation"
	"chess-engine/src/engine/moves"
	"fmt"
	"time"
)


const inf int = 9999999
const matescore int = 100000
const mateThreshold int = matescore - 1024

var searchAbandoned bool

var bestMoves map[uint64]*moves.Move

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


func negamax(state *board.GameState, isWhite bool, depth int, plyFromRoot int, pvLine []*moves.Move, alpha int, beta int, timeLeft time.Duration) (int, *moves.Move) {
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

		var bestMove *moves.Move
		if plyFromRoot == 0 {
			bestMove = searchTable.lookupMove(state.ZobristHash)
		}
		
		return ttEval, bestMove
	}
	
	if depth == 0 {return quiescenceSearch(state, isWhite, alpha, beta, timeLeft), &moves.Move{}}

	posSearched++

	moveList := moves.GenerateAllMoves(state, false)
	prevBestMove := bestMoves[state.ZobristHash]

	orderMoves(state, moveList, prevBestMove, searchMoveOrder)

	//TODO: draws by repetition and 50 move rule etc.
	if len(moveList) == 0 {return checkWin(state, plyFromRoot, isWhite), &moves.Move{}}  //deal with checkmates and draws

	nodeType := allNode

	var bestMove *moves.Move

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		negScore, _ := negamax(state, !isWhite, depth - 1, plyFromRoot + 1, pvLine, -beta, -alpha, timeLeft - elapsed)
		score := -negScore

		moves.UnMakeLastMove(state)

		//fail-hard cutoff (prune position)
		if score >= beta {
			if !searchAbandoned {
				searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, beta, cutNode, bestMove)  //we do not actually know if the value of bestMove is best for this position
			}

			return beta, bestMove
		}  

		//new best move found
		if score > alpha {
			alpha = score
			bestMove = move
			nodeType = pvNode

			if len(pvLine) > 0 {pvLine[plyFromRoot] = bestMove}  //NOTE: this will be overwritten later if a better move is found
		}
	}

	bestMoves[state.ZobristHash] = bestMove

	if !searchAbandoned {searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, alpha, nodeType, bestMove)}  //if the search was abandoned, the eval cannot be trusted

	return alpha, bestMove
}


func quiescenceSearch(state *board.GameState, isWhite bool, alpha int, beta int, timeLeft time.Duration) int {
	//TODO: include checks here
	//IMPORTANT: this will not work with checkmates because eval does not return the checkmate values
	if timeLeft < 0 {
		searchAbandoned = true
		return 0
	}

	if inCheck(state, isWhite) {
		//if we are in check, search through all evasive moves. This helps stop obvious blunders.
		eval, _ := negamax(state, isWhite, 1, 0, []*moves.Move{}, alpha, beta, timeLeft)
		return eval
	}
	
	startTime := time.Now()

	staticEval := evaluation.Eval(state, isWhite)

	if staticEval >= beta {return beta}  //prune
	if staticEval > alpha {alpha = staticEval}

	posSearched++
	
	moveList := moves.GenerateAllMoves(state, true)
	orderMoves(state, moveList, &moves.Move{}, quiSearchMoveOrder)  //TODO: maybe prev best move here??

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


func uciSearchInfo(depth int, score int, nodes int, timeMs int64, pvLine []*moves.Move) {
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
	fmt.Printf(" time %v", timeMs)
	fmt.Printf(" pv %s", pvMoves)
	fmt.Print("\n")
}


func GetBestMove(state *board.GameState, moveTime int) *moves.Move {
	searchAbandoned = false
	bestMoves = make(map[uint64]*moves.Move)

	startTime := time.Now()
	timeLeft := time.Duration(time.Millisecond * time.Duration(moveTime))  //NOTE: change to 500ms for testing

	depth := 1

	searchedDepthOne := false

	var pvLine []*moves.Move

	var bestMove *moves.Move
	var elapsed time.Duration
	for timeLeft > 0 {
		posSearched = 0
		ttLookups = 0

		pvLine = append(pvLine, &moves.Move{})

		elapsed = time.Since(startTime)

		timeLeft -= elapsed

		score, searchBestMove := negamax(state, state.WhiteToMove, depth, 0, pvLine, -inf, inf, timeLeft)  //NOTE: don't need to -score because this call is from the POV of the engine

		uciSearchInfo(depth, score, posSearched, elapsed.Milliseconds(), pvLine)

		if !searchAbandoned {
			bestMove = searchBestMove
			searchedDepthOne = true
		} else {
			break
		}

		depth++
	}

	if !searchedDepthOne {
		//could not even search to depth one, just play the first available move
		moveList := moves.GenerateAllMoves(state, false)
		bestMove = moveList[0]
	}

	if bestMove.PieceValue == 0 {panic("Best move is an empty move")}

	return bestMove  //NOTE: if no moves are available (in mate), this will be an empty move
}