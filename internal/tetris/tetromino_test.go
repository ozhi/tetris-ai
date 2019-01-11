package tetris_test

import (
	"testing"

	"github.com/ozhi/tetris-ai/internal/tetris"
	"github.com/stretchr/testify/assert"
)

func TestTetrominoesCount(t *testing.T) {
	assert.Equal(t, 7, tetris.TetrominoesCount)
}

func TestUninitializedTetromino(t *testing.T) {
	var tetromino tetris.Tetromino
	assert.Equal(t, tetris.TetrominoEmpty, tetromino)
}

func TestTetrominoes(t *testing.T) {
	assert.Equal(t, []tetris.Tetromino{1, 2, 3, 4, 5, 6, 7}, tetris.Tetrominoes())
}

func TestTetrominoIsNonEmpty(t *testing.T) {
	assert.False(t, tetris.Tetromino(-1).IsNonEmpty())
	assert.False(t, tetris.Tetromino(tetris.TetrominoesCount+1).IsNonEmpty())

	assert.False(t, tetris.TetrominoEmpty.IsNonEmpty())
	for _, tetromino := range tetris.Tetrominoes() {
		assert.True(t, tetromino.IsNonEmpty())
	}
}
