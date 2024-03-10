package evaluation


const pawnWeight int = 10
const knightWeight int = 30
const bishopWeight int = 30
const rookWeight int = 50
const kingWeight int = 100  //arbitrary value because both sides will have a king
const queenWeight int = 90


var pieceWeight [6]int = [6]int{pawnWeight, knightWeight, bishopWeight, rookWeight, kingWeight, queenWeight}
