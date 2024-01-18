package moves


func horBlock(startX int, startY int, kingX int, kingY int) uint64 {
	start := min(startY, kingY)
	end := max(startY, kingY)

	var bitBoard uint64
	for i := start; i <= end; i++ {
		pos := startX * 8 + i
		setBitBoard(&bitBoard, pos)
	}

	return bitBoard
}


func vertBlock(startX int, startY int, kingX int, kingY int) uint64 {
	start := min(startX, kingX)
	end := max(startX, kingX)

	var bitBoard uint64
	for i := start; i <= end; i++ {
		pos := i * 8 + startY
		setBitBoard(&bitBoard, pos)
	}

	return bitBoard
}


func diagBlock(startX int, startY int, kingX int, kingY int) uint64 {
	xStep := -1
	diff := startX - kingX
	if startX < kingX {
		xStep = 1
		diff = kingX - startX
	}

	yStep := -1
	if startY < kingY {
		yStep = 1
	}

	var bitBoard uint64
	for i := 0; i <= diff; i++ {
		x := startX + i * xStep
		y := startY + i * yStep

		setBitBoard(&bitBoard, x * 8 + y)
	}

	return bitBoard
}


func getBlockBitBoards(move Move) uint64 {
	sX := move.StartX
	sY := move.StartY
	kX := move.EndX
	kY := move.EndY

	var bitBoard uint64
	setBitBoard(&bitBoard, sX * 8 + sY)  //allow capture no of attacking piece matter what

	if sY == kY {
		bitBoard |= vertBlock(sX, sY, kX, kY)
	} else if sX == kX {
		bitBoard |= horBlock(sX, sY, kX, kY)
	} else if move.PieceValue == 3 || move.PieceValue == 9 || move.PieceValue == 6 || move.PieceValue == 12 {
		bitBoard |= diagBlock(sX, sY, kX, kY)
	}

	return bitBoard
}


func getKingAttackBlock(kingX int, kingY int, otherMoves []Move) []uint64 {
	//return a slice of bitboards containing the moves that will block an attack on the king

	var bitBoards []uint64
	for _, i := range otherMoves {
		if i.EndX == kingX && i.EndY == kingY {
			bb := getBlockBitBoards(i)
			bitBoards = append(bitBoards, bb)
		}
	}

	return bitBoards
}