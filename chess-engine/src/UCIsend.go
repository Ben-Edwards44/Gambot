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


func sendBestMove(bestMove moves.Move) {
	//TODO: add ponder move

	startFile := string(files[bestMove.StartY])
	startRank := strconv.Itoa(8 - bestMove.StartX)
	endFile := string(files[bestMove.EndY])
	endRank := strconv.Itoa(8 - bestMove.EndX)

	sendStr("bestmove " + startFile + startRank + endFile + endRank)
}