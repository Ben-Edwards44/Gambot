package moves


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//TODO: add flags like en passant, promotions etc.
	PawnDoubleMove bool
}