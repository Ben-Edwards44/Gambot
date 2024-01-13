package moves


//array matching piece values to their appropriate move functions
var moveFunctions [7]func([64]int, int, int, int) []move = [7]func([64]int, int, int, int) []move{emptyMove, emptyMove, emptyMove, emptyMove, rookMoves, emptyMove, emptyMove}


var dists [8 * 8 * 4]int
func InitPrecalculate(edgeDists [256]int) {
	dists = edgeDists
}


func emptyMove(board [64]int, x int, y int, pieceValue int) []move {
	//no legal moves
	return []move{}
}


func canMove(board [64]int, x int, y int, pieceValue int) (bool, bool) {
	inx := x * 8 + y
	sqValue := board[inx]

	if sqValue == 0 {
		return true, false
	} else {
		return (sqValue > 6) != (pieceValue > 6), true
	}
}


func rookMoves(board [64]int, x int, y int, pieceValue int) []move {
	dirInx := x * 32 + y * 4

	var moves []move
	for dir := 0; dir < 4; dir++ {
		edgeDist := dists[dirInx + dir]

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x
			newY := y

			if dir == 0 {
				newX = x - offset
			} else if dir == 1 {
				newX = x + offset
			} else if dir == 2 {
				newY = y - offset
			} else if dir == 3 {
				newY = y + offset
			}

			goodSq, capture := canMove(board, newX, newY, pieceValue) 

			if goodSq {
				m := move{x, y, newX, newY, pieceValue}
				moves = append(moves, m)
			}
			if capture {
				break
			}
		}
	}

	return moves
}


func GetPieceMoves(board [64]int, x int, y int) []move {
	pieceValue := board[x * 8 + y]
	if pieceValue != 0 {
		var inx int

		if pieceValue < 7 {
			inx = pieceValue
		} else {
			inx = pieceValue - 6
		}

		moveFunc := moveFunctions[inx]
		moves := moveFunc(board, x, y, pieceValue)

		return moves
	} else {
		panic("piece value not 0")
	}
}