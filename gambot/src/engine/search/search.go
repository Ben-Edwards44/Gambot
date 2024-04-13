package search


import (
	"fmt"
	"time"
	"gambot/src/engine/moves"
	"gambot/src/engine/board"
	"gambot/src/engine/evaluation"
)


const inf int = 9999999
const mateScore int = 100000
const mateThreshold int = mateScore - 1024

const reductionCutoff int = 4  //the move at which moves start getting reduced
const steepReductionCutoff int = 6  //the move at which moves start being more heavily reduced

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
		return -mateScore + plyFromRoot //negative because being checkmated is bad, also a larger ply from root is good
	} else {
		return 0  //draw
	}
}


func getReduction(state *board.GameState, isWhite bool, isCapt bool, isProm bool, moveInx int, depth int) int {
	//Late move reductions - we assume moves at the end of the ordered list to be bad, so we can search them at a lower depth. (https://www.chessprogramming.org/Late_Move_Reductions)
	check := inCheck(state, isWhite)

	//TODO: tune the values, or use a non-linear function
	if depth < 3 || moveInx < reductionCutoff || isCapt || isProm || check {return 0}  //we don't want to reduce potentially good moves

	if moveInx < steepReductionCutoff {
		return 1
	} else {
		return depth / 3
	}
}


