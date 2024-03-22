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


func searchMoveOrder(state *board.GameState, move *moves.Move) int {
	score := 0

	currentPiece := move.PieceValue
	promotion := move.PromotionValue
	captPiece := state.Board[move.EndX * 8 + move.EndY] - 6  //-6 because it is opposite colour to current

	if move.EnPassant {captPiece = 7}  //special case
	
	if currentPiece > 6 {
		currentPiece -= 6
		promotion -= 6
		captPiece += 6
	}

	if captPiece > 0 {score += captPiece - currentPiece}  //capturing high value pieces with low value ones is good

	score += promotion  //promotions are good (if not promotion, this will just add 0 to score)

	//var posBB uint64
	//posBB |= 1 << (move.StartX * 8 + move.StartY)
	//if posBB & state.NoKingMoveBitBoard != 0 {score -= currentPiece}  //moving to an attacked square is not good

	var posBB uint64
	posBB = 1 << (move.EndX * 8 + move.EndY)
	if posBB & state.Bitboards.AttackedSquares != 0 {score -= currentPiece}  //moving to an attacked square is not good

	posBB = 1 << (move.StartX * 8 + move.StartY)
	if posBB & state.Bitboards.AttackedSquares != 0 {score += currentPiece}  //moving out of an attacked square is good

	return score
}


func quiSearchMoveOrder(state *board.GameState, move *moves.Move) int {
	//NOTE: assumes move is a capture

	return searchMoveOrder(state, move)

	/*
	score := 0

	pieceVal := move.PieceValue
	captPieceVal := state.Board[move.EndX * 8 + move.EndY]
	aggressInx := pieceVal - 1
	victimInx := captPieceVal - 7
	if pieceVal > 6 {
		//black to move
		aggressInx = pieceVal - 7
		victimInx = captPieceVal - 1
	}

	if move.EnPassant {victimInx = 0}  //special case where captPieceVal will be 0

	score += mvvLva[victimInx * 6 + aggressInx]

	//maybe include this code to check whether the square is protected
	var posBB uint64
	posBB = 1 << (move.EndX * 8 + move.EndY)
	if posBB & state.NoKingMoveBitBoard != 0 {score -= pieceVal}  //moving to an attacked square is not good

	posBB = 1 << (move.StartX * 8 + move.StartY)
	if posBB & state.NoKingMoveBitBoard != 0 {score += pieceVal}  //moving out of an attacked square is good

	return score
	*/
}


func orderMoves(state *board.GameState, moveList []*moves.Move, prevBestMove *moves.Move, scoreFunc func(*board.GameState, *moves.Move) int) {
	//slices are passed by reference, so no need to return

	var moveScores []int
	for _, i := range moveList {
		if i == prevBestMove {
			moveScores = append(moveScores, inf)  //we want to evaluate the best move from the last search first
		} else {
			moveScores = append(moveScores, scoreFunc(state, i))
		}
	}

	quickSort(moveList, moveScores, 0, len(moveList) - 1)
}