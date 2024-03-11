package evaluation


const centerRank int = 4
const centerFile int = 4

var centerManhattanDist [64]int
var squareManhattanDists [64 * 64]int


func abs(a int) int {
	if a > 0 {
		return a
	} else {
		return -a
	}
}


func manhattanDist(rank1 int, file1 int, rank2 int, file2 int) int {
	return abs(file1 - file2) + abs(rank1 - rank2)
}


func centDists() {
	//compute the manhattan distances between the center and every square on the board
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			pos := file * 8 + rank
			mDist := manhattanDist(rank, file, centerRank, centerFile)

			centerManhattanDist[pos] = mDist
		}
	}
}


func squareDists() {
	//compute the manhattan distances between every square on the board
	for pos1 := 0; pos1 < 64; pos1++ {
		for pos2 := 0; pos2 < 64; pos2++ {
			mDist := manhattanDist(int(pos1 / 8), pos1 % 8, int(pos2 / 8), pos2 % 8)

			squareManhattanDists[pos1 * 64 + pos2] = mDist
			squareManhattanDists[pos2 * 64 + pos1] = mDist
		}
	}
}


func PrecalculateDists() {
	centDists()
	squareDists()
}