package evaluation


import "chess-engine/src/engine/board"


var gameMatInfo matInfo


type matInfo struct {
	whiteMat int
	blackMat int
	matScore int
	gamePhase int
}


func (info *matInfo) updateMatInfo(state *board.GameState, perspective int) {
	whiteMat := 0
	blackMat := 0
	phase := totalPhase

	wPieces := &board.PieceLists.WhitePieceSquares
	bPieces := &board.PieceLists.BlackPieceSquares

	for i := 0; i < len(wPieces); i++ {
		wPos := wPieces[i]
		bPos := bPieces[i]

		if wPos != -1 {
			wPiece := state.Board[wPos]

			whiteMat += pieceWeights[wPiece - 1]
			phase -= phaseWeights[wPiece - 1]
		}
		if bPos != -1 {
			bPiece := state.Board[bPos]

			blackMat += pieceWeights[bPiece - 7]
			phase -= phaseWeights[bPiece - 7]
		}
	}

	info.whiteMat = whiteMat
	info.blackMat = blackMat
	info.matScore = (whiteMat - blackMat) * perspective
	info.gamePhase = (phase * 256 + (totalPhase / 2)) / totalPhase  //I have no idea why this works. I got it from here: https://www.chessprogramming.org/Tapered_Eval
}