package cli

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ozhi/tetris-ai/internal/tetris"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	board *tetris.Board
}

func NewGame() *Game {
	return &Game{
		board: tetris.NewBoard(),
	}
}

func (g *Game) Start() {
	// moves := []Move{
	// 	Move{tetris.TetrominoI, 1, 0},
	// 	Move{tetris.TetrominoI, 1, 5},
	// 	Move{tetris.TetrominoT, 1, 8},
	// 	Move{tetris.TetrominoJ, 3, 2},
	// }

	for {
		printBoard(g.board)
		time.Sleep(500 * time.Millisecond)

		move := randMove()
		if err := g.board.Drop(move.tetromino, move.rotation, move.column); err != nil {
			fmt.Println(err)
		}
	}

	printBoard(g.board)
	fmt.Println("Game has ended.")
}

type Move struct {
	tetromino tetris.Tetromino
	rotation  int
	column    int
}

func randMove() *Move {
	move := Move{
		tetromino: tetris.RandomTetromino(),
	}

	var possibleRotations int
	switch move.tetromino {
	case tetris.TetrominoI:
		possibleRotations = 2
	case tetris.TetrominoJ:
		possibleRotations = 4
	case tetris.TetrominoL:
		possibleRotations = 4
	case tetris.TetrominoO:
		possibleRotations = 1
	case tetris.TetrominoS:
		possibleRotations = 2
	case tetris.TetrominoT:
		possibleRotations = 4
	case tetris.TetrominoZ:
		possibleRotations = 2
	}

	move.rotation = rand.Intn(possibleRotations)

	move.column = rand.Intn(8)

	return &move
}

func printBoard(board *tetris.Board) {
	const clearScreen = "\033[2J"

	fmt.Println(clearScreen)
	for row := 0; row < board.Height(); row++ {
		for col := 0; col < board.Width(); col++ {
			fmt.Print(CellToString(board.At(row, col)))
		}
		fmt.Println()
	}

	for col := 0; col < board.Width(); col++ {
		fmt.Print(col)
	}
	fmt.Printf("   %d\n", board.ClearedLines())
}

func CellToString(tetromino tetris.Tetromino) string {
	switch tetromino {
	case tetris.TetrominoEmpty:
		return "."
	case tetris.TetrominoI:
		return "I"
	case tetris.TetrominoJ:
		return "J"
	case tetris.TetrominoL:
		return "L"
	case tetris.TetrominoO:
		return "O"
	case tetris.TetrominoS:
		return "S"
	case tetris.TetrominoT:
		return "T"
	case tetris.TetrominoZ:
		return "Z"
	default:
		return "?"
	}
}
