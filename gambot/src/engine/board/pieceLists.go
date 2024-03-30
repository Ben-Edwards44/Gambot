package board


var PieceLists pieceList


type pieceList struct {
	pieceInxs [64]int
	WhitePieceSquares [16]int
	BlackPieceSquares [16]int
	
	WhiteKingPos int
	BlackKingPos int

	prvPieceInxs [][64]int
	prvWPieceSq [][16]int
	prvBPieceSq [][16]int
	prvWKPos []int
	prvBKPos []int
}


func (list *pieceList) setPrevVals() {
	list.prvPieceInxs = append(list.prvPieceInxs, list.pieceInxs)
	list.prvWPieceSq = append(list.prvWPieceSq, list.WhitePieceSquares)
	list.prvBPieceSq = append(list.prvBPieceSq, list.BlackPieceSquares)
	list.prvWKPos = append(list.prvWKPos, list.WhiteKingPos)
	list.prvBKPos = append(list.prvBKPos, list.BlackKingPos)
}


func (list *pieceList) restorePrev() {
	//restore the values
	list.pieceInxs = list.prvPieceInxs[len(list.prvPieceInxs) - 1]
	list.WhitePieceSquares = list.prvWPieceSq[len(list.prvWPieceSq) - 1]
	list.BlackPieceSquares = list.prvBPieceSq[len(list.prvBPieceSq) - 1]
	list.WhiteKingPos = list.prvWKPos[len(list.prvWKPos) - 1]
	list.BlackKingPos = list.prvBKPos[len(list.prvBKPos) - 1]

	//pop end of slice
	list.prvPieceInxs = list.prvPieceInxs[:len(list.prvPieceInxs) - 1]
	list.prvWPieceSq = list.prvWPieceSq[:len(list.prvWPieceSq) - 1]
	list.prvBPieceSq = list.prvBPieceSq[:len(list.prvBPieceSq) - 1]
	list.prvWKPos = list.prvWKPos[:len(list.prvWKPos) - 1]
	list.prvBKPos = list.prvBKPos[:len(list.prvBKPos) - 1]
}


func pieceCapture(captInx int, captPieceWhite bool) {
	//remove captured piece from list
	if captPieceWhite {
		PieceLists.WhitePieceSquares[captInx] = -1
	} else {
		PieceLists.BlackPieceSquares[captInx] = -1
	}
}


func Castle(rookStart int, rookEnd int, pieceVal int) {
	//move the rook's position after a castle. NOTE: king will already have been moved
	inx := PieceLists.pieceInxs[rookStart]

	if inx == -1 {panic("Piece does not exist at start square")}

	isWhite := pieceVal < 7
	if isWhite {
		PieceLists.WhitePieceSquares[inx] = rookEnd
	} else {
		PieceLists.BlackPieceSquares[inx] = rookEnd
	}

	//update the index map
	PieceLists.pieceInxs[rookEnd] = inx
	PieceLists.pieceInxs[rookStart] = -1
}


func EnPassant(captPos int, captPieceWhite bool) {
	//remove the pawn captured by en passant
	inx := PieceLists.pieceInxs[captPos]

	pieceCapture(inx, captPieceWhite)

	PieceLists.pieceInxs[captPos] = -1  //remove pawn from piece index map as well
}


func MovePiecePosition(start int, end int, pieceVal int) {
	//update the position of a piece after is has been moved
	PieceLists.setPrevVals()  //in case we unmake the move

	inx := PieceLists.pieceInxs[start]

	if inx == -1 {panic("Piece does not exist at start square")}

	isWhite := pieceVal < 7
	if isWhite {
		PieceLists.WhitePieceSquares[inx] = end
	} else {
		PieceLists.BlackPieceSquares[inx] = end
	}

	captInx := PieceLists.pieceInxs[end]
	if captInx != -1 {pieceCapture(captInx, !isWhite)}

	//update the index map
	PieceLists.pieceInxs[end] = inx
	PieceLists.pieceInxs[start] = -1

	//deal with king
	if pieceVal == 5 {
		PieceLists.WhiteKingPos = end
	} else if pieceVal == 11 {
		PieceLists.BlackKingPos = end
	}
}


func UnMoveLastPiece() {
	PieceLists.restorePrev()
}


func InitPieceLists(state *GameState) {
	//NOTE: this function is slow and should only be called on the creation of a new GameState obj
	whiteInx := 0
	blackInx := 0

	var pieceInxs [64]int
	var whitePieceSquares [16]int
	var blackPieceSquares [16]int

	var whiteKingPos int
	var blackKingPos int

	for i := 0; i < 64; i++ {
		pieceVal := state.Board[i]

		if pieceVal != 0 {
			if pieceVal < 7 {
				//white
				pieceInxs[i] = whiteInx
				whitePieceSquares[whiteInx] = i
				whiteInx++
			} else {
				//black
				pieceInxs[i] = blackInx
				blackPieceSquares[blackInx] = i
				blackInx++
			}

			if pieceVal == 5 {whiteKingPos = i}
			if pieceVal == 11 {blackKingPos = i}
		} else {
			pieceInxs[i] = -1
		}
	}

	//fill in rest of arrays
	for i := whiteInx; i < len(whitePieceSquares); i++ {
		whitePieceSquares[i] = -1
	}
	for i := blackInx; i < len(blackPieceSquares); i++ {
		blackPieceSquares[i] = -1
	}

	PieceLists = pieceList{pieceInxs: pieceInxs, WhitePieceSquares: whitePieceSquares, BlackPieceSquares: blackPieceSquares, WhiteKingPos: whiteKingPos, BlackKingPos: blackKingPos}
}