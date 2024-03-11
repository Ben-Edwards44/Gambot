package evaluation


//material evals
const pawnWeight int = 10
const knightWeight int = 30
const bishopWeight int = 30
const rookWeight int = 50
const queenWeight int = 90

//phase weights
const pawnPhase int = 0
const knightPhase int = 1
const bishopPhase int = 1
const rookPhase int = 2
const queenPhase int = 4
const totalPhase int = 16 * pawnPhase + 4 * knightPhase + 4 * bishopPhase + 4 * rookPhase + 2 * queenPhase  //NOTE: this is doubled because white and black are accounted for


//NOTE: the king weight and phase has been left as 0
var pieceWeights [6]int = [6]int{pawnWeight, knightWeight, bishopWeight, rookWeight, 0, queenWeight}
var phaseWeights [6]int = [6]int{pawnPhase, knightPhase, bishopPhase, rookPhase, 0, queenPhase}