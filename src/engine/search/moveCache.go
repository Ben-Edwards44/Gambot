package search

import (
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
	"fmt"
	"strconv"
)

//TODO: use bits to represent moves and use array instead of map
var moveCache map[string][]moves.Move
var captureCache map[string][]moves.Move


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


func getMoveList(state *board.GameState, moveChain string, onlyCaptures bool) []moves.Move {
	//get the ordered list of moves from a previous iteration of iterative deepening

	cache := moveCache
	if onlyCaptures {cache = captureCache}

	moveList, exists := cache[moveChain]

	if !exists {
		//moves were not cached, so we need to actually calculate them
		moveList = manuallyGenerateMoves(state, moveChain, onlyCaptures)
	}

	return moveList
}


func manuallyGenerateMoves(state *board.GameState, moveChain string, onlyCaptures bool) []moves.Move {
	moveList := moves.GenerateAllMoves(state, onlyCaptures)
	orderMoves(state, moveList, moves.Move{})

	appendToCache(moveChain, moveList, onlyCaptures)  //Add the move to the cache

	return moveList
}


func appendToCache(moveChain string, moveList []moves.Move, onlyCaptures bool) {
	cache := moveCache
	if onlyCaptures {cache = captureCache}

	cache[moveChain] = moveList

	if onlyCaptures {
		captureCache = cache
	} else {
		moveCache = cache
	}
}


func updateFirstMove(moveChain string, newFirstInx int, onlyCaptures bool) {
	//update the move order with the new best move at the front

	var list []moves.Move
	if onlyCaptures {
		list = captureCache[moveChain]
		fmt.Println(list)
	} else {
		list = moveCache[moveChain]
	}

	moveVal := list[newFirstInx]

	//shift other elements
	for i := newFirstInx - 1; i >= 0; i-- {
		list[i + 1] = list[i]
	}

	list[0] = moveVal  //shift the new best into first place

	if onlyCaptures {
		captureCache[moveChain] = list
	} else {
		moveCache[moveChain] = list
	}
}


func clearCaches() {
	moveCache = make(map[string][]moves.Move)
	captureCache = make(map[string][]moves.Move)
}