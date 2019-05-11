package gui

import (
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/ozhi/tetris-ai/internal/tetris"
	"golang.org/x/image/font"
)

type fontOptions struct {
	normal font.Face
	title  font.Face
}

// visualizationOptions are options for the visualization of the tetris game on the screen.
type visualizationOptions struct {
	logoSize int

	cellSize int

	titleBarHeight int
	boardWidth     int
	boardHeight    int
	buttonSize     int

	screenWidth  int
	screenHeight int

	windowTitle string
	scale       float64

	font      *fontOptions
	textColor color.Color

	tetrominoMatrices map[tetris.Tetromino]tetris.TetrominoMatrix
	tetrominoColors   map[tetris.Tetromino]color.Color
	background        color.Color
	boardBackground   color.Color
	borderColor       color.Color
}

// getvisualizationOptions returns the visualizationOptions that the GUI will use.
func getvisualizationOptions() *visualizationOptions {
	cellSize := 40
	titleBarHeight := 3 * cellSize
	boardWidth, boardHeight := 10*cellSize, 20*cellSize
	buttonSize := 5 * cellSize

	return &visualizationOptions{
		logoSize: 8 * cellSize,
		cellSize: cellSize,

		titleBarHeight: titleBarHeight,
		boardWidth:     boardWidth,
		boardHeight:    boardHeight,
		buttonSize:     buttonSize,

		screenWidth:  boardWidth + buttonSize,
		screenHeight: titleBarHeight + boardHeight,

		windowTitle: "Tetris AI",
		scale:       1,

		font:      loadFonts(),
		textColor: color.RGBA{200, 200, 200, 255},

		tetrominoMatrices: loadTetrominoMatrices(),
		tetrominoColors:   loadTetrominoColors(),
		background:        color.RGBA{0, 0, 0, 255},
		boardBackground:   color.RGBA{14, 17, 17, 255},
		borderColor:       color.RGBA{100, 100, 100, 255},
	}
}

// loadTetrominoColors returns a map of the colors of each Tetromino.
func loadTetrominoColors() map[tetris.Tetromino]color.Color {
	return map[tetris.Tetromino]color.Color{
		tetris.TetrominoEmpty: color.RGBA{0, 0, 0, 0},
		tetris.TetrominoI:     color.RGBA{238, 99, 82, 255},
		tetris.TetrominoJ:     color.RGBA{8, 178, 227, 255},
		tetris.TetrominoL:     color.RGBA{49, 136, 139, 255},
		tetris.TetrominoO:     color.RGBA{33, 87, 237, 255},
		tetris.TetrominoS:     color.RGBA{87, 167, 115, 255},
		tetris.TetrominoT:     color.RGBA{76, 101, 99, 255},
		tetris.TetrominoZ:     color.RGBA{128, 35, 142, 255},
	}
}

// loadTetrominoMatrices returns a map of the TetrominoMatrix of each Tetromino.
func loadTetrominoMatrices() map[tetris.Tetromino]tetris.TetrominoMatrix {
	tm := make(map[tetris.Tetromino]tetris.TetrominoMatrix)
	matrices := tetris.TetrominoMatrices()

	for _, tetromino := range tetris.Tetrominoes() {
		// The zeroth rotation of the next tetromino is displayed.
		tm[tetromino] = matrices[tetromino][0]
	}

	return tm
}

// loadFonts loads the fonts to be used throughout the GUI.
func loadFonts() *fontOptions {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(fmt.Errorf("loadFont: %s", err))
	}

	return &fontOptions{
		normal: truetype.NewFace(tt, &truetype.Options{
			Size:    18,
			DPI:     72,
			Hinting: font.HintingFull,
		}),

		title: truetype.NewFace(tt, &truetype.Options{
			Size:    72,
			DPI:     72,
			Hinting: font.HintingFull,
		}),
	}
}
