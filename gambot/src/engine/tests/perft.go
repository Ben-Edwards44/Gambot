package tests


import (
	"fmt"
	"time"
	"strconv"
	"gambot/src/engine/moves"
	"gambot/src/engine/board"
)


const runAllTests bool = false


var fileNames [8]string = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}
var promotionValues [7]string = [7]string{"", "", "n", "b", "r", "", "q"}  //empty strings are for pieces that you cannot promote to (or just no promotion at all)

var whitePieces [6]rune = [6]rune{'P', 'N', 'B', 'R', 'K', 'Q'}
var blackPieces [6]rune = [6]rune{'p', 'n', 'b', 'r', 'k', 'q'}


func getFenStr(state *board.GameState) string {
	//return the FEN string of the board for use when debugging
	fen := ""
	
	for i := 0; i < 8; i++ {
		rank := state.Board[i * 8 : (i + 1) * 8]

		emptyRun := 0
		for _, pieceVal := range rank {
			if pieceVal == 0 {
				emptyRun++
				continue
			}

			if emptyRun > 0 {
				//no longer empty run
				fen += strconv.Itoa(emptyRun)
				emptyRun = 0
			}

			var char rune
			if pieceVal < 7 {
				char = whitePieces[pieceVal - 1]
			} else {
				char = blackPieces[pieceVal - 7]
			}

			fen += string(char)
		}

		if emptyRun > 0 {fen += strconv.Itoa(emptyRun)}
		if i != 7 {fen += "/"}
	}

	return fen
}


func getMoveStr(move moves.Move) string {
	startRank := strconv.Itoa(8 - move.StartX)
	startFile := fileNames[move.StartY]
	endRank := strconv.Itoa(8 - move.EndX)
	endFile := fileNames[move.EndY]

	pVal := move.PromotionValue
	if pVal > 6 {pVal -= 6}

	pStr := promotionValues[pVal]

	return startFile + startRank + endFile + endRank + pStr
}


func runTests(state *board.GameState, prevMove *moves.Move) {
	if !runAllTests {return}

	defer func() {
		r := recover()

		if r != nil {
			//print helpful debugging info
			fmt.Print("Current position: ")
			fmt.Println(getFenStr(state))

			fmt.Print("Last move: ")
			fmt.Println(prevMove)

			panic(r)
		}
	}()

	testPieceInxs(state)
	testBitboards(state)
}


func bulkCount(position *board.GameState, depth int) int {	
	moveList := moves.GenerateAllMoves(position, false)

	if depth == 1 {return len(moveList)}

	total := 0
	for _, i := range moveList {
		moves.MakeMove(position, i)

		runTests(position, i)

		total += bulkCount(position, depth - 1)
		moves.UnMakeLastMove(position)
	}

	return total
}


func dividePerft(stateObj *board.GameState, maxDepth int) int {
	initMoves := moves.GenerateAllMoves(stateObj, false)

	total := 0
	for _, i := range initMoves {
		str := getMoveStr(*i)

		moves.MakeMove(stateObj, i)

		runTests(stateObj, i)

		current := 1
		if maxDepth > 1 {
			current = bulkCount(stateObj, maxDepth - 1)
		} 

		total += current

		moves.UnMakeLastMove(stateObj)

		fmt.Print(str + ": ")
		fmt.Println(current)
	}

	return total
}


func Perft(stateObj *board.GameState, maxDepth int) {
	start := time.Now()

	if runAllTests {fmt.Println("Running tests: TRUE")}

	runTests(stateObj, &moves.Move{})
	
	nodes := dividePerft(stateObj, maxDepth)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Print("Time taken: ")
	fmt.Println(elapsed)

	fmt.Print("Nodes searched: ")
	fmt.Println(nodes)
}