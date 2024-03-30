package engine


import (
	"gambot/src/engine/moves"
	"gambot/src/engine/board"
	"gambot/src/engine/search"
	"gambot/src/engine/evaluation"
)


func Init() {
	//to be called at start of every new game
	moves.PrecalculateEdgeDists()
	board.PrecalculateZobristNums()
	evaluation.PrecalculateDists()

	search.NewTT()
}