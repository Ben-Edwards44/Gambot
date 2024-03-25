package search


import (
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/moves"
)


//prioritise attacking the most valuable pieces (MVV) with the least valuable piece (LVA): https://www.chessprogramming.org/MVV-LVA
//the actual values are taken from https://github.com/likeawizard/tofiks
var mvvLva [6 * 6]int = [6 * 6]int {
	10, 9, 8, 7, 6, 5,       //pawn victim
	30, 29, 28, 27, 26, 25,  //bishop victim
	20, 19, 18, 17, 16, 15,  //knight victim
	40, 39, 38, 37, 36, 35,  //rook victim
	0, 0, 0, 0, 0, 0,        //king victim
	50, 49, 48, 47, 46, 45,  //queen victim
}


const hashMoveScore int = 10000
const promotionOffset int = 1000
const mvvLvaOffset int = 100


func quickSort(moveList []*moves.Move, moveScores []int, low int, high int) {
	if low < high {
		pivot := partition(moveList, moveScores, low, high)

		quickSort(moveList, moveScores, low, pivot - 1)
		quickSort(moveList, moveScores, pivot + 1, high)
	}
}


func partition(moveList []*moves.Move, moveScores []int, low int, high int) int {
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


func scoreMove(state *board.GameState, move *moves.Move, hashMove *moves.Move) int {
	//moves are ordered as follows: hash move / pv move (from tt), promotions, MVV/LVA for captures, quiet moves
	
	captVal := state.Board[move.EndX * 8 + move.EndY]
	if move.EnPassant {captVal = 1}
	
	if move == hashMove {
		//this also implicitly accounts for pv moves
		return hashMoveScore
	} else if move.PromotionValue > 0 {
		promVal := move.PromotionValue
		if promVal > 6 {promVal -= 6}

		return promotionOffset + promVal
	} else if captVal > 0 {
		victimInx := captVal - 1
		if captVal > 6 {victimInx -= 6}

		aggressInx := move.PieceValue - 1
		if move.PieceValue > 6 {aggressInx -= 6}

		return mvvLvaOffset + mvvLva[victimInx * 6 + aggressInx]
	} else {
		//not a pv move, hash move, promotion or capture. Just a regular ol' move (TODO: add history heuristic)
		score := 0
		
		posBB := uint64(1 << (move.EndX * 8 + move.EndY))
		if posBB & state.Bitboards.PawnAttacks != 0 {score -= move.PieceValue}  //moving to a square attacked by enemy pawn is not good
	
		return score
	}
}


func orderMoves(state *board.GameState, moveList []*moves.Move, hashMove *moves.Move) {
	//slices are passed by reference, so no need to return

	var moveScores []int
	for _, i := range moveList {
		moveScores = append(moveScores, scoreMove(state, i, hashMove))  //get the move's score
	}

	quickSort(moveList, moveScores, 0, len(moveList) - 1)
}