package search


import (
	"unsafe"
	"chess-engine/src/engine/moves"
)


const pvNode int = 0
const allNode int = 1
const cutNode int = 2

const ttEnabled bool = true

const ttSizeMib int = 64
const ttLen uint64 = uint64((1024 * 1024 * ttSizeMib) / int(unsafe.Sizeof(ttEntry{})))  //using uint64 means we don't have to convert later


var ttEntries [ttLen]ttEntry


type ttEntry struct {
	zobHash uint64
	depthSearched int
	eval int
	nodeType int
	bestMove *moves.Move  //for depth 1 (in terms of iterative deepening lookups)
	//TODO: age (so we know when to clear)
}


func lookupEval(zobHash uint64, currentDepth int, alpha int, beta int) (bool, int) {
	if !ttEnabled {return false, 0}

	inx := zobHash % ttLen

	entry := ttEntries[inx]

	if entry.zobHash != zobHash {return false, 0}  //lookup failed

	if entry.depthSearched >= currentDepth {
		if entry.nodeType == pvNode {
			//we have stored the exact evaluation, so no problem
			return true, entry.eval
		} else if entry.nodeType == allNode {
			//node is an upper bound
			//TODO: corrent mate scores
			if entry.eval <= alpha {return true, entry.eval}
		} else if entry.nodeType == cutNode {
			//node is a lower bound
			//TODO: correct mate scores
			if entry.eval >= beta {return true, entry.eval}
		}
	} 

	return false, 0  //lookup failed
}


func lookupMove(zobHash uint64) *moves.Move {
	//This is for when the position we are currently searching is in the transposition table
	inx := zobHash % ttLen

	return ttEntries[inx].bestMove
}


func storeEntry(zobHash uint64, searchDepth int, eval int, nodeType int, bestMove *moves.Move) {
	entry := ttEntry{zobHash: zobHash, depthSearched: searchDepth, eval: eval, nodeType: nodeType, bestMove: bestMove}
	inx := zobHash % ttLen

	ttEntries[inx] = entry
}


func NewTT() {
	ttEntries = [ttLen]ttEntry{}
}