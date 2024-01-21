package moves


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//special moves
	EnPassant bool
	KingCastle bool
	QueenCastle bool
	promotionValue int
}