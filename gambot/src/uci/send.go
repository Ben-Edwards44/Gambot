package uci


import (
	"fmt"
	"strconv"
	"gambot/src/engine/moves"
)


const engineName string = "Gambot v1"
const engineAuthor string = "Ben Edwards"


func sendStr(text string) {
	fmt.Println(text)
}


func uciOk(engine *bot) {
	sendStr("id name " + engineName)
	sendStr("id author " + engineAuthor)

	sendOptions(engine)

	sendStr("uciok")
}


func sendBestMove(bestMove *moves.Move) {
	//TODO: add ponder move
	moveStr := bestMove.MoveStr()
	
	sendStr("bestmove " + moveStr)
}


func sendSpinOpt(option *spinOption) {
	str := "option name " + option.name
	str += " type spin"
	str += " default " + strconv.Itoa(option.defaultVal)
	str += " min " + strconv.Itoa(option.min)
	str += " max " + strconv.Itoa(option.max)

	sendStr(str)
}


func sendOptions(engine *bot) {
	sendStr("")
	sendSpinOpt(engine.ttSize)
}