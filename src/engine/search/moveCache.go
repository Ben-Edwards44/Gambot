package search

import (
	"strconv"
	"chess-engine/src/engine/moves"
)

//TODO: use bits to represent moves and use array instead of map
var moveCache []map[string][]moves.Move


func hashMove(move moves.Move) string {
	//convert move to a string so a chain of moves can be used as a map key. Spaces (or other seperators) are needed becuase some numbers are 2 digits long, so can overlap.
	
	sx := strconv.Itoa(move.StartX)
	sy := strconv.Itoa(move.StartY)
	ex := strconv.Itoa(move.EndX)
	ey := strconv.Itoa(move.EndY)
	piece := strconv.Itoa(move.PieceValue)
	dp := strconv.FormatBool(move.DoublePawnMove)
	ep := strconv.FormatBool(move.EnPassant)
	kc := strconv.FormatBool(move.KingCastle)
	qc := strconv.FormatBool(move.QueenCastle)
	pr := strconv.Itoa(move.PromotionValue)

	return " " + sx + " " + sy + " " + ex + " " + ey + " " + piece + " " + dp + " " + ep + " " + kc + " " + qc + " " + pr + " " 
}


func getMoveList(depthInx int, moveChain string) ([]moves.Move, bool) {
	//get the ordered list of moves from a previous iteration of iterative deepening

	if depthInx >= len(moveCache) {return []moves.Move{}, false}

	moveList, exists := moveCache[depthInx][moveChain]

	return moveList, exists
}


func appendToCache(depthInx int, moveChain string, moveList []moves.Move) {
	if depthInx < len(moveCache) {
		//just adding a new move onto an existing depth
		moveCache[depthInx][moveChain] = moveList
	} else {
		//first move on a new depth

		newMap := make(map[string][]moves.Move)
		newMap[moveChain] = moveList

		moveCache = append(moveCache, newMap)
	}
}


func updateFirstMove(depthInx int, moveChain string, newFirstInx int) {
	//update the move order with the new best move at the front

	list := moveCache[depthInx][moveChain]  //NOTE: slices are passed by reference

	moveVal := list[newFirstInx]

	//shift other elements
	for i := newFirstInx - 1; i >= 0; i-- {
		list[i + 1] = list[i]
	}

	list[0] = moveVal  //shift the new best into first place
}


func clearCache() {
	moveCache = make([]map[string][]moves.Move, 1)

	newMap := make(map[string][]moves.Move)
	moveCache[0] = newMap
}