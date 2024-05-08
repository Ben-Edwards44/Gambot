package tests


import (
	"fmt"
	"gambot/src/engine/board"
)


func inList(square int, list [16]int) bool {
	for _, i := range list {
		if i == square {return true}
	}

	return false
}


func testPieceInxs(state *board.GameState) {
	//test if all of the pieces on the board are on the piece list map
	for square := 0; square < 64; square++ {
		pieceVal := state.Board[square]

		isOnWList := inList(square, board.PieceLists.WhitePieceSquares)
		isOnBList := inList(square, board.PieceLists.BlackPieceSquares)

		if pieceVal == 0 {
			if isOnWList || isOnBList {panic("Piece indexes failed")}
		} else {
			if isOnWList != (pieceVal < 7) {panic("Piece indexes failed")}
			if isOnBList != (pieceVal > 6) {panic("Piece indexes failed")}
		}
	}
}


func testBitboards(state *board.GameState) {
	for inx, bb := range state.Bitboards.WPieces {
		pieceVal := inx + 1
		for square := 0; square < 64; square++ {
			onBB := bb & (1 << square) != 0
			onBoard := state.Board[square] == pieceVal

			if onBB != onBoard {
				fmt.Println(inx)
				fmt.Println(state.Bitboards.WPieces[0])
				fmt.Println(bb)
				fmt.Println(square)
				fmt.Println(state.Board[square])
				fmt.Println(pieceVal)
				panic("Bitboards do not match")
			}
		}
	}

	for inx, bb := range state.Bitboards.BPieces {
		pieceVal := inx + 7
		for square := 0; square < 64; square++ {
			onBB := bb & (1 << square) != 0
			onBoard := state.Board[square] == pieceVal

			if onBB != onBoard {
				panic("Bitboards do not match")
			}
		}
	}
}