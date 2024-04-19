package board


const WkCastle uint8 = 8
const WqCastle uint8 = 4
const BkCastle uint8 = 2
const BqCastle uint8 = 1

const InvWkCastle uint8 = ^WkCastle
const InvWqCastle uint8 = ^WqCastle
const InvBkCastle uint8 = ^BkCastle
const InvBqCastle uint8 = ^BqCastle


type GameState struct {
	Board [64]int

	WhiteToMove bool

	CastleRights uint8  //Least significant 4 bits act as flags: [X][X][X][X][Wk][Wq][Bk][Bq]

	PrevPawnDouble [2]int

	Bitboards *Bitboard

	DoubleChecked bool
	EnPassantPin bool

	ZobristHash uint64

	//for unmaking moves (end index is the most recent)
	prvBoard stack[[64]int]
	prvCastleRights stack[uint8]
	prvPrevPawnDouble stack[[2]int]
	prvBitboard stack[*Bitboard]
	prvDoubleCheck stack[bool]
	prvEnPassantPin stack[bool]
	PrvZobHash stack[uint64]
}


type Bitboard struct {
	AttackedSquares uint64  //NOTE: this will not include all attacked squares for sliding pieces (only enough to filter out illegal moves)
	PawnAttacks uint64
	AttacksOnKing uint64
	
	PinArray [64]uint64
}


type stack[T any] struct {
	data []T
	top int  //NOTE: the actual index is 1 less than this
	totalSize int
}


func (state *GameState) SetPrevVals() {
	//set prev values
	state.prvBoard.push(state.Board)
	state.prvCastleRights.push(state.CastleRights)
	state.prvPrevPawnDouble.push(state.PrevPawnDouble)
	state.prvBitboard.push(state.Bitboards)
	state.prvDoubleCheck.push(state.DoubleChecked)
	state.prvEnPassantPin.push(state.EnPassantPin)
	state.PrvZobHash.push(state.ZobristHash)
}


func (state *GameState) RestorePrev() {
	//restore the previous values
	state.Board = state.prvBoard.pop()
	state.CastleRights = state.prvCastleRights.pop()
	state.PrevPawnDouble = state.prvPrevPawnDouble.pop()
	state.Bitboards = state.prvBitboard.pop()
	state.DoubleChecked = state.prvDoubleCheck.pop()
	state.EnPassantPin = state.prvEnPassantPin.pop()
	state.ZobristHash = state.PrvZobHash.pop()
}


func (s *stack[T]) push(data T) {
	//for the sake of speed, we assume the stack has not reached its max size
	s.top++

	if s.totalSize < s.top {
		s.totalSize++
		s.data = append(s.data, data)
	} else {
		s.data[s.top - 1] = data
	}
}


func (s *stack[T]) pop() T {
	//for the sake of speed, we assume the stack is not empty without actually checking
	s.top--
	
	return s.data[s.top]
}
