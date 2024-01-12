package src


import (
	"chess-engine/src/api"
	"chess-engine/src/engine"
	"fmt"
)


func engineMove(currentPosition [8][8]int) {
	newBoard := engine.CalculateMove(currentPosition)
	api.WriteBoardState(newBoard)
}


func Main() {
	json, parsedBoard := api.LoadData()
	action := json["task"]

	if action == "move_gen" {
		engineMove(parsedBoard)
	} else if action == "legal_moves" {
		fmt.Println("not implemented yet")
		fmt.Println(json["piece_x"])
		fmt.Println(json["piece_y"])
	}
}