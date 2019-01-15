package cli

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ozhi/tetris-ai/internal/ai"
	"github.com/ozhi/tetris-ai/internal/tetris"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CLI is the command-line interface of Tetris-ai.
// CLI encapsulates an AI that plays tetris and visualization logic.
// The zero value of CLI is not usable, method New should be used.
type CLI struct {
	ai *ai.AI
}

// New creates and initializes a new CLI.
func New() *CLI {
	return &CLI{
		ai: ai.New(),
	}
}

// Start starts the AI's game and the visualization loop.
func (cli *CLI) Start() {
	cli.ai.SetNext(tetris.RandomTetromino())
	for {
		board := cli.ai.Board()
		printBoard(board)
		time.Sleep(500 * time.Millisecond)

		err := cli.ai.DropSetNext(tetris.RandomTetromino())
		if err != nil {
			break
		}
	}

	printBoard(cli.ai.Board())
	fmt.Println("Game over")
}

func printBoard(board *tetris.Board) {
	const clearScreen = "\033[2J"

	fmt.Println(clearScreen)
	for row := 0; row < board.Height(); row++ {
		for col := 0; col < board.Width(); col++ {
			if board.At(row, col) == tetris.TetrominoEmpty {
				fmt.Print(".")
			} else {
				fmt.Printf("%s", board.At(row, col))
			}
		}
		fmt.Println()
	}

	for col := 0; col < board.Width(); col++ {
		fmt.Print(col)
	}
	fmt.Printf("   %d\n", board.ClearedLines())
}
