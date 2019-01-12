package tetris

import (
	"fmt"
	"math/rand"
)

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

func Tetrominoes() []Tetromino {
	return []Tetromino{
		TetrominoI,
		TetrominoJ,
		TetrominoL,
		TetrominoO,
		TetrominoS,
		TetrominoT,
		TetrominoZ,
	}
}

func RandomTetromino() Tetromino {
	tetrominoes := Tetrominoes()
	return tetrominoes[rand.Intn(TetrominoesCount)]
}

func (t Tetromino) String() string {
	switch t {
	case TetrominoEmpty:
		return "Empty"
	case TetrominoI:
		return "I"
	case TetrominoJ:
		return "J"
	case TetrominoL:
		return "L"
	case TetrominoO:
		return "O"
	case TetrominoS:
		return "S"
	case TetrominoT:
		return "T"
	case TetrominoZ:
		return "Z"
	default:
		panic(fmt.Errorf("Tetromno.String: invalid tetromino %d provided", t))
	}
}

func (t Tetromino) IsNonEmpty() bool {
	return 1 <= t && t <= TetrominoesCount
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

func loadTetrominoMatrices() [][]tetrominoMatrix {
	return [][]tetrominoMatrix{
		// Empty
		{},

		// I
		{
			{
				[]bool{true},
				[]bool{true},
				[]bool{true},
				[]bool{true},
			},
			{
				[]bool{true, true, true, true},
			},
		},

		// J
		{
			{
				[]bool{false, true},
				[]bool{false, true},
				[]bool{true, true},
			},
			{
				[]bool{true, false, false},
				[]bool{true, true, true},
			},
			{
				[]bool{true, true},
				[]bool{true, false},
				[]bool{true, false},
			},
			{
				[]bool{true, true, true},
				[]bool{false, false, true},
			},
		},

		// L
		{
			{
				[]bool{true, false},
				[]bool{true, false},
				[]bool{true, true},
			},
			{
				[]bool{true, true, true},
				[]bool{true, false, false},
			},
			{
				[]bool{true, true},
				[]bool{false, true},
				[]bool{false, true},
			},
			{
				[]bool{false, false, true},
				[]bool{true, true, true},
			},
		},

		// O
		{
			{
				[]bool{true, true},
				[]bool{true, true},
			},
		},

		// S
		{
			{
				[]bool{false, true, true},
				[]bool{true, true, false},
			},
			{
				[]bool{true, false},
				[]bool{true, true},
				[]bool{false, true},
			},
		},

		// T
		{
			{
				[]bool{true, true, true},
				[]bool{false, true, false},
			},
			{
				[]bool{false, true},
				[]bool{true, true},
				[]bool{false, true},
			},
			{
				[]bool{false, true, false},
				[]bool{true, true, true},
			},
			{
				[]bool{true, false},
				[]bool{true, true},
				[]bool{true, false},
			},
		},

		// Z
		{
			{
				[]bool{true, true, false},
				[]bool{false, true, true},
			},
			{
				[]bool{false, true},
				[]bool{true, true},
				[]bool{true, false},
			},
		},
	}
}
