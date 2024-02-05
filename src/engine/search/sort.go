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