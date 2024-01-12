package api


import (
	"os"
	"strconv"
	"strings"
)


func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}


func strToList(str string) [8][8]int {
	//convert the string "[[1, 2], [3, 4]]" to array [[1, 2], [3, 4]]

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

	return outList
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


func stateToString(boardState [8][8]int) string {
	str := "\"["
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
	str += "]\""

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


func LoadData() (map[string]string, [8][8]int) {
	file, err := os.Open("src/api/interface.json")

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

	json := jsonLoad(str)
	board := strToList(json["board"])

	return json, board
}


func writeToJson(writeStr string) {
	writeData := []byte(writeStr)

	//open file in read/write mode and overwrite existing contents
	file, err := os.Create("src/api/interface.json")
	panicErr(err)

	defer file.Close()

	_, err = file.Write(writeData)
	panicErr(err)
}


func WriteBoardState(boardState [8][8]int) {
	str := stateToString(boardState)
	writeStr := "{\"board\": " + str + "}"

	writeToJson(writeStr)
}


func WriteLegalMoves(moveCoords [][2]int) {
	str := coordsToString(moveCoords)
	writeStr := "{\"moves\": " + str + "}"

	writeToJson(writeStr)
}