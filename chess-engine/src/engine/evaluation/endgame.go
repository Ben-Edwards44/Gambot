package evaluation


func mopUpScore(whiteKingPos int, blackKingPos int, perspective int) int {
	//in endgames where we are up material, we want kings to be close together and near corners (assumes mat info has been updated)
	friendKingPos := whiteKingPos
	enemyKingPos := blackKingPos
	if perspective == -1 {
		friendKingPos = blackKingPos
		enemyKingPos = whiteKingPos
	}
	
	score := 0
	if gameMatInfo.matScore >= 2 * pawnWeight {
		score += centerManhattanDist[enemyKingPos]  //we want the enemy king far away from the center
		score += 14 - squareManhattanDists[friendKingPos * 64 + enemyKingPos]  //we want our king close to the enemy king
	}

	return score
}