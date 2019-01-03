package ai

import (
	"fmt"
	"math/rand"

	"github.com/ozhi/tetris-ai/internal/tetris"
)

type AI struct {
	board *tetris.Board

	current tetris.Tetromino
}

func New() *AI {
	return &AI{
		board: tetris.NewBoard(),
	}
}

func (ai *AI) Board() *tetris.Board {
	return ai.board
}

func evaluate(board *tetris.Board, depth int) float64 {
	const inf = 10000.0
	if depth < 0 {
		panic("dai polojitelna dalbochina we")
	}

	if depth == 0 || board.GameOver() {
		return utility(board)
	}

	minEval := inf
	for tetromino := tetris.Tetromino(1); tetromino <= tetris.Tetromino(tetris.TetrominoesCount); tetromino++ {
		maxEval := -inf
		for rotation := 0; rotation < tetromino.RotationsCount(); rotation++ {
			for column := 0; column < board.Width(); column++ {
				newBoard := tetris.NewBoardFromBoard(board)
				if err := newBoard.Drop(tetromino, rotation, column); err != nil {
					// Invalid column.
					continue
				}

				if eval := evaluate(newBoard, depth-1); eval > maxEval {
					maxEval = eval
				}
			}
		}

		if maxEval < minEval {
			minEval = maxEval
		}
	}
	return minEval
}

func (ai *AI) Drop(tetromino tetris.Tetromino) error {
	if ai.board.GameOver() {
		panic(fmt.Errorf("game is over we"))
	}

	type Move struct {
		rotation int
		column   int
	}

	var (
		bestEval  = -10000.0
		bestMoves []Move
	)

	for rotation := 0; rotation < tetromino.RotationsCount(); rotation++ {
		for column := 0; column < ai.board.Width(); column++ {
			newBoard := tetris.NewBoardFromBoard(ai.board)
			if err := newBoard.Drop(tetromino, rotation, column); err != nil {
				// column is invalid.
				continue
			}

			eval := evaluate(newBoard, 1)
			if eval > bestEval {
				bestEval = eval
				bestMoves = []Move{
					Move{
						rotation: rotation,
						column:   column,
					},
				}
			} else if eval == bestEval {
				bestMoves = append(bestMoves, Move{
					rotation: rotation,
					column:   column,
				})
			}
		}
	}

	if len(bestMoves) == 0 {
		bestMoves = []Move{Move{rotation: 0, column: 0}} // will lose
	}

	move := bestMoves[rand.Intn(len(bestMoves))]

	if err := ai.board.Drop(tetromino, move.rotation, move.column); err != nil {
		return fmt.Errorf("AI.Drop: could not drop: %s", err)
	}

	return nil
}

func (ai *AI) SetCurrent(current tetris.Tetromino) {
	ai.current = current
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func utility(board *tetris.Board) float64 {
	if board.GameOver() {
		return -1e5
	}

	var (
		columnHeights = board.ColumnHeights()
		columnHoles   = board.ColumnHoles()

		heightsSum  int
		heightsDiff int
		holes       int
	)

	for col := 0; col < board.Width(); col++ {
		heightsSum += columnHeights[col]
		holes += columnHoles[col]
		if col != 0 {
			heightsDiff += abs(columnHeights[col] - columnHeights[col-1])
		}
	}

	utility := -0.510066*float64(heightsSum) +
		0.760666*float64(board.ClearedLines()) +
		-0.35663*float64(holes) +
		-0.184483*float64(heightsDiff)

	if board.GameOver() {
		utility -= 1e4
	}

	return utility
}
