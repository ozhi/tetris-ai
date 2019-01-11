package ai

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ozhi/tetris-ai/internal/tetris"
)

const (
	minUtility = -float64(1e5)
	maxUtility = float64(1e5)
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

	const evaluationDepth = 1
	var (
		bestEval  = minUtility - 1
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

					eval := evaluate(nextBoard, evaluationDepth, bestEval, maxUtility)
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

func evaluate(board *tetris.Board, depth int, alpha, beta float64) float64 {
	if depth < 0 {
		panic("evaluate: depth should be nonnegative")
	}

	if depth == 0 || board.GameOver() {
		return utility(board)
	}

	minEval := maxUtility + 1
	for _, tetromino := range tetris.Tetrominoes() {
		maxEval := minUtility - 1
		for rotation := 0; rotation < tetromino.RotationsCount(); rotation++ {
			for column := 0; column < board.Width(); column++ {
				newBoard := tetris.NewBoardFromBoard(board)
				if err := newBoard.Drop(tetromino, rotation, column); err != nil {
					// Invalid column.
					continue
				}

				eval := evaluate(newBoard, depth-1, alpha, beta)
				maxEval = math.Max(maxEval, eval)
				alpha = math.Min(alpha, maxEval)
				if alpha > beta {
					break
				}
			}
		}

		minEval = math.Min(minEval, maxEval)
		beta = math.Min(beta, minEval)
		if alpha > beta {
			break
		}
	}
	return minEval
}

// bigger is better
func utility(board *tetris.Board) float64 {
	if board.GameOver() {
		return minUtility
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

	if utility < minUtility || utility > maxUtility {
		panic(fmt.Errorf("Invalid utility %f returned", utility))
	}

	return utility
}
