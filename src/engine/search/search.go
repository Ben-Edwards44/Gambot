package search


import "chess-engine/src/engine/moves"


func GetBestMove(state *moves.GameState) moves.Move {
	moveList := moves.GenerateAllMoves(state)

	//for now, just make the first move
	return moveList[0]
}