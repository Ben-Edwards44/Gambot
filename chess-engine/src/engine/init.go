package engine


import (
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/search"
	"chess-engine/src/engine/evaluation"
)


func Init() {
	//to be called at start of every new game
	moves.PrecalculateEdgeDists()
	board.PrecalculateZobristNums()
	evaluation.PrecalculateDists()

	search.NewTT()
}