package tetris_test

// type Move struct {
// 	tetromino tetris.Tetromino
// 	rotation  int
// 	column    int
// }

// func clearLineMoves() []Move {
// 	return []Move{
// 		Move{tetris.TetrominoI, 1, 0},
// 		Move{tetris.TetrominoI, 1, 5},
// 		Move{tetris.TetrominoT, 1, 8},
// 		Move{tetris.TetrominoJ, 3, 2},
// 	}
// }

// func dropFromTopToBottomMoves() []Move {
// 	return []Move{
// 		Move{tetris.TetrominoI, 0, 4},
// 		Move{tetris.TetrominoI, 0, 8},
// 		Move{tetris.TetrominoI, 1, 4},
// 		Move{tetris.TetrominoO, 0, 5},
// 	}
// }

// func randomMoves() []Move {
// 	var moves []Move
// 	for i := 0; i < 20; i++ {
// 		move := Move{
// 			tetromino: tetris.Tetromino(1 + rand.Intn(7)),
// 		}

// 		var possibleRotations int
// 		switch move.tetromino {
// 		case tetris.TetrominoI:
// 			possibleRotations = 2
// 		case tetris.TetrominoJ:
// 			possibleRotations = 4
// 		case tetris.TetrominoL:
// 			possibleRotations = 4
// 		case tetris.TetrominoO:
// 			possibleRotations = 1
// 		case tetris.TetrominoS:
// 			possibleRotations = 2
// 		case tetris.TetrominoT:
// 			possibleRotations = 4
// 		case tetris.TetrominoZ:
// 			possibleRotations = 2
// 		}

// 		move.rotation = rand.Intn(possibleRotations)

// 		move.column = rand.Intn(8)

// 		moves = append(moves, move)
// 	}
// 	return moves
// }
