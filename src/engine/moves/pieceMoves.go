package moves


//array matching piece values to their appropriate move functions
var moveFunctions [7]func([64]int, int, int, int) []move = [7]func([64]int, int, int, int) []move{emptyMove, pawnMoves, knightMoves, bishopMoves, rookMoves, kingMoves, queenMoves}

//array matching distance index to their x, y multipliers
var xMults [8]int = [8]int{-1, 1, 0, 0, -1, -1, 1, 1}
var yMults [8]int = [8]int{0, 0, -1, 1, -1, 1, -1, 1}


var dists [512]int
func InitPrecalculate(edgeDists [512]int) {
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
	dirInx := x * 64 + y * 8

	var moves []move
	for dir := 0; dir < 4; dir++ {
		//get from precalculated
		edgeDist := dists[dirInx + dir]

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir]
			newY := y + offset * yMults[dir]

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


func bishopMoves(board [64]int, x int, y int, pieceValue int) []move {
	dirInx := x * 64 + y * 8

	var moves []move
	for dir := 0; dir < 4; dir++ {
		//get from precalculated (+4 since we are looking at diagonal)
		edgeDist := dists[dirInx + dir + 4]

		for offset := 1; offset <= edgeDist; offset++ {
			newX := x + offset * xMults[dir + 4]
			newY := y + offset * yMults[dir + 4]

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


func queenMoves(board [64]int, x int, y int, pieceValue int) []move {
	rook := rookMoves(board, x, y, pieceValue)
	bishop := bishopMoves(board, x, y, pieceValue)
	allMoves := append(rook, bishop...)

	return allMoves
}


func kingMoves(board [64]int, x int, y int, pieceValue int) []move {
	edgeInx := x * 64 + y * 8

	var moves []move
	for dir := 0; dir < 8; dir++ {
		edgeDist := dists[edgeInx + dir]

		if edgeDist > 0 {
			newX := x + xMults[dir]
			newY := y + yMults[dir]
			
			good, _ := canMove(board, newX, newY, pieceValue)

			if good {
				m := move{x, y, newX, newY, pieceValue}
				moves = append(moves, m)
			}
		}
	}

	return moves
}


func knightMoves(board [64]int, x int, y int, pieceValue int) []move {
	var moves []move
	for xStep := 1; xStep < 3; xStep++ {
		for xMult := -1; xMult < 2; xMult += 2 {
			newX := x + xStep * xMult

			if newX < 0 || newX > 7 {break}

			//xStep 1 => yStep 2, xStep 2 => yStep 1
			yStep := 3 - xStep
			for yMult := -1; yMult < 2; yMult += 2 {
				newY := y + yStep * yMult

				if newY < 0 || newY > 7 {break}

				good, _ := canMove(board, newX, newY, pieceValue)

				if good {
					m := move{x, y, newX, newY, pieceValue}
					moves = append(moves, m)
				}
			}
		}
	}

	return moves
}


func pawnMoves(board [64]int, x int, y int, pieceValue int) []move {
	if x == 0 || x == 7 {
		//on back rank
		return []move{}
	}

	isWhite := pieceValue < 7
	onStart := (isWhite && x == 6) || (!isWhite && x == 1)

	xMult := 1
	if isWhite {
		xMult = -1
	}
	
	maxStep := 1
	if onStart {
		maxStep = 2
	}

	//normal moves - no capture
	var moves []move
	for i := 1; i <= maxStep; i++ {
		newX := x + i * xMult

		good, capture := canMove(board, newX, y, pieceValue)
		if good && !capture {
			m := move{x, y, newX, y, pieceValue}
			moves = append(moves, m)
		}
	}

	//capture moves
	newX := x + xMult
	for i := -1; i < 2; i += 2 {
		newY := y + i

		if 0 <= newY && newY < 8 {
			good, capture := canMove(board, newX, newY, pieceValue)

			if good && capture {
				m := move{x, y, newX, newY, pieceValue}
				moves = append(moves, m)
			}
		}
	}
	
	
	return moves
}


func GetPieceMoves(board [64]int, x int, y int) []move {
	pieceValue := board[x * 8 + y]
	if pieceValue != 0 {
		var inx int

		//accounts for white/black
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