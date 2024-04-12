package search


import "gambot/src/engine/moves"


const maxLength int = 256


var repTable searchRepTable
var playedReps playedRepTable


type searchRepTable struct {
	seenHashes [maxLength]uint64
	length int
}


type playedRepTable struct {
	//this should not be used within the search
	playedHashes map[uint64]bool
	repetitions searchRepTable
}


func (table *searchRepTable) push(hash uint64) {	
	table.seenHashes[table.length] = hash
	table.length++
}


func (table *searchRepTable) pop() {
	table.length--
}


func (table *searchRepTable) seen(hash uint64) bool {
	for i := 0; i < table.length; i ++ {
		if table.seenHashes[i] == hash {return true}
	}

	return false
}


func initRepTable(rootHash uint64) {
	repTable = searchRepTable{}

	//add any actually played repetitions to the table so that these are scored as draws
	repTable.length = playedReps.repetitions.length
	repTable.seenHashes = playedReps.repetitions.seenHashes

	if !repTable.seen(rootHash) {repTable.push(rootHash)}  //add the root hash
}


func ResetRepetitions(initialHash uint64) {
	//to be called when a new position command is given
	playedMap := make(map[uint64]bool)
	playedMap[initialHash] = true

	playedReps = playedRepTable{playedHashes: playedMap}
}


func AddPlayedRepetition(playedMove *moves.Move, newHash uint64, whiteToMove bool, prvBoard [64]int) {
	//see if we need to add a played repetition. This is done when a position fen ... moves ... command is given
	isRep := playedReps.playedHashes[newHash]

	if isRep {
		if !playedReps.repetitions.seen(newHash) {
			//add the repetition to the slice
			playedReps.repetitions.push(newHash)
		}
	} else {
		playedReps.playedHashes[newHash] = true
	}

	//If a move is a capture or castle, the positions before it are unreachable.
	//To save time during the search, we can just remove these now unreachable positions.
	captVal := prvBoard[playedMove.EndX * 8 + playedMove.EndY]
	if captVal != 0 || playedMove.KingCastle || playedMove.QueenCastle {
		playedReps.repetitions = searchRepTable{}
	}
}