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
	prvBoard [][64]int
	prvWhiteToMove []bool
	prvCastleRights []uint8
	prvPrevPawnDouble [][2]int
	prvBitboard []*Bitboard
	prvDoubleCheck []bool
	prvEnPassantPin []bool
	prvZobHash []uint64
}


type Bitboard struct {
	AttackedSquares uint64  //NOTE: this will not include all attacked squares for sliding pieces (only enough to filter out illegal moves)
	PawnAttacks uint64
	AttacksOnKing uint64
	
	PinArray [64]uint64
}


func (state *GameState) SetPrevVals() {
	//set prev values
	//NOTE: any slices will be passed by reference, so must be manually copied
	state.prvBoard = append(state.prvBoard, state.Board)
	state.prvWhiteToMove = append(state.prvWhiteToMove, state.WhiteToMove)
	state.prvCastleRights = append(state.prvCastleRights, state.CastleRights)
	state.prvPrevPawnDouble = append(state.prvPrevPawnDouble, state.PrevPawnDouble)
	state.prvBitboard = append(state.prvBitboard, state.Bitboards)
	state.prvDoubleCheck = append(state.prvDoubleCheck, state.DoubleChecked)
	state.prvEnPassantPin = append(state.prvEnPassantPin, state.EnPassantPin)
	state.prvZobHash = append(state.prvZobHash, state.ZobristHash)
}


func (state *GameState) RestorePrev() {
	//restore the previous values. Note that the slices are only shallow copies
	
	//restore value
	state.Board = state.prvBoard[len(state.prvBoard) - 1]
	state.WhiteToMove = state.prvWhiteToMove[len(state.prvWhiteToMove) - 1]
	state.CastleRights = state.prvCastleRights[len(state.prvCastleRights) - 1]
	state.PrevPawnDouble = state.prvPrevPawnDouble[len(state.prvPrevPawnDouble) - 1]
	state.Bitboards = state.prvBitboard[len(state.prvBitboard) - 1]
	state.DoubleChecked = state.prvDoubleCheck[len(state.prvDoubleCheck) - 1]
	state.EnPassantPin = state.prvEnPassantPin[len(state.prvEnPassantPin) - 1]
	state.ZobristHash = state.prvZobHash[len(state.prvZobHash) - 1]

	//pop end of slice
	state.prvBoard = state.prvBoard[:len(state.prvBoard) - 1]
	state.prvWhiteToMove = state.prvWhiteToMove[:len(state.prvWhiteToMove) - 1]
	state.prvCastleRights = state.prvCastleRights[:len(state.prvCastleRights) - 1]
	state.prvPrevPawnDouble = state.prvPrevPawnDouble[:len(state.prvPrevPawnDouble) - 1]
	state.prvBitboard = state.prvBitboard[:len(state.prvBitboard) - 1]
	state.prvDoubleCheck = state.prvDoubleCheck[:len(state.prvDoubleCheck) - 1]
	state.prvEnPassantPin = state.prvEnPassantPin[:len(state.prvEnPassantPin) - 1]
	state.prvZobHash = state.prvZobHash[:len(state.prvZobHash) - 1]
}