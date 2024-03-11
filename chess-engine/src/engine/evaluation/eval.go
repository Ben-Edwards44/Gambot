package evaluation


import "chess-engine/src/engine/board"


func getTaperedEval(openingEval int, endgameEval int, phase int) int {
	return ((openingEval * (256 - phase)) + (endgameEval * phase)) / 256  //https://mediocrechess.blogspot.com/2011/10/guide-tapered-eval.html
}


func Eval(state *board.GameState, whiteToMove bool) int {
	//NOTE: with negamax, the eval should always be in the perspective of the current player (so times by -1 for black)
	//eval = (good for white - good for black) * perspective
	
	perspective := -1
	if whiteToMove {perspective = 1}
	
	gameMatInfo.updateMatInfo(state, perspective)

	openingEval := 0  //TODO: get opening eval
	endgameEval := 0//mopUpScore(board.PieceLists.WhiteKingPos, board.PieceLists.BlackKingPos, perspective)

	openingEval += gameMatInfo.matScore
	endgameEval += gameMatInfo.matScore

	finalEval := getTaperedEval(openingEval, endgameEval, gameMatInfo.gamePhase)

	return finalEval
}