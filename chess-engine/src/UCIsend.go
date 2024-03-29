package src


import (
	"fmt"
	"chess-engine/src/engine/moves"
)


const engineName string = "Gambot v1"
const engineAuthor string = "Ben Edwards"


func sendStr(text string) {
	fmt.Println(text)
}


func uciOk() {
	sendStr("id name " + engineName)
	sendStr("id author " + engineAuthor)
	sendStr("uciok")
}


func sendBestMove(bestMove *moves.Move) {
	//TODO: add ponder move
	moveStr := bestMove.MoveStr()
	
	sendStr("bestmove " + moveStr)
}