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
	
	if depth == 0 {return quiescenceSearch(state, isWhite, plyFromRoot, alpha, beta, timeLeft), &moves.Move{}}

	posSearched++

	moveList := moves.GenerateAllMoves(state, false)
	hashMove := searchTable.lookupMove(state.ZobristHash)

	orderMoves(state, moveList, hashMove)

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

	if !searchAbandoned {searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, alpha, nodeType, bestMove)}  //if the search was abandoned, the eval cannot be trusted

	return alpha, bestMove
}


func quiescenceSearch(state *board.GameState, isWhite bool, plyFromRoot int, alpha int, beta int, timeLeft time.Duration) int {
	if timeLeft < 0 {
		searchAbandoned = true
		return 0
	}

	if inCheck(state, isWhite) {
		//if we are in check, search through all evasive moves. This helps stop obvious blunders and detect mates.
		eval, _ := negamax(state, isWhite, 1, plyFromRoot, []*moves.Move{}, alpha, beta, timeLeft)
		return eval
	}
	
	startTime := time.Now()

	staticEval := evaluation.Eval(state, isWhite)

	if staticEval >= beta {return beta}  //prune
	if staticEval > alpha {alpha = staticEval}

	posSearched++
	
	moveList := moves.GenerateAllMoves(state, true)
	hashMove := searchTable.lookupMove(state.ZobristHash)

	orderMoves(state, moveList, hashMove)

	for _, move := range moveList {
		moves.MakeMove(state, move)

		elapsed := time.Since(startTime)
		score := -quiescenceSearch(state, !isWhite, plyFromRoot + 1, -beta, -alpha, timeLeft - elapsed)
	
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

		if searchAbandoned {
			break
		} else {
			bestMove = searchBestMove
			searchedDepthOne = true
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


//6k1/p4pp1/2p1pn2/8/6p1/2p5/r7/1r1N1K2 b - - 1 56
//position fen "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" moves e2e4 b8c6 d2d4 g8f6 e4e5 f6d5 c2c4 d5b6 c4c5 b6d5 f1c4 e7e6 b1c3 d5c3 b2c3 d7d5 c4b5 c8d7 d1f3 c6e5 b5d7 e5d7 c5c6 b7c6 c3c4 f8b4 e1f1 e8g8 f3e2 b4c3 c1b2 c3b2 e2b2 d5c4 a1c1 d7b6 h2h4 a8b8 h1h3 b6d5 b2a1 d5f4 h3h2 d8d5 g2g3 f4d3 c1d1 b8d8 g1e2 f8e8 e2c3 d5d4 c3e2 d4d5 a1b1 c6c5 e2f4 d5f3 f4h3 d8b8 b1c2 b8b2 c2a4 c7c6 f1g1 h7h6 d1f1 c4c3 a4a5 c3c2 a5c3 f3d1 f1d1 c2d1q g1g2 d1e2 h2h1 d3f2 h1e1 f2d1 e1e2 b2e2 g2f3 d1c3 a2a3 e2a2 h4h5 a2a3 h3f2 c3d5 f3g2 d5f6 f2d1 f6h5 g3g4 h5f6 d1f2 c5c4 g4g5 h6g5 f2d1 c4c3 d1e3 a3a2 g2f1 e8b8 e3d1 b8b1 f1e1 g5g4 e1f1 g4g3 f1e1 a2b2 e1f1 g3g2 f1g1 b1d1 g1h2 g2g1q h2h3 b2b3 h3h4 g1h2 h4g5 f6d7 g5g4 h2g1 g4h5 b3b5 h5h4 b5a5 h4h3 a5b5 h3h4 b5a5