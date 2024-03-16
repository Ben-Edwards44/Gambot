package evaluation


import "chess-engine/src/engine/board"


var gameMatInfo matInfo


type matInfo struct {
	whiteMat int
	blackMat int
	matScore int

	whiteOpTbl int
	whiteEgTbl int
	blackOpTbl int
	blackEgTbl int

	opPieceSqScore int
	egPieceSqScore int

	gamePhase int
}


func (info *matInfo) updateMatInfo(state *board.GameState, perspective int) {
	//loops through piece positions and calculates material advantage, piece square table score, and game phase

	whiteMat := 0
	blackMat := 0

	wOpTbl := 0
	wEgTbl := 0
	bOpTbl := 0
	bEgTbl := 0

	phase := totalPhase

	wPieces := &board.PieceLists.WhitePieceSquares
	bPieces := &board.PieceLists.BlackPieceSquares

	for i := 0; i < len(wPieces); i++ {
		wPos := wPieces[i]
		bPos := bPieces[i]

		if wPos != -1 {
			wPiece := state.Board[wPos]
			inx := wPiece - 1

			whiteMat += pieceWeights[inx]
			phase -= phaseWeights[inx]

			wOpTbl += opPieceSqTables[inx][wPos]
			wEgTbl += egPieceSqTables[inx][wPos]
		}
		if bPos != -1 {
			bPiece := state.Board[bPos]
			inx := bPiece - 7

			blackMat += pieceWeights[inx]
			phase -= phaseWeights[inx]

			tblInx := bPos ^ 56  //this is needed because the piece square tables are from the POV of white
			bOpTbl += opPieceSqTables[inx][tblInx]
			bEgTbl += egPieceSqTables[inx][tblInx]
		}
	}

	info.whiteMat = whiteMat
	info.blackMat = blackMat
	info.matScore = (whiteMat - blackMat) * perspective

	info.whiteOpTbl = wOpTbl
	info.whiteEgTbl = wEgTbl
	info.blackOpTbl = bOpTbl
	info.blackEgTbl = bEgTbl

	info.opPieceSqScore = (wOpTbl - bOpTbl) * perspective
	info.egPieceSqScore = (wEgTbl - bEgTbl) * perspective

	info.gamePhase = (phase * 256 + (totalPhase / 2)) / totalPhase  //I have no idea why this works. I got it from here: https://www.chessprogramming.org/Tapered_Eval
}


func mopUpScore(whiteKingPos int, blackKingPos int, perspective int, mInfo *matInfo) int {
	//in endgames where we are up material, we want kings to be close together and near corners
	friendKingPos := whiteKingPos
	enemyKingPos := blackKingPos
	if perspective == -1 {
		friendKingPos = blackKingPos
		enemyKingPos = whiteKingPos
	}
	
	score := 0
	if mInfo.matScore >= 2 * pawnWeight {
		score += centerManhattanDist[enemyKingPos]  //we want the enemy king far away from the center
		score += 14 - squareManhattanDists[friendKingPos * 64 + enemyKingPos]  //we want our king close to the enemy king
	}

	return score
}


func openingEval(mInfo *matInfo) int {
	//assumes mInfo has been updated
	eval := mInfo.matScore
	eval += mInfo.opPieceSqScore

	return eval
}


func endgameEval(mInfo *matInfo, perspective int) int {
	//assumes mInfo has been updated
	eval := mInfo.matScore
	eval += mInfo.egPieceSqScore

	eval += mopUpScore(board.PieceLists.WhiteKingPos, board.PieceLists.BlackKingPos, perspective, mInfo)

	return eval
}


func getTaperedEval(openingEval int, endgameEval int, phase int) int {
	return ((openingEval * (256 - phase)) + (endgameEval * phase)) / 256  //https://mediocrechess.blogspot.com/2011/10/guide-tapered-eval.html
}


func Eval(state *board.GameState, whiteToMove bool) int {
	//NOTE: with negamax, the eval should always be in the perspective of the current player (so times by -1 for black)
	//eval = (good for white - good for black) * perspective
	
	perspective := -1
	if whiteToMove {perspective = 1}
	
	gameMatInfo.updateMatInfo(state, perspective)

	openingEval := openingEval(&gameMatInfo)
	endgameEval := endgameEval(&gameMatInfo, perspective)

	finalEval := getTaperedEval(openingEval, endgameEval, gameMatInfo.gamePhase)

	return finalEval
}