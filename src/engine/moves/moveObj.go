package moves


type move struct {
	startX int
	startY int

	EndX int
	EndY int

	pieceValue int

	//TODO: add flags like en passant, promotions etc.
}