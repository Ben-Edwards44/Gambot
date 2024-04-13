package moves


import "strconv"


var pieces [6]rune = [6]rune{'p', 'n', 'b', 'r', 'k', 'q'}
var files [8]rune = [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}


type Move struct {
	StartX int
	StartY int

	EndX int
	EndY int

	PieceValue int

	//special moves
	DoublePawnMove bool
	EnPassant bool
	KingCastle bool
	QueenCastle bool
	PromotionValue int
}


func (move *Move) MoveStr() string {
	//return a string representation like e2e4 of a7a8q
	if move == nil || move.PieceValue == 0 {
		return "0000"  //null move
	}
	
	startFile := string(files[move.StartY])
	startRank := strconv.Itoa(8 - move.StartX)
	endFile := string(files[move.EndY])
	endRank := strconv.Itoa(8 - move.EndX)
	
	promotion := ""
	if move.PromotionValue != 0 {
		inx := move.PromotionValue - 1
		if move.PromotionValue > 6 {inx -= 6}
	
		promotion = string(pieces[inx])
	}
	
	return startFile + startRank + endFile + endRank + promotion
}