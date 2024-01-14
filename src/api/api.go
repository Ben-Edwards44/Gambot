package api


import (
	"os"
	"strconv"
	"strings"
	"chess-engine/src/engine"
)


const FILE_PATH string = "src/api/interface.json"


func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}


func flattenBoard(board [8][8]int) [64]int {
	var flattened [64]int

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			flattened[x * 8 + y] = board[x][y]
		}
	}

	return flattened
}


func unflattenBoard(flattened [64]int) [8][8]int {
	var board [8][8]int
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			value := flattened[x * 8 + y]
			board[x][y] = value
		}
	}

	return board
}


func boardStrToList(str string) [64]int {
	//convert the string "[[1, 2, ...], [3, 4, ...], ...]" to array [[1, 2], [3, 4]]

	//remove the [[ and ]] at end
	str = str[2 : len(str) - 2]

	strLs := strings.Split(str, "], [")

	var outList [8][8]int
	for i := 0; i < 8; i++ {
		nums := strings.Split(strLs[i], ", ")

		for x := 0; x < 8; x++ {
			num, err := strconv.Atoi(nums[x])

			panicErr(err)

			outList[i][x] = num
		}
	}

	return flattenBoard(outList)
}


func coordStrToList(str string) [2]int {
	//convert a single coord string "[x, y]" to array [x, y]

	//remove [ and ]
	str = str[1 : len(str) - 1]
	nums := strings.Split(str, ", ")

	var coords [2]int
	for i, x := range nums {
		num := strToInt(x)
		coords[i] = num
	}

	return coords
}


func removeFirstLast(str string) string {
	newStr := str[1 : len(str) - 1]

	return newStr
}


func reverseString(str string) string {
	var chars []string

	for i := len(str) - 1; i >= 0; i-- {
		char := string(str[i])
		chars = append(chars, char)
	}

	reversed := strings.Join(chars, "")

	return reversed
}


func splitJson(str string) []string {
	//str in form '"..." : [[...], ...], "..." : "..."' - assumes that every k/v is a string

	var colonInxs []int

	for i, x := range str {
		if x == ':' {
			colonInxs = append(colonInxs, i)
		}
	}

	var splitted []string
	for _, inx := range colonInxs {
		//work forwards/back until 2nd double quote

		key := ""
		backInx := inx - 1
		numQuote := 0

		for backInx >= 0 && numQuote < 2 {
			key += string(str[backInx])

			if str[backInx] == '"' {
				numQuote++
			}

			backInx--
		}

		value := ""
		//add 2 becuase there is a space after colon
		forInx := inx + 2
		numQuote = 0

		for forInx < len(str) && numQuote < 2 {
			value += string(str[forInx])

			if str[forInx] == '"' {
				numQuote++
			}

			forInx++
		}

		//because the chars were added last to first (but not with key)
		key = reverseString(key)

		splitted = append(splitted, key)
		splitted = append(splitted, value)
	}

	return splitted
}


func jsonLoad(str string) map[string]string {
	//str will look like {"board" : [[...], ...], "..." : "..."} (no need for nested {})

	//remove {}
	str = removeFirstLast(str)

	kvPairs := splitJson(str)
	json := make(map[string]string)

	for i := 0; i < len(kvPairs); i += 2 {
		k := kvPairs[i]
		v := kvPairs[i + 1]

		//remove ""
		key := removeFirstLast(k)
		value := removeFirstLast(v)

		json[key] = value
	}

	return json
}


func boardToString(board [64]int) string {
	boardState := unflattenBoard(board)

	str := "["
	for i, line := range boardState {
		str += "["

		for i, num := range line {
			str += strconv.Itoa(num)

			if i < len(line) - 1 {
				str += ", "
			}
		}

		str += "]"

		if i < len(boardState) - 1 {
			str += ", "
		}
	}

	//add final ]" for 2d array
	str += "]"

	return str
}


