package ai

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ozhi/tetris-ai/internal/tetris"
)

type AI struct {
	board *tetris.Board

	nextTetromino tetris.Tetromino
}

func New() *AI {
	return &AI{
		board: tetris.NewBoard(),
	}
}

func (ai *AI) Board() *tetris.Board {
	return ai.board
}
func (ai *AI) DropSetNext(next tetris.Tetromino) error {
	if ai.nextTetromino == tetris.TetrominoEmpty {
		ai.nextTetromino = next
		return nil
	}

	if ai.board.GameOver() {
		panic(fmt.Errorf("AI.DropSetNext: can not drop, game is over"))
	}

	type Move struct {
		rotation int
		column   int
	}

	var (
		bestEval  = -10000.0
		bestMoves []Move
	)

	for curRot := 0; curRot < ai.nextTetromino.RotationsCount(); curRot++ {
		for curCol := 0; curCol < ai.board.Width(); curCol++ {
			curMove := Move{
				rotation: curRot,
				column:   curCol,
			}

			curBoard := tetris.NewBoardFromBoard(ai.board)
			if err := curBoard.Drop(ai.nextTetromino, curRot, curCol); err != nil {
				// Column is invalid.
				continue
			}

			for nextRot := 0; nextRot < next.RotationsCount(); nextRot++ {
				for nextCol := 0; nextCol < ai.board.Width(); nextCol++ {
					nextBoard := tetris.NewBoardFromBoard(curBoard)
					if err := nextBoard.Drop(next, nextRot, nextCol); err != nil {
						// Column is invalid.
						continue
					}

					eval := evaluate(nextBoard, 0)
					if eval > bestEval {
						bestEval = eval
						bestMoves = []Move{curMove}
					} else if eval == bestEval {
						bestMoves = append(bestMoves, curMove)
					}
				}
			}
		}
	}

	if len(bestMoves) == 0 {
		return fmt.Errorf("AI.DRopSetNext: can not drop tetromino %d", ai.nextTetromino)
	}

	move := bestMoves[rand.Intn(len(bestMoves))]
	if err := ai.board.Drop(ai.nextTetromino, move.rotation, move.column); err != nil {
		return fmt.Errorf("AI.Drop: could not drop: %s", err)
	}

	ai.nextTetromino = next

	return nil
}

func evaluate(board *tetris.Board, depth int) float64 {
	const inf = 1e5
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
			heightsDiff += int(math.Abs(float64(columnHeights[col] - columnHeights[col-1])))
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
