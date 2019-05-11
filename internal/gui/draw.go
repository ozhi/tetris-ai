package gui

import (
	"fmt"
	"time"

	"github.com/ozhi/tetris-ai/internal/tetris"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

// draw draws all images on the screen.
func (gui *GUI) draw(screen *ebiten.Image) {
	switch gui.screen {
	case ScreenWelcome:
		draw(screen, gui.welcomeScreen(), 0, 0)

	case ScreenPlay:
		draw(screen, gui.playScreen(), 0, 0)

	default:
		panic(fmt.Errorf("GUI.update: invalid gui screen"))
	}
}

// welcomeScreen returns the image of the welcome screen.
func (gui *GUI) welcomeScreen() *ebiten.Image {
	screenWidth, screenHeight := gui.visualization.screenWidth, gui.visualization.screenHeight

	image, _ := ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.background)

	text.Draw(image,
		"Tetris AI",
		gui.visualization.font.title,
		(screenWidth-350)/2,
		screenHeight/2,
		gui.visualization.textColor)

	draw(image, gui.tetrominoImage(tetris.TetrominoT), (screenWidth-gui.visualization.buttonSize)/2, screenHeight/2)

	return image
}

// playScreen returns the image of the welcome screen.
func (gui *GUI) playScreen() *ebiten.Image {
	image, _ := ebiten.NewImage(
		gui.visualization.screenWidth,
		gui.visualization.screenHeight,
		ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.background)

	draw(image, gui.titleImage(), 0, 0)
	draw(image, gui.boardImage(), 0, gui.visualization.titleBarHeight)
	draw(image, gui.automaticModeButtonImage(), gui.visualization.boardWidth, gui.visualization.titleBarHeight)
	draw(image, gui.nextTetrominoButtonImage(), gui.visualization.boardWidth, gui.visualization.titleBarHeight+gui.visualization.buttonSize)

	strings := []string{
		fmt.Sprintf("Dropped: %d", gui.ai.Board().DroppedTetrominoes()),
		fmt.Sprintf("Cleared lines: %d", gui.ai.Board().ClearedLines()),
		fmt.Sprintf("Time: %s", displayTime(time.Since(gui.gameStart))),
	}
	for i := range strings {

		text.Draw(
			image,
			strings[i],
			gui.visualization.font.normal,
			gui.visualization.boardWidth+10,
			gui.visualization.titleBarHeight+2*gui.visualization.buttonSize+(i+1)*30,
			gui.visualization.textColor)
	}

	return image
}

// tetrominoimage returns the image of thje given tetromino.
func (gui *GUI) tetrominoImage(tetromino tetris.Tetromino) *ebiten.Image {
	cellSize := gui.visualization.cellSize
	buttonSize := gui.visualization.buttonSize

	image, _ := ebiten.NewImage(buttonSize, buttonSize, ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.background)

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)
	_ = cell.Fill(gui.visualization.tetrominoColors[tetromino])

	offset := int(0.5 * float64(cellSize))
	matrix := gui.visualization.tetrominoMatrices[tetromino]
	for row := range matrix {
		for col := range matrix[row] {
			if matrix[row][col] {
				draw(image, cell, col*cellSize+offset, row*cellSize+offset)
			}
		}
	}

	cell.Dispose()

	return image
}

// titleImage creates the image for the information box.
func (gui *GUI) titleImage() *ebiten.Image {
	image, _ := ebiten.NewImage(
		gui.visualization.screenWidth-gui.visualization.titleBarHeight,
		gui.visualization.titleBarHeight,
		ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.background)

	text.Draw(
		image,
		"Tetris AI",
		gui.visualization.font.title,
		(gui.visualization.screenWidth-300)/2,
		gui.visualization.titleBarHeight*2/3,
		gui.visualization.textColor)

	return image
}

// boardImage creates the image of the tetris board.
func (gui *GUI) boardImage() *ebiten.Image {
	cellSize := gui.visualization.cellSize

	image, _ := ebiten.NewImage(10*cellSize, 20*cellSize, ebiten.FilterDefault)
	image.Fill(gui.visualization.boardBackground)

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)

	board := gui.ai.Board()
	for row := 0; row < board.Height(); row++ {
		for col := 0; col < board.Width(); col++ {
			cell.Fill(gui.visualization.tetrominoColors[board.At(row, col)])
			draw(image, cell, col*cellSize, row*cellSize)
		}
	}

	cell.Dispose()
	return image
}

// automaticModeButtonImage creates the image of the "automatic mode" button.
func (gui *GUI) automaticModeButtonImage() *ebiten.Image {
	buttonSize := gui.visualization.buttonSize

	borderedImage, _ := ebiten.NewImage(buttonSize, buttonSize, ebiten.FilterDefault)
	_ = borderedImage.Fill(gui.visualization.borderColor)

	foreground, _ := ebiten.NewImage(buttonSize-6, buttonSize-6, ebiten.FilterDefault)
	_ = foreground.Fill(gui.visualization.background)

	automaticMode := "off"
	if gui.automaticMode {
		automaticMode = "on"
	}
	text.Draw(
		foreground,
		fmt.Sprintf("Automatic mode: %s", automaticMode),
		gui.visualization.font.normal,
		10,
		buttonSize/2+5,
		gui.visualization.textColor)
	draw(borderedImage, foreground, 3, 3)

	return borderedImage
}

// nextTetrominoButtonImage creates the image of the "next tetromino" button.
func (gui *GUI) nextTetrominoButtonImage() *ebiten.Image {
	buttonSize := gui.visualization.buttonSize

	borderedImage, _ := ebiten.NewImage(buttonSize, buttonSize, ebiten.FilterDefault)
	_ = borderedImage.Fill(gui.visualization.borderColor)

	foreground, _ := ebiten.NewImage(buttonSize-6, buttonSize-6, ebiten.FilterDefault)
	_ = foreground.Fill(gui.visualization.background)

	draw(foreground, gui.tetrominoImage(gui.nextTetromino), 0, 0)
	draw(borderedImage, foreground, 3, 3)

	return borderedImage
}

// draw is a helper that draws the given image on top of the target with an offset.
func draw(target, image *ebiten.Image, offsetX, offsetY int) error {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(offsetX), float64(offsetY))
	return target.DrawImage(image, &opts)
}

func displayTime(t time.Duration) string {
	format := func(f float64) string {
		n := int(f) % 60

		if n == 0 {
			return "00"
		}

		if n < 10 {
			return fmt.Sprintf("0%d", int(n))
		}

		return fmt.Sprintf("%d", int(n))
	}

	return fmt.Sprintf("%s:%s", format(t.Minutes()), format(t.Seconds()))
}
