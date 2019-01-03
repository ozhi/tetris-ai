package tetris

import "fmt"

const TetrominoesCount = 7

type Tetromino int

const (
	TetrominoEmpty Tetromino = iota
	TetrominoI
	TetrominoJ
	TetrominoL
	TetrominoO
	TetrominoS
	TetrominoT
	TetrominoZ
)

func (t Tetromino) IsValid() bool {
	return 0 <= t && t <= TetrominoesCount
}

func (t Tetromino) RotationsCount() int {
	switch t {
	case TetrominoJ, TetrominoL, TetrominoT:
		return 4
	case TetrominoI, TetrominoS, TetrominoZ:
		return 2
	case TetrominoO:
		return 1
	default:
		panic(fmt.Errorf("RotationsCount: invalid tetromino %d", t))
	}
}

type tetrominoMatrix [][]bool

// todo make this a map[TetrominoInstance]tetrominoMatrix
func loadTetrominoMatrices() [][]tetrominoMatrix {
	return [][]tetrominoMatrix{
		// Empty
		[]tetrominoMatrix{},

		// I
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, true, true, true},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// J
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{false, true, false, false},
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, false, false, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, true, true, false},
				[]bool{false, false, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// L
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, true, true, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{false, false, true, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// O
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{true, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// S
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{false, true, true, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
		},

		// T
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{true, true, true, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{false, true, false, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// Z
		[]tetrominoMatrix{
			tetrominoMatrix{
				[]bool{true, true, false, false},
				[]bool{false, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			tetrominoMatrix{
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
		},
	}
}
