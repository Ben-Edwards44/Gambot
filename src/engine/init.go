package engine


import (
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
)


func findDistToEdge() [512]int {
	//1d array is faster than 2d
	var dists [8 * 8 * 8]int

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			up := x
			down := 7 - x
			left := y
			right := 7 - y

			inx := x * 64 + y * 8

			dists[inx] = up
			dists[inx + 1] = down
			dists[inx + 2] = left
			dists[inx + 3] = right
			dists[inx + 4] = min(up, left)
			dists[inx + 5] = min(up, right)
			dists[inx + 6] = min(down, left)
			dists[inx + 7] = min(down, right)
		}
	}

	return dists
}


func PrecomputeValues() {
	//to be called at start of execution
	edgeDists := findDistToEdge()
	moves.InitPrecalculate(edgeDists)

	board.PrecalculateZobristNums()
}