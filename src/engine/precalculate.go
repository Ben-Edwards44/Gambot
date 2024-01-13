package engine


import "chess-engine/src/engine/moves"


func findDistToEdge() [256]int {
	//1d array is faster than 2d
	var dists [8 * 8 * 4]int

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			up := x
			down := 7 - x
			left := y
			right := 7 - y

			inx := x * 32 + y * 4

			dists[inx] = up
			dists[inx + 1] = down
			dists[inx + 2] = left
			dists[inx + 3] = right
		}
	}

	return dists
}


func PrecomputeValues() {
	//to be called at start of execution
	edgeDists := findDistToEdge()
	moves.InitPrecalculate(edgeDists)
}