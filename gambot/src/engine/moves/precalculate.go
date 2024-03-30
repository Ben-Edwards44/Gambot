package moves


var edgeDists [512]int


func PrecalculateEdgeDists() {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			up := x
			down := 7 - x
			left := y
			right := 7 - y

			inx := x * 64 + y * 8

			edgeDists[inx] = up
			edgeDists[inx + 1] = down
			edgeDists[inx + 2] = left
			edgeDists[inx + 3] = right
			edgeDists[inx + 4] = min(up, left)
			edgeDists[inx + 5] = min(up, right)
			edgeDists[inx + 6] = min(down, left)
			edgeDists[inx + 7] = min(down, right)
		}
	}
}