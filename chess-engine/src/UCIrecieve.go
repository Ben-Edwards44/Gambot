package src

import (
	"bufio"
	"chess-engine/src/engine/board"
	"chess-engine/src/engine/moves"
	"os"
	"strconv"
)


func findInx(list []string, value string) int {
	for i, x := range list {
		if x == value {
			return i
		}
	}

	return -1
}


func splitArgs(cmd string) []string {
	currentArg := ""

	var args []string
	for _, i := range cmd {
		//check for quotes for commands like: position fen "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" moves e2e4
		if len(currentArg) > 0 && currentArg[0] == '"' {
			if i == '"' {
				args = append(args, currentArg[1:])
				currentArg = ""
			} else {
				currentArg += string(i)
			}
		} else {
			if i == ' ' {
				args = append(args, currentArg)
				currentArg = ""
			} else {
				currentArg += string(i)
			}
		}
	}

	args = append(args, currentArg)

	return args
}


func getXCoord(rank byte) int {
	coord, err := strconv.Atoi(string(rank))
	if err != nil {panic(err)}

	return 8 - coord
}


func getYCoord(file byte) int {
	rFile := rune(file)

	for i, x := range files {
		if x == rFile {return i}
	}

	panic("Invalid file")
}


func parseMove(stateObj *board.GameState, move string) moves.Move {
	startX := getXCoord(move[1])
	startY := getYCoord(move[0])
	endX := getXCoord(move[3])
	endY := getYCoord(move[2])

	pieceVal := stateObj.Board[startX * 8 + startY]
	endVal := stateObj.Board[endX * 8 + endY]

	doublePawn := (pieceVal == 1 || pieceVal == 7) && (startX - endX == 2 || startX - endX == -2)
	ep := (pieceVal == 1 || pieceVal == 7) && (startY != endY) && (endVal == 0)
	kingCastle := (pieceVal == 5 || pieceVal == 11) && (endY - startY == 2)
	queenCastle := (pieceVal == 5 || pieceVal == 11) && (startY - endY == 2)

	promotionVal := 0
	if len(move) > 4 {
		//NOTE: the promoted piece is always lowercase
		for i, x := range blackPieces {
			if x == rune(move[4]) {
				if pieceVal > 6 {
					promotionVal = i + 7
				} else {
					promotionVal = i + 1
				}

				break
			}
		}
	}

	moveObj := moves.Move{StartX: startX, StartY: startY, EndX: endX, EndY: endY, PieceValue: pieceVal, DoublePawnMove: doublePawn, EnPassant: ep, KingCastle: kingCastle, QueenCastle: queenCastle, PromotionValue: promotionVal}

	return moveObj
}


func posCmd(splitCmd []string) {
	var fen string
	if splitCmd[1] == "fen" {
		fen = splitCmd[2]
	} else if splitCmd[1] == "startpos" {
		fen = startFen
	}

	stateObj := parseFen(fen)

	//play the moves
	inx := findInx(splitCmd, "moves")
	if inx != -1 {
		for i := inx + 1; i < len(splitCmd); i++ {
			moveObj := parseMove(&stateObj, splitCmd[i])

			moves.MakeMove(&stateObj, moveObj)
		}
	}

	chessEngine.setPosition(stateObj)
}


func getIntArg(cmd []string, key string) int {
	for i, x := range cmd {
		if x == key {
			val := cmd[i + 1]
			num, err := strconv.Atoi(val)

			if err != nil {panic(err)}

			return num
		}
	}

	return -1
}


func goCmd(splitCmd []string) {
	var isPerft bool
	for _, i := range splitCmd {
		if i == "perft" {
			isPerft = true
			break
		}
	}

	if isPerft {
		//engine must perform perft
		depth, err := strconv.Atoi(splitCmd[2])

		if err != nil {panic(err)}

		chessEngine.runPerft(depth)
	} else {
		//engine must search for best move

		moveTime := getIntArg(splitCmd, "movetime")
		if moveTime == -1 {
			//the engine has to calculate the move time itself
			wClock := getIntArg(splitCmd, "wtime")
			bClock := getIntArg(splitCmd, "btime")
			wInc := getIntArg(splitCmd, "winc")
			bInc := getIntArg(splitCmd, "binc")
	
			chessEngine.calcMoveTime(wClock, bClock, wInc, bInc)
		} else {
			chessEngine.updateMoveTime(moveTime)
		}

		bestMove := chessEngine.runBestMove()

		sendBestMove(bestMove)
	}
}


func interpretCmd(cmd string) {
	if cmd == "" {return}

	splitted := splitArgs(cmd)

	switch splitted[0] {
	case "uci":
		uciOk()
	case "isready":
		sendStr("readyok")
	case "ucinewgame":
		chessEngine.newGame()
	case "position":
		posCmd(splitted)
	case "go":
		goCmd(splitted)
	case "stop":
		stop = true
	case "quit":
		stop = true
	default:
		panic("Unrecognised command: " + cmd)
	}
}


func recieveCmd() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()  //read stdin

	cmd := scanner.Text()

	if scanner.Err() != nil {panic(scanner.Err())}

	interpretCmd(cmd)
}