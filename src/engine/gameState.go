package engine


type GameState struct {
	Board [64]int

	WhiteToMove bool

	WhiteKingCastle bool
	WhiteQueenCastle bool

	BlackKingCastle bool
	BlackQueenCastle bool

	PrevPawnDouble [2]int
}