func negamax(state *board.GameState, isWhite bool, depth int, plyFromRoot int, alpha int, beta int, timeLeft time.Duration) (int, *moves.Move) {
	//core minimax/negamax search routine
	if timeLeft < 0 {
		//out of time
		searchAbandoned = true
		return 0, nil
	}

	startTime := time.Now()

	//check for draw by threefold repetition.
	//NOTE: we only need to check for 1 repetition in the search because there is no reason we would choose differently if we saw the same pos again
	if plyFromRoot > 0 && repTable.seen(state.ZobristHash) {return 0, nil}

	ttSuccess, ttEval := searchTable.lookupEval(state.ZobristHash, depth, plyFromRoot, alpha, beta)

	if ttSuccess {
		//this position is in transposition table. We don't need to search it again
		ttLookups++
		bestMove := searchTable.lookupMove(state.ZobristHash)

		return ttEval, bestMove
	}

	//mate distance pruning - a mate has been found in n plies, so we can prune if we have not yet found a mate in <= n plies
	if alpha > mateThreshold && plyFromRoot > 0 {
		currentMateScore := mateScore - plyFromRoot

		if currentMateScore < beta {beta = currentMateScore}
		if alpha >= currentMateScore {return currentMateScore, nil}
	}
	
	if depth == 0 {return quiescenceSearch(state, isWhite, plyFromRoot, alpha, beta, timeLeft), nil}

	posSearched++

	moveList := moves.GenerateAllMoves(state, false)
	hashMove := searchTable.lookupMove(state.ZobristHash)

	orderMoves(state, moveList, hashMove, plyFromRoot)

	//check for wins and draws
	if len(moveList) == 0 {return checkWin(state, plyFromRoot, isWhite), nil}

	if plyFromRoot > 0 {repTable.push(state.ZobristHash)}

	nodeType := allNode
	
	var bestMove *moves.Move

	for inx, move := range moveList {
		moveIsCapt := state.Board[move.EndX * 8 + move.EndY] != 0
		moveIsProm := move.PromotionValue != 0

		elapsed := time.Since(startTime)

		//PVS search - due to move ordering, we assume the best move to be the first.
		//Therefore, we can search all other moves with a smaller window and just make sure they are as bad as we think they are.
		var score int
		if inx == 0 {
			moves.MakeMove(state, move)

			negScore, _ := negamax(state, !isWhite, depth - 1, plyFromRoot + 1, -beta, -alpha, timeLeft - elapsed)
			score = -negScore
		} else {
			reduction := getReduction(state, isWhite, moveIsCapt, moveIsProm, inx, depth)

			moves.MakeMove(state, move)  //getting the reduction relies on the current board state, so only make the move after getting the reduction

			reducedDepth := depth - reduction - 1
			if reducedDepth < 0 {reducedDepth = 0}

			negScore, _ := negamax(state, !isWhite, reducedDepth, plyFromRoot + 1, -alpha - 1, -alpha, timeLeft - elapsed)
			score = -negScore

			needFullSearch := score > alpha
			if needFullSearch {
				//this move is not as bad as we first thought; we need to search it fully without reductions and with a full window
				negScore, _ := negamax(state, !isWhite, depth - 1, plyFromRoot + 1, -beta, -alpha, timeLeft - elapsed)
				score = -negScore
			}
		}

		moves.UnMakeLastMove(state)

		if searchAbandoned {return 0, bestMove}

		if plyFromRoot > 0 && repTable.seenHashes[repTable.length - 1] != state.ZobristHash {panic("Reptable not popped")}  //for testing

		//fail-hard cutoff (prune position)
		if score >= beta {
			if plyFromRoot > 0 {repTable.pop()}
			
			searchTable.storeEntry(state.ZobristHash, depth, plyFromRoot, beta, cutNode, bestMove)  //we do not actually know if the value of bestMove is best for this position
			
			if !moveIsCapt && !moveIsProm {addKiller(move, plyFromRoot)}  //promotions and captures are ordered before killers anyway, so only add quiet moves to killers

			return beta, bestMove
		}

		//new best move found
		if score > alpha {
			alpha = score
			bestMove = move
			nodeType = pvNode
		}
	}

	if plyFromRoot > 0 {repTable.pop()}
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
		eval, _ := negamax(state, isWhite, 1, plyFromRoot, alpha, beta, timeLeft)
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


func playPvLine(state *board.GameState, prevPvMove *moves.Move, pvLine *[]*moves.Move, depthSearched int) {
	//use the transposition table to get the pv line
	if len(*pvLine) > depthSearched {return}  //we cannot have a pv line longer than the depth we searched it to

	*pvLine = append(*pvLine, prevPvMove)

	moves.MakeMove(state, prevPvMove)

	pvMove := searchTable.lookupPvMove(state.ZobristHash)
	if pvMove != nil {playPvLine(state, pvMove, pvLine, depthSearched)}

	moves.UnMakeLastMove(state)
}


func getPvLine(state *board.GameState, bestMove *moves.Move, depthSearched int) (pvLine *[]*moves.Move) {
	pvLine = &[]*moves.Move{}
	
	defer func() {
		//In some rare cases, the tt will be overwritten by another pv node.
		//This will give a non-legal move, causing a panic when we try to make this move it.
		//Since the pv is not used, we can just ignore the panic (x_x) and not crash the entire engine.
		r := recover()

		if r != nil {
			*pvLine = (*pvLine)[:len(*pvLine) - 1]

			for i := 0; i <= len(*pvLine); i++ {
				moves.UnMakeLastMove(state)
			}
		}
	}()

	playPvLine(state, bestMove, pvLine, depthSearched)

	return  //NOTE: named return value has been used
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

	initRepTable(state.ZobristHash)

	var bestMove *moves.Move
	var elapsed time.Duration

	for timeLeft > 0 {
		posSearched = 0
		ttLookups = 0

		elapsed = time.Since(startTime)

		timeLeft -= elapsed

		score, searchBestMove := negamax(state, state.WhiteToMove, depth, 0, -inf, inf, timeLeft)  //NOTE: don't need to -score because this call is from the POV of the engine
		
		if searchBestMove != nil {
			bestMove = searchBestMove
			searchedDepthOne = true

			pvLine = *getPvLine(state, searchBestMove, depth)
		}

		uciSearchInfo(depth, score, posSearched, ttLookups, elapsed.Milliseconds(), pvLine)

		if searchAbandoned {break}
		if score >= mateScore - depth {break}  //we have found a mate in the given depth, no need to look further

		depth++
	}

	if !searchedDepthOne {bestMove = getAnyMove(state)}
	if bestMove.PieceValue == 0 {panic("Best move is an empty move")}  //NOTE: this may be because we are in mate

	return bestMove
}