func coordsToString(moveCoords [][2]int) string {
	str := "\"["

	for i, coord := range moveCoords {
		x := strconv.Itoa(coord[0])
		y := strconv.Itoa(coord[1])

		str += "[" + x + ", " + y + "]"

		if i < len(moveCoords) - 1 {
			str += ", "
		}
	}

	str += "]\""

	return str
}


func strToInt(str string) int {
	i, err := strconv.Atoi(str)
	panicErr(err)

	return i
}


func formatAttr(name string, value string) string {
	qName := "\"" + name + "\""
	qVal := "\"" + value + "\""

	str := qName + ": " + qVal

	return str
}


func jsonToState(json map[string]string) engine.GameState {
	board := boardStrToList(json["board"])
	whiteMove := json["white_to_move"] == "true"
	whiteKingCastle := json["white_king_castle"] == "true"
	whiteQueenCastle := json["white_queen_castle"] == "true"
	blackKingCastle := json["black_king_castle"] == "true"
	blackQueenCastle := json["black_queen_castle"] == "true"
	pawnDouble := coordStrToList(json["prev_pawn_double"])

	stateObj := engine.GameState{Board: board, WhiteToMove: whiteMove, WhiteKingCastle: whiteKingCastle, WhiteQueenCastle: whiteQueenCastle, BlackKingCastle: blackKingCastle, BlackQueenCastle: blackQueenCastle, PrevPawnDouble: pawnDouble}
	
	return stateObj
}


func stateToJson(stateObj engine.GameState) string {
	board := boardToString(stateObj.Board)
	whiteMove := strconv.FormatBool(stateObj.WhiteToMove)
	whiteKingCastle := strconv.FormatBool(stateObj.WhiteKingCastle)
	whiteQueenCastle := strconv.FormatBool(stateObj.WhiteQueenCastle)
	blackKingCastle := strconv.FormatBool(stateObj.BlackKingCastle)
	blackQueenCastle := strconv.FormatBool(stateObj.BlackQueenCastle)

	pDoubleSlice := coordsToString([][2]int{stateObj.PrevPawnDouble})
	pDouble := pDoubleSlice[2 : len(pDoubleSlice) - 2] //remove the "[ and ]" at each end

	bAttr := formatAttr("board", board)
	wmAttr := formatAttr("white_to_move", whiteMove)
	wkAttr := formatAttr("white_king_castle", whiteKingCastle)
	wqAttr := formatAttr("white_queen_castle", whiteQueenCastle)
	bkAttr := formatAttr("black_king_castle", blackKingCastle)
	bqAttr := formatAttr("black_queen_castle", blackQueenCastle)
	pdAttr := formatAttr("prev_pawn_double", pDouble)

	attrs := []string{bAttr, wmAttr, wkAttr, wqAttr, bkAttr, bqAttr, pdAttr}
	str := strings.Join(attrs, ", ")

	return str
}


func readFile() string {
	file, err := os.Open(FILE_PATH)

	panicErr(err)

	defer file.Close()

	buffer := make([]byte, 1024)

	//keep reading bytes until there are none left to read
	for {
		readBytes, err := file.Read(buffer)

		if readBytes == 0 {
			break
		} else {
			panicErr(err)
		}
	}

	str := string(buffer)

	return str
}


func LoadGameState() (map[string]string, engine.GameState) {
	str := readFile()
	json := jsonLoad(str)
	
	state := jsonToState(json)

	return json, state
}


func writeToJson(writeStr string) {
	writeData := []byte(writeStr)

	//open file in read/write mode and overwrite existing contents
	file, err := os.Create(FILE_PATH)
	panicErr(err)

	defer file.Close()

	_, err = file.Write(writeData)
	panicErr(err)
}


func WriteState(stateObj engine.GameState) {
	str := stateToJson(stateObj)
	writeStr := "{" + str + "}"

	writeToJson(writeStr)
}


func WriteLegalMoves(moveCoords [][2]int) {
	str := coordsToString(moveCoords)
	writeStr := "{\"moves\": " + str + "}"

	writeToJson(writeStr)
}