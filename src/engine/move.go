package engine


import (
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/search"
)


func CalculateMove(stateObj *moves.GameState) moves.GameState {
	//TODO: use pointer for return rather than value

	move := search.GetBestMove(stateObj)

	moves.MakeMove(stateObj, move)
	
	return *stateObj
}


func GetLegalMoves(stateObj *moves.GameState, x int, y int) [][2]int {
	var legalMoves []moves.Move
	
	moves.GetPieceMoves(stateObj, x, y, &legalMoves)

	//convert move structs to list of coords
	var coords [][2]int
	for _, i := range legalMoves {
		coord := [2]int{i.EndX, i.EndY}
		coords = append(coords, coord)
	}
	
	return coords
}