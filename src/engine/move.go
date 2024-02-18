package engine


import (
	"chess-engine/src/engine/moves"
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/search"
)


func CheckWin(stateObj *board.GameState) string {
	legalMoves := moves.GenerateAllMoves(stateObj, false)

	if len(legalMoves) > 0 {return "not_terminal"}
	
	kingPos := stateObj.BlackPiecePos[4][0]
	if stateObj.WhiteToMove {kingPos = stateObj.WhitePiecePos[4][0]}

	pos := kingPos[0] * 8 + kingPos[1]

	//set bitboard at king's position
	var kingPosBB uint64
	kingPosBB |= 1 << pos

	inCheck := (kingPosBB & stateObj.NoKingMoveBitBoard) != 0

	if inCheck {
		if stateObj.WhiteToMove {
			return "black_win"
		} else {
			return "white_win"
		}
	} else {
		return "draw"
	}
}


func CalculateMove(stateObj *board.GameState) board.GameState {
	//TODO: use pointer for return rather than value

	move := search.GetBestMove(stateObj)

	if move.PieceValue != 0 {moves.MakeMove(stateObj, move)}  //If in checkmate, the piece value will be 0
	
	return *stateObj
}


func GetLegalMoves(stateObj *board.GameState, x int, y int) [][2]int {
	var legalMoves []moves.Move
	
	moves.GetPieceMoves(stateObj, x, y, &legalMoves, false)

	//convert move structs to list of coords
	var coords [][2]int
	for _, i := range legalMoves {
		coord := [2]int{i.EndX, i.EndY}
		coords = append(coords, coord)
	}
	
	return coords
}