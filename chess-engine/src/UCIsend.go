package src


import (
	"fmt"
	"strconv"
	"chess-engine/src/engine/moves"
)


const engineName string = "chess-engine v1"
const engineAuthor string = "Ben Edwards"


func sendStr(text string) {
	fmt.Println(text)
}


func uciOk() {
	sendStr("id name " + engineName)
	sendStr("id author " + engineAuthor)
	sendStr("uciok")
}


func convertMove(move *moves.Move) string {
	//convert a move obj to a string like e2e4
	if move.PieceValue == 0 {
		return "0000"  //null move
	}

	startFile := string(files[move.StartY])
	startRank := strconv.Itoa(8 - move.StartX)
	endFile := string(files[move.EndY])
	endRank := strconv.Itoa(8 - move.EndX)

	return startFile + startRank + endFile + endRank
}


func sendBestMove(bestMove *moves.Move) {
	//TODO: add ponder move
	moveStr := convertMove(bestMove)
	
	sendStr("bestmove " + moveStr)
}