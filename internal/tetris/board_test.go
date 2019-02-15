package tetris_test

import (
	"fmt"
	"testing"

	"github.com/ozhi/tetris-ai/internal/tetris"
	"github.com/stretchr/testify/assert"
)

const (
	Empty = tetris.TetrominoEmpty
	I     = tetris.TetrominoI
	J     = tetris.TetrominoJ
	L     = tetris.TetrominoL
	O     = tetris.TetrominoO
	S     = tetris.TetrominoS
	T     = tetris.TetrominoT
	Z     = tetris.TetrominoZ
)

type Move struct {
	tetromino tetris.Tetromino
	rotation  int
	column    int
}

func ExampleBoard() {
	board := tetris.NewBoard()

	board.Drop(J, 1, 0)
	board.Drop(J, 1, 3)
	board.Drop(O, 0, 6)
	board.Drop(L, 0, 8)

	fmt.Println(board.ClearedLines(), board.At(19, 0))
	// Output: 1 J
}

func TestBoardSize(t *testing.T) {
	board := tetris.NewBoard()
	assert.Equal(t, 10, board.Width())
	assert.Equal(t, 20, board.Height())
}

func TestNewBoardFromBoardDoesNotShareCellsWithOriginal(t *testing.T) {
	original := tetris.NewBoard()
	original.Drop(O, 0, 0)

	copy := tetris.NewBoardFromBoard(original)

	assert.Equal(t, original.At(18, 1), O)
	assert.Equal(t, copy.At(18, 1), O)

	original.Drop(O, 0, 0)

	assert.Equal(t, original.At(16, 1), O)
	assert.Equal(t, copy.At(16, 1), Empty)
}

func TestBoardDropPanicsOnInvalidTetromino(t *testing.T) {
	board := tetris.NewBoard()

	tests := []struct {
		tetromino tetris.Tetromino
		panics    bool
	}{
		{tetromino: tetris.Tetromino(-1), panics: true},
		{tetromino: tetris.Tetromino(8), panics: true},
		{tetromino: Empty, panics: true},
		{tetromino: S, panics: false},
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
		{tetromino: I, rotation: 1, panics: false},
		{tetromino: I, rotation: 2, panics: true},
		{tetromino: J, rotation: 3, panics: false},
		{tetromino: J, rotation: 4, panics: true},
		{tetromino: O, rotation: 0, panics: false},
		{tetromino: O, rotation: 1, panics: true},
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

func TestBoardDropReturnsErrorOnGameOver(t *testing.T) {
	board := tetris.NewBoard()

	for i := 0; i < 9; i++ {
		err := board.Drop(T, 0, 4)
		assert.Nil(t, err)
	}

	err := board.Drop(L, 0, 3)
	assert.NotNil(t, err)
}

func TestBoardDropKeepsStatistics(t *testing.T) {
	board := tetris.NewBoard()

	moves := []Move{
		{tetromino: O, rotation: 0, column: 0},
		{tetromino: O, rotation: 0, column: 2},
		{tetromino: O, rotation: 0, column: 4},
		{tetromino: O, rotation: 0, column: 6},
		{tetromino: O, rotation: 0, column: 8},
		// 2 cleared lines.
		{tetromino: L, rotation: 1, column: 0},
		{tetromino: L, rotation: 3, column: 3},
	}

	for _, move := range moves {
		err := board.Drop(move.tetromino, move.rotation, move.column)
		fmt.Println(board.ColumnHeights())
		assert.Nil(t, err)
	}

	assert.Equal(t, 2, board.ClearedLines())
	assert.Equal(t, []int{2, 2, 2, 1, 1, 2, 0, 0, 0, 0}, board.ColumnHeights())
	assert.Equal(t, []int{0, 1, 1, 0, 0, 0, 0, 0, 0, 0}, board.ColumnHoles())
}

func TestBoardDropClearsMultipleLines(t *testing.T) {
	board := tetris.NewBoard()
	moves := []Move{
		{tetromino: J, rotation: 1, column: 0},
		{tetromino: L, rotation: 3, column: 3},
		{tetromino: Z, rotation: 0, column: 0},
		{tetromino: T, rotation: 2, column: 2},
		{tetromino: L, rotation: 2, column: 5},
		{tetromino: O, rotation: 0, column: 7},
		{tetromino: Z, rotation: 0, column: 6},
		{tetromino: I, rotation: 0, column: 9},
	}
	for _, move := range moves {
		board.Drop(move.tetromino, move.rotation, move.column)
	}

	expected18 := []tetris.Tetromino{Empty, Empty, Empty, T, Empty, Empty, Z, Z, Empty, I}
	expected19 := []tetris.Tetromino{J, Z, Z, Empty, Empty, L, L, O, O, I}

	for col := 0; col < board.Width(); col++ {
		for row := 0; row <= 17; row++ {
			assert.Equal(t, Empty, board.At(row, col))
		}
		assert.Equal(t, expected18[col], board.At(18, col))
		assert.Equal(t, expected19[col], board.At(19, col))
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

func TestBoardAtNumbersCellsFromTopLeft(t *testing.T) {
	board := tetris.NewBoard()
	board.Drop(I, 0, 4)
	board.Drop(L, 3, 3)

	assert.Equal(t, I, board.At(19, 4))
	assert.Equal(t, I, board.At(18, 4))
	assert.Equal(t, I, board.At(17, 4))
	assert.Equal(t, I, board.At(16, 4))

	assert.Equal(t, L, board.At(15, 3))
	assert.Equal(t, L, board.At(15, 4))
	assert.Equal(t, L, board.At(15, 5))
	assert.Equal(t, L, board.At(14, 5))

	assert.Equal(t, Empty, board.At(0, 0))
	assert.Equal(t, Empty, board.At(19, 9))
	assert.Equal(t, Empty, board.At(18, 3))
	assert.Equal(t, Empty, board.At(14, 4))
	assert.Equal(t, Empty, board.At(15, 6))
}
