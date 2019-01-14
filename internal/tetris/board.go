package tetris

import (
	"fmt"
)

// The default tetris board size.
const (
	defaultBoardWidth  = 10
	defaultBoardHeight = 20
)

// tetrominoMatrices is the slice of matrices for each rotation of each tetromino.
// tetrominoMatrices is a global variable because it is shared between all boards.
// TODO: Is there a more elegant way to do this?
var tetrominoMatrices [][]TetrominoMatrix

func init() {
	tetrominoMatrices = TetrominoMatrices()
}

// Board is a tetris Board.
// Board allows dropping tetrominoes and keeps statistics about the game so far and the current state.
// The zero value of Board is not usable, a new board must be created with the NewBoard method.
type Board struct {
	width  int
	height int

	// cells contains the cells of the board.
	// Cells are indexed from 0 from left to right and top to bottom.
	cells [][]Tetromino

	gameOver bool

	clearedLines       int
	droppedTetrominoes int
	columnHeights      []int
	holes              []int
}

// NewBoard returns a pointer to a new Board.
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

// NewBoardFromBoard returns a pointer to an independent copy of the given board.
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

// Width returns the width of the board.
func (b *Board) Width() int {
	return b.width
}

// Height returns the height of the board.
func (b *Board) Height() int {
	return b.height
}

// GameOver returns true if the game of the current board is over.
func (b *Board) GameOver() bool {
	return b.gameOver
}

// ClearedLines returns the number of lines cleared on the given board.
func (b *Board) ClearedLines() int {
	return b.clearedLines
}

// DroppedTetrominoes returns the number of tetrominoes that have been dropped on the given board.
func (b *Board) DroppedTetrominoes() int {
	return b.droppedTetrominoes
}

func (b *Board) ColumnHeights() []int {
	return b.columnHeights
}

func (b *Board) ColumnHoles() []int {
	return b.holes
}

// At returns the tetromino at the given position of the board.
// If the given position is not occupied, TetrominoEmpty is returned.
// The cells of the board are indexed from 0, left to right and top to bottom.
// At panics if invalid coordinates are provided.
func (b *Board) At(row, col int) Tetromino {
	if !b.isValidCell(row, col) {
		panic(fmt.Errorf("Board.At: invalid coordinates (%d, %d) provided", row, col))
	}

	return b.cells[row][col]
}

// Drop drops a specified rotation of a tetromino such that the leftmost cell is in the given column.
// Drop returns error if tetromino is dropped, but that leads to game over.
// Drop panics if the given tetromino, rotation or column are invalid or if the board's game is already over.
func (b *Board) Drop(tetromino Tetromino, rotation int, column int) error {
	if !tetromino.Valid() {
		panic(fmt.Errorf("Board.Drop: invalid tetromino %d provided", tetromino))
	}

	if rotation < 0 || rotation >= tetromino.RotationsCount() {
		panic(fmt.Errorf(
			"Board.Drop: invalid rotation %d provided for tetromino %s",
			rotation, tetromino,
		))
	}

	if b.gameOver {
		panic(fmt.Errorf("Board.Drop: can not drop: game is over"))
	}

	tetrominoMatrix := tetrominoMatrices[tetromino][rotation]

	if column+len(tetrominoMatrix[0]) > b.width {
		panic(fmt.Errorf(
			"Board.Drop: can not drop: invalid column %d provided for tetromino %s, rotation %d",
			column, tetromino, rotation,
		))
	}

	if !b.canBePut(tetrominoMatrix, 0, column) {
		b.gameOver = true
		for i := range tetrominoMatrix {
			for j := range tetrominoMatrix[i] {
				if tetrominoMatrix[i][j] {
					b.cells[i][column+j] = tetromino
				}
			}
		}
		return fmt.Errorf("Board.Drop: the game just ended.")
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

	b.droppedTetrominoes++
	someRowsCleared := b.clearFullRows()

	fromCol := 0
	toCol := b.width
	if !someRowsCleared {
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

// canBePut returns true if the given tetromino matrix's can be put on the board
// with its top left cell at coordinates row, col.
// If the matrix sticks out of the board or overlaps a non-empty cell on the board, false is returned.
func (b *Board) canBePut(tetrominoMatrix TetrominoMatrix, row, col int) bool {
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

// isValidCell returns true if the given coordinates are valid for the board.
func (b *Board) isValidCell(row, col int) bool {
	return 0 <= row && row < b.height &&
		0 <= col && col < b.width
}

// isValidCell returns true if the given coordinates are valid for the board.
// isFullRow returns true if all of the cells in the given row of the board are non-empty.
func (b *Board) isFullRow(row int) bool {
	for _, cell := range b.cells[row] {
		if cell == TetrominoEmpty {
			return false
		}
	}
	return true
}

// clearFullRows traverses the board and clears any full rows.
// The rows above are then shifted down and the board is filled with emoty rows at the top.
// clearFullRows returns true if at least one row was cleared.
func (b *Board) clearFullRows() bool {
	var (
		someRowsCleared = false
		idxFrom         = b.height - 1
		idxTo           = b.height - 1
	)
	for idxTo >= 0 {
		if idxFrom < 0 {
			// Insert a new empty rows.
			b.cells[idxTo] = make([]Tetromino, b.width)
			idxTo--
			continue
		}

		if b.isFullRow(idxFrom) {
			// Full rows are cleared.
			b.clearedLines++
			someRowsCleared = true
			idxFrom--
			continue
		}

		// Non-full rows are just moved down.
		b.cells[idxTo] = b.cells[idxFrom]
		idxFrom--
		idxTo--
	}

	return someRowsCleared
}
