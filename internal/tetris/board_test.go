package tetris_test

import (
	"testing"

	"github.com/ozhi/tetris-ai/internal/tetris"
	"github.com/stretchr/testify/assert"
)

type Move struct {
	tetromino tetris.Tetromino
	rotation  int
	column    int
}

func TestBoardSize(t *testing.T) {
	board := tetris.NewBoard()
	assert.Equal(t, 10, board.Width())
	assert.Equal(t, 20, board.Height())
}

func TestNewBoardFromBoardDoesNotShareCellsWithOriginal(t *testing.T) {
	original := tetris.NewBoard()
	original.Drop(tetris.TetrominoO, 0, 0)

	copy := tetris.NewBoardFromBoard(original)

	assert.Equal(t, original.At(18, 1), tetris.TetrominoO)
	assert.Equal(t, copy.At(18, 1), tetris.TetrominoO)

	original.Drop(tetris.TetrominoO, 0, 0)

	assert.Equal(t, original.At(16, 1), tetris.TetrominoO)
	assert.Equal(t, copy.At(16, 1), tetris.TetrominoEmpty)
}

func TestBoardDropPanicsOnInvalidTetromino(t *testing.T) {
	board := tetris.NewBoard()

	tests := []struct {
		tetromino tetris.Tetromino
		panics    bool
	}{
		{tetromino: tetris.Tetromino(-1), panics: true},
		{tetromino: tetris.Tetromino(8), panics: true},
		{tetromino: tetris.TetrominoEmpty, panics: true},
		{tetromino: tetris.TetrominoS, panics: false},
	}

	for _, test := range tests {
		t.Run("Board.Drop panics on invalid tetromino provided", func(t *testing.T) {
			defer func() {
				r := recover()

				if test.panics {
					assert.NotNil(t, r)
					return
				}

				assert.Nil(t, r)
			}()

			board.Drop(test.tetromino, 0, 0)
		})
	}
}

func TestBoardDropPanicsOnInvalidRotation(t *testing.T) {
	board := tetris.NewBoard()

	tests := []struct {
		tetromino tetris.Tetromino
		rotation  int
		panics    bool
	}{
		{tetromino: tetris.TetrominoI, rotation: 1, panics: false},
		{tetromino: tetris.TetrominoI, rotation: 2, panics: true},
		{tetromino: tetris.TetrominoJ, rotation: 3, panics: false},
		{tetromino: tetris.TetrominoJ, rotation: 4, panics: true},
		{tetromino: tetris.TetrominoO, rotation: 0, panics: false},
		{tetromino: tetris.TetrominoO, rotation: 1, panics: true},
	}

	for _, test := range tests {
		t.Run("Board.Drop panics on invalid rotation provided", func(t *testing.T) {
			defer func() {
				r := recover()

				if test.panics {
					assert.NotNil(t, r)
					return
				}

				assert.Nil(t, r)
			}()

			board.Drop(test.tetromino, test.rotation, 4)
		})
	}
}

func TestBoardDrop(t *testing.T) {
	board := tetris.NewBoard()
	moves := []Move{
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 1},
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 2},
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 3},
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 4},
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 5},
		{tetromino: tetris.TetrominoJ, rotation: 2, column: 6},
	}
	for _, move := range moves {
		board.Drop(move.tetromino, move.rotation, move.column)

	}
}

func TestBoardAtPanicsOnInvalidCoordinates(t *testing.T) {
	board := tetris.NewBoard()

	tests := []struct {
		row int
		col int
	}{
		{-5, 7},
		{-1, 7},
		{20, 7},
		{25, 7},
		{7, -5},
		{7, -1},
		{7, 10},
		{7, 15},
	}

	for _, test := range tests {
		t.Run("Board.At panics on invalid coordinates", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r)
			}()

			board.At(test.row, test.col)
		})
	}
}

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
