package api


import (
	"fmt"
	"os"
    "strings"
    "strconv"
)


func strToList(str string) [8][8]int {
    //convert the string "[[1, 2], [3, 4]]" to array [[1, 2], [3, 4]]

    //remove the [[ and ]] at end (2 because there are 2 chars that must be removed)
    str = removeFirstLast(str)
    str = removeFirstLast(str)

    strLs := strings.Split(str, "], [")

    var outList [8][8]int
    for i := 0; i < 8; i++ {
        nums := strings.Split(strLs[i], ", ")

        for x := 0; x < 8; x++ {
            num, err := strconv.Atoi(nums[x])

            if err != nil {
                panic(err)
            }

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
        backInx := inx - 2
        numQuote := 0

        for backInx >= 0 && numQuote < 2 {
            key += string(str[backInx])

            if str[backInx] == '"' {
                numQuote++
            }

            backInx--
        }

        value := ""
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

    fmt.Println(kvPairs)

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


func LoadBoardState() [8][8]int {    
    file, err := os.Open("src/API/interface.json")

    if err != nil {
        panic(err)
    }

    defer file.Close()

    buffer := make([]byte, 1024)

    //keep reading bytes until there are none left to read
    for {
        readBytes, err := file.Read(buffer)

        if readBytes == 0 {
            break
        } else if err != nil {
            panic(err)
        }
    }

    str := string(buffer)
    fmt.Println(str)
    json := jsonLoad(str)

    fmt.Println(json)

    board := strToList(json["board"])

    fmt.Println(board)

    return board
}