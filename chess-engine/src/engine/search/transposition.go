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


var searchTable ttTable
var qSearchTable ttTable


type ttEntry struct {
	zobHash uint64
	depthSearched int
	eval int
	nodeType int
	bestMove *moves.Move
}


type ttTable struct {
	entries [ttLen]ttEntry
}


func (table *ttTable) lookupEval(zobHash uint64, currentDepth int, alpha int, beta int) (bool, int) {
	if !ttEnabled {return false, 0}

	inx := zobHash % ttLen

	entry := &table.entries[inx]

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


func (table *ttTable) lookupMove(zobHash uint64) *moves.Move {
	//This is so that we search the best move from the previous depth first
	inx := zobHash % ttLen
	entry := &table.entries[inx]

	if entry.zobHash == zobHash {
		return entry.bestMove
	} else {
		return &moves.Move{}
	}
}


func (table *ttTable) storeEntry(zobHash uint64, searchDepth int, eval int, nodeType int, bestMove *moves.Move) {
	entry := ttEntry{zobHash: zobHash, depthSearched: searchDepth, eval: eval, nodeType: nodeType, bestMove: bestMove}
	inx := zobHash % ttLen

	table.entries[inx] = entry
}


func NewTT() {
	searchTable = ttTable{}
	qSearchTable = ttTable{}
}