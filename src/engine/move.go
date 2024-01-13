package engine


import (
	"chess-engine/src/engine/moves"
)


func CalculateMove(board [64]int) [64]int {
	//TODO: return a new board ([64]int) with the engine's move
	return board
}


func GetLegalMoves(board [64]int, x int, y int, prevMove moves.Move) [][2]int {
	moves := moves.GetPieceMoves(board, x, y, prevMove)

	//convert move structs to list of coords
	var coords [][2]int
	for _, i := range moves {
		coord := [2]int{i.EndX, i.EndY}
		coords = append(coords, coord)
	}
	
	return coords
}