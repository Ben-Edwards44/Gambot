package uci


import (
	"gambot/src/engine"
	"gambot/src/engine/board"
	"gambot/src/engine/moves"
	"gambot/src/engine/tests"
	"gambot/src/engine/evaluation"
)


const maxMoveTime int = 5 * 60 * 1000  //5 mins
const moveFract int = 40  //1 / 40 of remaining time
const incMult float32 = 0.8  //80% of the increment


var chessEngine bot


type bot struct {
	currentPos board.GameState
	moveTime int
	ttSize *spinOption
}


func (b *bot) setPosition(state board.GameState) {
	b.currentPos = state
}


func (b *bot) calcMoveTime(wClock int, bClock int, wInc int, bInc int) {
	//calculate and update move time (in ms)

	clock := bClock
	inc := bInc
	if b.currentPos.WhiteToMove {
		clock = wClock
		inc = wInc
	}

	moveTime := clock / moveFract

	if inc != -1 && clock > inc {
		incAdd := float32(inc) + incMult
		moveTime += int(incAdd)
	}

	if moveTime > maxMoveTime {moveTime = maxMoveTime}

	b.moveTime = moveTime
}


func (b *bot) updateMoveTime(moveTime int) {
	//move time is given by UCI, so no need to calculate it
	b.moveTime = moveTime
}


func (b *bot) runPerft(depth int) {
	//run perft - assumes the position has been updated
	tests.Perft(&b.currentPos, depth)
}


func (b *bot) runBestMove() *moves.Move {
	//calculate the best move - assumes position and move time have been updated
	bestMove := engine.CalculateMove(&b.currentPos, b.moveTime)

	return bestMove
}


func (b *bot) runEval() int {
	//evaluate the current position - assumes position has been set
	eval := evaluation.Eval(&b.currentPos, b.currentPos.WhiteToMove)

	return eval
}


func (b *bot) newGame() {
	//start of new game
	engine.Init(b.ttSize.setVal)
}


func InitEngine() {
	//TODO: init engine when it is supposed to (according to UCI)
	ttSize := createTTOpt()

	chessEngine = bot{ttSize: ttSize}

	engine.Init(ttSize.defaultVal)
}