package engine


import (
	"math/rand"
	"chess-engine/src/engine/moves"
)


func CalculateMove(currentPosition [8][8]int) [8][8]int {
	//return a new board ([8][8]int) with the engine's move
	newBoard := test(currentPosition)

	return newBoard
}


func GetLegalMoves(board [64]int, x int, y int) [][2]int {
	moves := moves.GetPieceMoves(board, x, y)

	//convert move structs to list of coords
	var coords [][2]int
	for _, i := range moves {
		coord := [2]int{i.EndX, i.EndY}
		coords = append(coords, coord)
	}
	
	return coords
}


func test(board [8][8]int) [8][8]int {
	//place random piece at random pos
	x := rand.Intn(8)
	y := rand.Intn(8)
	piece := rand.Intn(11) + 1

	board[x][y] = piece

	return board
}