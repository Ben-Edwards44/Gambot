package uci


import (
	"strconv"
	"strings"
	"gambot/src/engine/board"
	"gambot/src/engine/moves"
)


const startFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var whitePieces [6]rune = [6]rune{'P', 'N', 'B', 'R', 'K', 'Q'}
var blackPieces [6]rune = [6]rune{'p', 'n', 'b', 'r', 'k', 'q'}

var files [8]rune = [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}


func getPieceValue (piece rune) int {
	for i, x := range whitePieces {
		if x == piece {return i + 1}
	}

	for i, x := range blackPieces {
		if x == piece {return i + 7}
	}

	return 0
}


func parseBoard(fenBoard string) [64]int {
	var inx int
	var board [64]int
	for _, char := range fenBoard {
		if char != '/' {
			pieceVal := getPieceValue(char)

			if pieceVal == 0 {
				num, err := strconv.Atoi(string(char))
				if err != nil {panic(err)}

				inx += num
			} else {
				board[inx] = pieceVal
				inx++
			}
		}
	}

	return board
}


func parseEp(epTarget string, whiteToMove bool) [2]int {
	if epTarget == "-" {return [2]int{-1, -1}}

	x, err := strconv.Atoi(string(epTarget[1]))
	if err != nil {panic(err)}

	var y int
	for i, x := range files {
		if x == rune(epTarget[0]) {
			y = i
			break
		}
	}

	if whiteToMove {
		return [2]int{x + 1, y}
	} else {
		return [2]int{x - 1, y}
	}
}


func parseCastle(castle string) uint8 {
	var castleRights uint8
	for _, i := range castle {
		if i == 'K' {
			castleRights |= board.WkCastle
		} else if i == 'Q' {
			castleRights |= board.WqCastle
		} else if i == 'k' {
			castleRights |= board.BkCastle
		} else if i == 'q' {
			castleRights |= board.BqCastle
		}
	}

	return castleRights
}


func parseFen(fen string) board.GameState {
	splitted := strings.Split(fen, " ")
	
	board := parseBoard(splitted[0])
	whiteToMove := splitted[1] == "w"
	castleRights := parseCastle(splitted[2])
	prevPawnDouble := parseEp(splitted[3], whiteToMove)

	//TODO: parse half moves and full moves

	stateObj := moves.CreateGameState(board, whiteToMove, castleRights, prevPawnDouble)

	return stateObj
}