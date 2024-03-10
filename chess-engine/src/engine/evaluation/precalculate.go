package evaluation


const centerRank int = 4
const centerFile int = 4

var centerManhattanDist [64]int
var centerChebyshevDist [64]int


func abs(a int) int {
	if a > 0 {
		return a
	} else {
		return -a
	}
}


func PrecalculateDists() {
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			pos := file * 8 + rank

			mDist := abs(file - centerFile) + abs(rank - centerRank)
			cDist := max(abs(file - centerFile), abs(rank - centerRank))

			centerManhattanDist[pos] = mDist
			centerChebyshevDist[pos] = cDist
		}
	}
}