package search


import (
	"unsafe"
	"gambot/src/engine/moves"
)


const pvNode int = 0
const allNode int = 1
const cutNode int = 2

const ttEnabled bool = true

const DefaultTTSizeMib int = 64


var searchTable ttTable


type ttEntry struct {
	zobHash uint64
	depthSearched int
	eval int
	nodeType int
	bestMove *moves.Move
}


type ttTable struct {
	length uint64
	entries []ttEntry
}


func correctMateScore(score int, plyFromRoot int) int {
	if score > mateThreshold {
		return score - plyFromRoot
	} else if score < -mateThreshold {
		return score + plyFromRoot
	} else {
		return score
	}
}


func (table *ttTable) lookupEval(zobHash uint64, currentDepth int, plyFromRoot int, alpha int, beta int) (bool, int) {
	if !ttEnabled {return false, 0}

	inx := zobHash % table.length

	entry := &table.entries[inx]

	if entry.zobHash != zobHash {return false, 0}  //lookup failed

	score := correctMateScore(entry.eval, plyFromRoot)

	if entry.depthSearched >= currentDepth {
		if entry.nodeType == pvNode {
			//we have stored the exact evaluation, so no problem
			return true, score
		} else if entry.nodeType == allNode {
			//node is an upper bound
			if entry.eval <= alpha {return true, score}
		} else if entry.nodeType == cutNode {
			//node is a lower bound
			if entry.eval >= beta {return true, score}
		}
	} 

	return false, 0  //lookup failed
}


func (table *ttTable) lookupMove(zobHash uint64) *moves.Move {
	//This is so that we search the best move from the previous depth first (no need to check for searched depth)
	inx := zobHash % table.length
	entry := &table.entries[inx]

	if entry.zobHash == zobHash {
		return entry.bestMove
	} else {
		return &moves.Move{}
	}
}


func (table *ttTable) lookupPvMove(zobHash uint64) *moves.Move {
	//Get the PV move. Need to ensure the node is a pv node
	inx := zobHash % table.length
	entry := &table.entries[inx]

	if entry.zobHash == zobHash && entry.nodeType == pvNode {
		return entry.bestMove
	} else {
		return nil
	}
}


func (table *ttTable) storeEntry(zobHash uint64, searchDepth int, plyFromRoot int, eval int, nodeType int, bestMove *moves.Move) {
	correctedScore := correctMateScore(eval, -plyFromRoot)  //the - is because we want to increase (not decrease) the magnitude of the stored score if it is a mate

	entry := ttEntry{zobHash: zobHash, depthSearched: searchDepth, eval: correctedScore, nodeType: nodeType, bestMove: bestMove}
	inx := zobHash % table.length

	table.entries[inx] = entry
}


func NewTT(sizeMib int) {
	length := uint64((1024 * 1024 * sizeMib) / int(unsafe.Sizeof(ttEntry{})))
	entries := make([]ttEntry, length)

	searchTable = ttTable{length: length, entries: entries}
}