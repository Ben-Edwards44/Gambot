package search


import "chess-engine/src/engine/moves"


func quickSort(moveList []moves.Move, moveScores []int, low int, high int) {
	if low < high {
		pivot := partition(moveList, moveScores, low, high)

		quickSort(moveList, moveScores, low, pivot - 1)
		quickSort(moveList, moveScores, pivot + 1, high)
	}
}


func partition(moveList []moves.Move, moveScores []int, low int, high int) int {
	pivot := moveScores[high]
	i := low - 1

	for j := low; j < high; j++ {
		if moveScores[j] > pivot {
			i++

			//swap elements
			moveList[i], moveList[j] = moveList[j], moveList[i]
			moveScores[i], moveScores[j] = moveScores[j], moveScores[i]
		}
	}

	//swap last element
	moveList[i + 1], moveList[high] = moveList[high], moveList[i + 1]
	moveScores[i + 1], moveScores[high] = moveScores[high], moveScores[i + 1]

	return i + 1
}


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