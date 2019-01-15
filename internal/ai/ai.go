// Package ai contains artificial intelligence that plays tetris.
package ai

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ozhi/tetris-ai/internal/tetris"
)

// minUtility and maxUtility define the boundaries of AI's evaluation function.
// The methods utility and evaluate must ony return values in range [minUtility; MaxUtility].
const (
	minUtility      = -float64(1e5)
	maxUtility      = float64(1e5)
	evaluationDepth = 1
)

// AI encapsulates the artificial intelligence logic.
// AI has a reference to a tetris board and the next tetromino that should be dropped.
// By searching the space of potential boards, AI chooses how to rotate and where to drop each tetromino.
// AI uses the minimax algorithm with alpha beta pruning and a utility function.
// The zero value of AI is not usable, method New should be used to create a struct.
type AI struct {
	board    *tetris.Board
	next     tetris.Tetromino
	matrices [][]tetris.TetrominoMatrix
}

// New returns a pointer to a new AI struct.
func New() *AI {
	return &AI{
		board:    tetris.NewBoard(),
		matrices: tetris.TetrominoMatrices(),
	}
}

// Board returns a pointer to the AI's board.
func (ai *AI) Board() *tetris.Board {
	return ai.board
}

// SetNext sets the next tetromino to be dropped by the AI.
// SetNext is usually only called once, before dropping the first tetromino.
// SetNext overwrites if a next tetromino is already set.
// SetNext panics if an empty or invalid tetromino is provided.
func (ai *AI) SetNext(next tetris.Tetromino) {
	if !next.Valid() {
		panic(fmt.Errorf("Ai.SetNext: invalid tetromino %d provided", next))
	}
	ai.next = next
}

// DropSetNext drops the next tetromino and sets the given tetromino as next.
// The given next tetromino is taken into consideration.
// DropSetNext returns error if the tetromino is dropped but that leads to game over.
// DropSetNext panics if the given tetromino is empty or not valid.
// DropSetNext panics if the board is already i game over state. // TODO board should not have such a state.
func (ai *AI) DropSetNext(next tetris.Tetromino) error {
	if !next.Valid() {
		panic(fmt.Errorf("Ai.DropSetNext: invalid tetromino %d provided", next))
	}

	if ai.board.GameOver() {
		panic(fmt.Errorf("AI.DropSetNext: can not drop tetromino %s, game is already over", ai.next))
	}

	type Move struct {
		rotation int
		column   int
	}

	var (
		bestEval  = minUtility - 1
		bestMoves []Move
	)

	for curRot := 0; curRot < ai.next.RotationsCount(); curRot++ {
		curWidth := len(ai.matrices[ai.next][curRot][0])
		for curCol := 0; curCol <= ai.board.Width()-curWidth; curCol++ {
			curMove := Move{
				rotation: curRot,
				column:   curCol,
			}

			curBoard := tetris.NewBoardFromBoard(ai.board)
			if err := curBoard.Drop(ai.next, curRot, curCol); err != nil {
				continue // curBoard's game has just ended.
			}

			for nextRot := 0; nextRot < next.RotationsCount(); nextRot++ {
				nextWidth := len(ai.matrices[next][nextRot][0])
				for nextCol := 0; nextCol <= ai.board.Width()-nextWidth; nextCol++ {
					nextBoard := tetris.NewBoardFromBoard(curBoard)
					if err := nextBoard.Drop(next, nextRot, nextCol); err != nil {
						continue // curBoard's game has just ended.
					}

					eval := ai.evaluate(nextBoard, evaluationDepth, bestEval, maxUtility)
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
		return fmt.Errorf("AI.DropSetNext: can not drop tetromino %s, all moves lead to game over", ai.next)
	}

	move := bestMoves[rand.Intn(len(bestMoves))]
	if err := ai.board.Drop(ai.next, move.rotation, move.column); err != nil {
		return fmt.Errorf("AI.Drop: could not drop: %s", err) // TODO: whats this error?
	}

	ai.next = next

	return nil
}

func (ai *AI) evaluate(board *tetris.Board, depth int, alpha, beta float64) float64 {
	if depth == 0 || board.GameOver() {
		return utility(board)
	}

	minEval := maxUtility + 1
	for _, tetromino := range tetris.Tetrominoes() { // TODO: shuffle tetrominoes?
		maxEval := minUtility - 1
		for rotation := 0; rotation < tetromino.RotationsCount(); rotation++ {
			tetrominoWidth := len(ai.matrices[tetromino][rotation][0])
			for column := 0; column <= board.Width()-tetrominoWidth; column++ {
				newBoard := tetris.NewBoardFromBoard(board)
				if err := newBoard.Drop(tetromino, rotation, column); err != nil {
					// newBoard's game has just ended. Ignore.
				}

				eval := ai.evaluate(newBoard, depth-1, alpha, beta)
				maxEval = math.Max(maxEval, eval)
				alpha = math.Max(alpha, maxEval)
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
