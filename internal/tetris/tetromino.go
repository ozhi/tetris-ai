package tetris

import (
	"fmt"
	"math/rand"
)

// TetrominoesCount is the number of (valid, non-empty) tetrominoes.
const TetrominoesCount = 7

// Tetromino is a type of tetromino - I, J, L, O, S, T or Z without any rotation specified.
// The zero value of the the is the empty tetromino, which is not valid for use in most cases.
type Tetromino int

// The only valid values for a tetromino I, J, L, O, S, T or Z and the zero value of the type - empty.
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

// Tetrominoes returns a slice of all the valid non-empty tetrominoes.
// Convenient for ranging.
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

// RandomTetromino returns a pseudo-random valid non-empty tetromino with uniform distribution.
func RandomTetromino() Tetromino {
	tetrominoes := Tetrominoes()
	return tetrominoes[rand.Intn(TetrominoesCount)]
}

// Tetromino implements Stringer.
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

// Valid returns true if the given tetromino is valid and not empty.
// Returns false otherwise.
func (t Tetromino) Valid() bool {
	return 1 <= t && t <= TetrominoesCount
}

// RotationsCount returns the number of different rotations for the given tetromino.
// J, L and T have 4 rotations, I, S and Z have 2 and O has 1.RotationsCount
// The method panics if called on an empty or invalid tetromino.
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

// TetrominoMatrix is a description of some rotation of some tetromino in 2D space.
// A cell in the matrix is true only if the corresponding cell is occupied.
// A TetrominoMatrix must be recangular - all rows must be of the same length.
// A TetrominoMatrix must be as small as possible - no empty columns on the left or right and
// no empty rows on the top ot bottom.
type TetrominoMatrix [][]bool

// TetrominoMatrices returns a two-dimensional slice of tetromino matrices.
// The first index corresponds to the type of tetromino.
// The second index corresponds to its rotation.
// The zeroth rotation of each tetromino is the one that looks like the letter which describes it.
// The k+1-th rotation of each tetromino is like its k-th, but rotated 90 degrees clockwise.
// For example tetrominoMatrices[TetrominoL][2] is the matrix which describes the second rotation of the L tetromino.
func TetrominoMatrices() [][]TetrominoMatrix {
	return [][]TetrominoMatrix{
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
