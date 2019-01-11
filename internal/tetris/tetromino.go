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
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
			},
			{
				[]bool{true, true, true, true},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// J
		{
			{
				[]bool{false, true, false, false},
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, false, false, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, true, true, false},
				[]bool{false, false, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// L
		{
			{
				[]bool{true, false, false, false},
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, true, true, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{false, false, true, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// O
		{
			{
				[]bool{true, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// S
		{
			{
				[]bool{false, true, true, false},
				[]bool{true, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
		},

		// T
		{
			{
				[]bool{true, true, true, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{false, true, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{false, true, false, false},
				[]bool{true, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{true, false, false, false},
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
		},

		// Z
		{
			{
				[]bool{true, true, false, false},
				[]bool{false, true, true, false},
				[]bool{false, false, false, false},
				[]bool{false, false, false, false},
			},
			{
				[]bool{false, true, false, false},
				[]bool{true, true, false, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
			},
		},
	}
}
