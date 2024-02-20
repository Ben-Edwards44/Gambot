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

	WhitePiecePos [6][10][2]int
	BlackPiecePos [6][10][2]int

	NoKingMoveBitBoard uint64
	KingAttackBlocks []uint64

	PinArray [64]uint64
	EnPassantPin bool

	ZobristHash uint64

	//for unmaking moves (end index is the most recent)
	prvBoard [][64]int

	prvWhiteToMove []bool

	prvCastleRights []uint8

	prvPrevPawnDouble [][2]int

	PrvWhitePiecePos [][6][10][2]int
	PrvBlackPiecePos [][6][10][2]int

	prvNoKingMoveBitBoard []uint64
	prvKingAttackBlocks [][]uint64

	prvPinArray [][64]uint64
	prvEnPassantPin []bool
}


func (state *GameState) SetPrevVals() {
	//copy the slice (because slices are passed by reference)
	cpyKingAttackBlocks := make([]uint64, len(state.KingAttackBlocks))
	copy(cpyKingAttackBlocks, state.KingAttackBlocks)

	//set prev values
	state.prvBoard = append(state.prvBoard, state.Board)
	state.prvWhiteToMove = append(state.prvWhiteToMove, state.WhiteToMove)
	state.prvCastleRights = append(state.prvCastleRights, state.CastleRights)
	state.prvPrevPawnDouble = append(state.prvPrevPawnDouble, state.PrevPawnDouble)
	state.PrvWhitePiecePos = append(state.PrvWhitePiecePos, state.WhitePiecePos)
	state.PrvBlackPiecePos = append(state.PrvBlackPiecePos, state.BlackPiecePos)
	state.prvNoKingMoveBitBoard = append(state.prvNoKingMoveBitBoard, state.NoKingMoveBitBoard)
	state.prvKingAttackBlocks = append(state.prvKingAttackBlocks, cpyKingAttackBlocks)
	state.prvPinArray = append(state.prvPinArray, state.PinArray)
	state.prvEnPassantPin = append(state.prvEnPassantPin, state.EnPassantPin)
}


func (state *GameState) RestorePrev() {
	//restore the previous values. Note that the slices are only shallow copies
	
	//restore value
	state.Board = state.prvBoard[len(state.prvBoard) - 1]
	state.WhiteToMove = state.prvWhiteToMove[len(state.prvWhiteToMove) - 1]
	state.CastleRights = state.prvCastleRights[len(state.prvCastleRights) - 1]
	state.PrevPawnDouble = state.prvPrevPawnDouble[len(state.prvPrevPawnDouble) - 1]
	state.WhitePiecePos = state.PrvWhitePiecePos[len(state.PrvWhitePiecePos) - 1]
	state.BlackPiecePos = state.PrvBlackPiecePos[len(state.PrvBlackPiecePos) - 1]
	state.NoKingMoveBitBoard = state.prvNoKingMoveBitBoard[len(state.prvNoKingMoveBitBoard) - 1]
	state.KingAttackBlocks = state.prvKingAttackBlocks[len(state.prvKingAttackBlocks) - 1]
	state.PinArray = state.prvPinArray[len(state.prvPinArray) - 1]
	state.EnPassantPin = state.prvEnPassantPin[len(state.prvEnPassantPin) - 1]

	//pop end of slice
	state.prvBoard = state.prvBoard[:len(state.prvBoard) - 1]
	state.prvWhiteToMove = state.prvWhiteToMove[:len(state.prvWhiteToMove) - 1]
	state.prvCastleRights = state.prvCastleRights[:len(state.prvCastleRights) - 1]
	state.prvPrevPawnDouble = state.prvPrevPawnDouble[:len(state.prvPrevPawnDouble) - 1]
	state.PrvWhitePiecePos = state.PrvWhitePiecePos[:len(state.PrvWhitePiecePos) - 1]
	state.PrvBlackPiecePos = state.PrvBlackPiecePos[:len(state.PrvBlackPiecePos) - 1]
	state.prvNoKingMoveBitBoard = state.prvNoKingMoveBitBoard[:len(state.prvNoKingMoveBitBoard) - 1]
	state.prvKingAttackBlocks = state.prvKingAttackBlocks[:len(state.prvKingAttackBlocks) - 1]
	state.prvPinArray = state.prvPinArray[:len(state.prvPinArray) - 1]
	state.prvEnPassantPin = state.prvEnPassantPin[:len(state.prvEnPassantPin) - 1]
}