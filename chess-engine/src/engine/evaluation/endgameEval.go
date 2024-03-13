package evaluation


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

func getEndgameEval(perspective int, whiteKingPos int, blackKingPos int, mInfo *matInfo) int {
	eval := mInfo.matScore
	eval += mopUpScore(whiteKingPos, blackKingPos, perspective, mInfo)

	return eval
}