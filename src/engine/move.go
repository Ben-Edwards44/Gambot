package engine


import (
	"chess-engine/src/engine/moves"
)


func CalculateMove(stateObj GameState) GameState {
	//TODO: return a new game state with the engine's move
	return stateObj
}


func GetLegalMoves(stateObj GameState, x int, y int) [][2]int {
	var legalMoves []moves.Move
	
	moves.GetPieceMoves(stateObj.Board, x, y, stateObj.PrevPawnDouble, &legalMoves)

	//convert move structs to list of coords
	var coords [][2]int
	for _, i := range legalMoves {
		coord := [2]int{i.EndX, i.EndY}
		coords = append(coords, coord)
	}
	
	return coords
}