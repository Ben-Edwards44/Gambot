package engine


type GameState struct {
	Board [64]int

	WhiteToMove bool

	WhiteCanCastle bool
	BlackCanCastle bool

	PrevPawnDouble [2]int
}