package moves


//array matching piece values to their appropriate move functions
var moveFunctions [7]func([64]int, int, int, int) [][2]int = [7]func([64]int, int, int, int) [][2]int{emptyMove, emptyMove, emptyMove, emptyMove, rookMoves, emptyMove, emptyMove}


var dists [8 * 8 * 4]int
func InitPrecalculate(edgeDists [256]int) {
	dists = edgeDists
}


func emptyMove(board [64]int, x int, y int, pieceValue int) [][2]int {
	//for testing
	moves := [][2]int{{0, 0}}
	return moves
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


func rookMoves(board [64]int, x int, y int, pieceValue int) [][2]int {
	dirInx := x * 32 + y * 4

	var moves [][2]int
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
				move := [2]int{newX, newY}
				moves = append(moves, move)
			}
			if capture {
				break
			}
		}
	}

	return moves
}


func GetPieceMoves(board [64]int, x int, y int) [][2]int {
	pieceValue := board[x * 8 + y]

	if pieceValue != 0 {
		moveFunc := moveFunctions[pieceValue]
		moves := moveFunc(board, x, y, pieceValue)

		return moves
	} else {
		panic("piece value not 0")
	}
}