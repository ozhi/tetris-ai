package ai_test

import (
	"fmt"
	"testing"

	"github.com/ozhi/tetris-ai/internal/ai"
	"github.com/ozhi/tetris-ai/internal/tetris"
)

func ExampleAI() {
	ai := ai.New()

	ai.SetNext(tetris.TetrominoL)
	moves := []tetris.Tetromino{
		tetris.TetrominoZ,
		tetris.TetrominoT,
		tetris.TetrominoJ,
		tetris.TetrominoO,
	}
	for _, move := range moves {
		ai.DropSetNext(move)
	}

	fmt.Println(ai.Board().DroppedTetrominoes())
	// Output: 4
}

func benchmarkDropSetNext(tetrominoesToDrop int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		ai := ai.New() // Start with a fresh board each time.
		ai.SetNext(tetris.RandomTetromino())

		for t := 0; t < tetrominoesToDrop; t++ {
			ai.DropSetNext(tetris.RandomTetromino())
		}
	}
}

func BenchmarkDropSetNext1(b *testing.B) {
	benchmarkDropSetNext(1, b)
	benchmarkDropSetNext(10, b)
	benchmarkDropSetNext(100, b)
	benchmarkDropSetNext(1000, b)
}

func BenchmarkDropSetNext10(b *testing.B) {
	benchmarkDropSetNext(10, b)
}

func BenchmarkDropSetNext100(b *testing.B) {
	benchmarkDropSetNext(100, b)
}

func BenchmarkDropSetNext1000(b *testing.B) {
	benchmarkDropSetNext(1000, b)
}
func BenchmarkDropSetNext10000(b *testing.B) {
	benchmarkDropSetNext(10000, b)
}
