package tetris

import (
	"fmt"
)

const (
	defaultBoardWidth  = 10
	defaultBoardHeight = 20
)

var tetrominoMatrices [][]tetrominoMatrix

func init() {
	tetrominoMatrices = loadTetrominoMatrices()
}

type Board struct {
	width  int
	height int
	cells  [][]Tetromino

	gameOver bool

	clearedLines       int
	droppedTetrominoes int
	columnHeights      []int
	holes              []int
}

func NewBoard() *Board {
	board := Board{
		width:         defaultBoardWidth,
		height:        defaultBoardHeight,
		columnHeights: make([]int, defaultBoardWidth),
		holes:         make([]int, defaultBoardWidth),
	}

	board.cells = make([][]Tetromino, board.height)
	for row := range board.cells {
		board.cells[row] = make([]Tetromino, board.width)
	}

	return &board
}

func NewBoardFromBoard(other *Board) *Board {
	board := *other

	board.cells = make([][]Tetromino, other.height)
	for row := range board.cells {
		board.cells[row] = make([]Tetromino, other.width)
		copy(board.cells[row], other.cells[row])
	}

	board.columnHeights = make([]int, other.width)
	copy(board.columnHeights, other.columnHeights)

	board.holes = make([]int, other.width)
	copy(board.holes, other.holes)

	return &board
}

func (b *Board) Width() int {
	return b.width
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) GameOver() bool {
	return b.gameOver
}

func (b *Board) ClearedLines() int {
	return b.clearedLines
}

func (b *Board) DroppedTetrominoes() int {
	return b.droppedTetrominoes
}

func (b *Board) ColumnHeights() []int {
	return b.columnHeights
}

func (b *Board) ColumnHoles() []int {
	return b.holes
}

func (b *Board) At(row, col int) Tetromino {
	if !b.isValidCell(row, col) {
		panic(fmt.Errorf("Board.At: invalid coordinates (%d, %d) provided", row, col))
	}

	return b.cells[row][col]
}

func (b *Board) Drop(tetromino Tetromino, rotation int, column int) error {
	if !tetromino.IsNonEmpty() {
		panic(fmt.Errorf("Board.Drop: invalid tetromino %d provided", tetromino))
	}

	if rotation < 0 || rotation >= tetromino.RotationsCount() {
		panic(fmt.Errorf(
			"Board.Drop: invalid rotation %d provided for tetromino %d",
			rotation, tetromino,
		))
	}

	if b.gameOver {
		return fmt.Errorf("Board.Drop: can not drop: game is over")
	}

	tetrominoMatrix := tetrominoMatrices[tetromino][rotation]

	if !b.isValidColumn(column, tetrominoMatrix) {
		return fmt.Errorf(
			"Board.Drop: can not drop: invalid column %d provided for tetromino %d, rotation %d",
			column, tetromino, rotation,
		)
	}

	// TODO: drop tetrominoes from above the board, add tests for this. Or not?

	if !b.canBePut(tetrominoMatrix, 0, column) {
		b.gameOver = true
		for i := range tetrominoMatrix {
			for j := range tetrominoMatrix[i] {
				if tetrominoMatrix[i][j] {
					b.cells[i][column+j] = tetromino
				}
			}
		}
		return nil
		// return fmt.Errorf("Board.Drop: tetromino dropped, game is over")
	}

	var row int
	for row = range b.cells {
		if !b.canBePut(tetrominoMatrix, row, column) {
			row--
			break
		}
	}

	for i := range tetrominoMatrix {
		for j := range tetrominoMatrix[i] {
			if tetrominoMatrix[i][j] {
				b.cells[row+i][column+j] = tetromino
			}
		}
	}

	linesCleared := false
	for i := row; i < row+len(tetrominoMatrix); i++ {
		linesCleared = linesCleared || b.clearLineIfFull(i)
	}

	b.droppedTetrominoes++

	fromCol := 0
	toCol := b.width
	if !linesCleared {
		fromCol = column
		toCol = column + len(tetrominoMatrix[0])
	}

	for col := fromCol; col < toCol; col++ {
		if col >= b.width {
			break
		}

		inHole := false
		for row := 0; row < b.height; row++ {
			if inHole && b.cells[row][col] == TetrominoEmpty {
				b.holes[col]++
			}
			if !inHole && b.cells[row][col] != TetrominoEmpty {
				inHole = true
				b.columnHeights[col] = b.height - row
			}
		}
	}

	return nil
}

func (b *Board) canBePut(tetrominoMatrix tetrominoMatrix, row, col int) bool {
	for i := range tetrominoMatrix {
		for j := range tetrominoMatrix[i] {
			if tetrominoMatrix[i][j] {
				if !b.isValidCell(row+i, col+j) || b.cells[row+i][col+j] != TetrominoEmpty {
					return false
				}
			}
		}
	}
	return true
}

func (b *Board) isValidCell(row, col int) bool {
	return 0 <= row && row < b.height &&
		0 <= col && col < b.width
}

func (b *Board) clearLineIfFull(line int) bool {
	if line < 0 || line >= b.height {
		return false
	}

	lineIsFull := true
	for col := range b.cells[line] {
		if b.cells[line][col] == TetrominoEmpty {
			lineIsFull = false
			break
		}
	}

	if !lineIsFull {
		return false
	}

	for row := line; row >= 1; row-- {
		for col := range b.cells[row] {
			b.cells[row][col] = b.cells[row-1][col]
		}
	}
	for col := range b.cells[0] {
		b.cells[0][col] = TetrominoEmpty
	}

	b.clearedLines++
	return true
}

func (b *Board) isValidColumn(column int, tetrominoMatrix tetrominoMatrix) bool {
	if column < 0 {
		return false
	}

	tetrominoWidth := 0
	for row := range tetrominoMatrix {
		for col := range tetrominoMatrix[row] {
			if tetrominoMatrix[row][col] && col > tetrominoWidth {
				tetrominoWidth = col
			}
		}
	}

	return column+tetrominoWidth < b.width
}
