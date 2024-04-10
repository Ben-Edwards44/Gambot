package search


const maxLength int = 256


var repTable repetitionTable


type repetitionTable struct {
	seenHashes [maxLength]uint64
	length int
	startWhite bool
}


func (table *repetitionTable) push(hash uint64) {
	table.seenHashes[table.length] = hash
	table.length++
}


func (table *repetitionTable) pop() {
	table.length--
}


func (table *repetitionTable) seen(hash uint64, isWhite bool) bool {
	//NOTE: we only need to search position indices corresponding to the same colour as we are currently searching
	startInx := 1
	if table.startWhite == isWhite {startInx = 0}

	for i := startInx; i < table.length; i += 2 {
		if table.seenHashes[i] == hash {return true}
	}

	return false
}


func initRepTable(rootIsWhite bool) {
	repTable = repetitionTable{startWhite: rootIsWhite}
}