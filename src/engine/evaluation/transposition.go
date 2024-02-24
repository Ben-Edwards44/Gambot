package evaluation


import "chess-engine/src/engine/moves"


const PvNode int = 0
const AllNode int = 1
const CutNode int = 2

const ttEnabled bool = true


var ttEntries map[uint64]ttEntry = make(map[uint64]ttEntry)  //TODO: use array instead of map. Use inx = zobHash % length and then check that hash == actual hash we want


type ttEntry struct {
	zobHash uint64
	depthSearched int
	eval int
	nodeType int
	bestMove moves.Move  //for depth 1 (in terms of iterative deepening lookups)
	//TODO: age (so we know when to clear)
}


func LookupEval(zobHash uint64, currentDepth int, alpha int, beta int) (bool, int) {
	if !ttEnabled {return false, 0}

	entry, exists := ttEntries[zobHash]

	if !exists {return false, 0}  //lookup failed

	if entry.depthSearched >= currentDepth {
		if entry.nodeType == PvNode {
			//we have stored the exact evaluation, so no problem
			return true, entry.eval
		} else if entry.nodeType == AllNode {
			//node is an upper bound
			//TODO: corrent mate scores
			if entry.eval <= alpha {return true, entry.eval}
		} else if entry.nodeType == CutNode {
			//node is a lower bound
			//TODO: correct mate scores
			if entry.eval >= beta {return true, entry.eval}
		}
	} 

	return false, 0  //lookup failed
}


func LookupMove(zobHash uint64) moves.Move {
	//This is for when the position we are currently searching is in the transposition table
	return ttEntries[zobHash].bestMove
}


func StoreEntry(zobHash uint64, searchDepth int, eval int, nodeType int, bestMove moves.Move) {
	entry := ttEntry{zobHash: zobHash, depthSearched: searchDepth, eval: eval, nodeType: nodeType, bestMove: bestMove}

	ttEntries[zobHash] = entry
}


func ClearTT() {
	ttEntries = make(map[uint64]ttEntry)
}