package gui

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/ozhi/tetris-ai/internal/tetris"
)

// update updates the state of the GUI according to user input.
func (gui *GUI) update() {
	switch gui.screen {
	case ScreenWelcome:
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) || len(inpututil.JustPressedTouchIDs()) > 0 {
			gui.screen = ScreenPlay
			gui.gameStart = time.Now()
		}

	case ScreenPlay:
		if gui.isAutomaticModeJustToggled() {
			gui.automaticMode = !gui.automaticMode
			if gui.automaticMode {
				gui.automaticModeTurnedOn <- struct{}{}
			}
		}

		if !gui.automaticMode {
			if gui.isNextTetrominoJustPressed() {
				gui.nextTetromino = tetris.RandomTetromino()
				gui.ai.DropSetNext(gui.nextTetromino)
			}
		}

	default:
		panic(fmt.Errorf("GUI.update: invalid gui screen"))
	}
}

// isAutomaticModeJustToggled returns true if input for toggling automatic mode
// has been passed at the current frame.
func (gui *GUI) isAutomaticModeJustToggled() bool {
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		return true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > gui.visualization.boardWidth &&
			y > gui.visualization.titleBarHeight &&
			y < gui.visualization.titleBarHeight+gui.visualization.buttonSize {
			return true
		}
	}

	for _, t := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(t)

		if inpututil.TouchPressDuration(t) < 2 &&
			x > gui.visualization.boardWidth &&
			y > gui.visualization.titleBarHeight &&
			y < gui.visualization.titleBarHeight+gui.visualization.buttonSize {
			return true
		}
	}

	return false
}

// isNextTetrominoJustPressed returns true if input for next tetromino
// has been passed on the current frame.
func (gui *GUI) isNextTetrominoJustPressed() bool {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		return true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if x > gui.visualization.boardWidth &&
			y > gui.visualization.titleBarHeight+gui.visualization.buttonSize &&
			y < gui.visualization.titleBarHeight+2*gui.visualization.buttonSize {
			return true
		}
	}

	for _, t := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(t)

		if inpututil.TouchPressDuration(t) < 2 &&
			x > gui.visualization.boardWidth &&
			y > gui.visualization.titleBarHeight+gui.visualization.buttonSize &&
			y < gui.visualization.titleBarHeight+2*gui.visualization.buttonSize {
			return true
		}
	}

	return false
}
