package engine


import "math/rand"


func CalculateMove(currentPosition [8][8]int) [8][8]int {
	//return a new board ([8][8]int) with the engine's move
	newBoard := test(currentPosition)

	return newBoard
}


func test(board [8][8]int) [8][8]int {
	//place random piece at random pos
	x := rand.Intn(8)
	y := rand.Intn(8)
	piece := rand.Intn(11) + 1

	board[x][y] = piece

	return board
}