// Package gui contains the graphical user interface of Tetris-AI.
package gui

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/ozhi/tetris-ai/internal/ai"
	"github.com/ozhi/tetris-ai/internal/tetris"
	"golang.org/x/image/font"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// GUI is the graphical user interface of Tetris-AI.
// GUI encapsulates an AI that plays tetris and visualization logic.
// The zero value of GUI is not usable, function New should be used to create one.
type GUI struct {
	ai            *ai.AI
	nextTetromino tetris.Tetromino

	automaticMode         bool
	automaticModeTurnedOn chan struct{}

	visualization *visualizationDetails
}

// New creates and initializes a new GUI.
func New() *GUI {
	cellSize := 40
	boardWidth := 10 * cellSize
	boardHeight := 20 * cellSize
	sidebarWidth := 5 * cellSize

	visualization := visualizationDetails{
		cellSize:    cellSize,
		boardWidth:  boardWidth,
		boardHeight: boardHeight,

		sidebarWidth: sidebarWidth,

		windowWidth:  boardWidth + sidebarWidth,
		windowHeight: boardHeight,

		windowTitle: "Tetris AI",
		scale:       1,

		font:      loadFont(),
		textColor: color.RGBA{200, 200, 200, 255},

		tetrominoMatrices: loadTetrominoMatrices(),
		tetrominoColors:   loadTetrominoColors(),
		boardBackground:   color.RGBA{14, 17, 17, 255},
		sidebarBackground: color.RGBA{14, 17, 17, 255},
	}

	ai := ai.New()

	gui := GUI{
		ai:            ai,
		nextTetromino: tetris.RandomTetromino(),

		automaticMode:         false,
		automaticModeTurnedOn: make(chan struct{}),

		visualization: &visualization,
	}

	ai.SetNext(gui.nextTetromino)

	return &gui
}

// Start starts the AI's game and the visualization loop.
func (gui *GUI) Start() error {
	update := func(screen *ebiten.Image) error {
		gui.update()

		if !ebiten.IsDrawingSkipped() {
			gui.draw(screen)
		}

		return nil
	}

	go gui.automaticallyDropTetrominoes()

	err := ebiten.Run(
		update,
		gui.visualization.windowWidth,
		gui.visualization.windowHeight,
		gui.visualization.scale,
		gui.visualization.windowTitle)
	if err != nil {
		return fmt.Errorf("gui.Start: could not start GUI: %s", err)
	}

	return nil
}

// update updates the state of the GUI according to user input.
func (gui *GUI) update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		gui.automaticMode = !gui.automaticMode
		if gui.automaticMode {
			gui.automaticModeTurnedOn <- struct{}{}
		}
	}

	if !gui.automaticMode {
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
			gui.nextTetromino = tetris.RandomTetromino()
			gui.ai.DropSetNext(gui.nextTetromino)
		}
	}
}

// draw draws all images on the screen.
func (gui *GUI) draw(screen *ebiten.Image) {
	_ = screen.DrawImage(gui.boardImage(), &ebiten.DrawImageOptions{})

	nextTetromino := gui.nextTetrominoImage()
	nextTetrominoOptions := ebiten.DrawImageOptions{}
	nextTetrominoOptions.GeoM.Reset()
	nextTetrominoOptions.GeoM.Translate(float64(gui.visualization.boardWidth), 0)
	_ = screen.DrawImage(nextTetromino, &nextTetrominoOptions)
	nextTetromino.Dispose()

	info := gui.infoImage()
	infoOptions := ebiten.DrawImageOptions{}
	infoOptions.GeoM.Reset()
	infoOptions.GeoM.Translate(float64(gui.visualization.boardWidth), float64(gui.visualization.sidebarWidth))
	_ = screen.DrawImage(info, &infoOptions)
	info.Dispose()
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

			op := ebiten.DrawImageOptions{}
			op.GeoM.Reset()
			op.GeoM.Translate(
				float64(col*cellSize),
				float64(row*cellSize),
			)

			_ = image.DrawImage(cell, &op)
		}
	}

	cell.Dispose()
	return image
}

// nextTetrominoImage creates the image of the next tetromino to be dropped.
func (gui *GUI) nextTetrominoImage() *ebiten.Image {
	cellSize := gui.visualization.cellSize

	image, _ := ebiten.NewImage(
		gui.visualization.sidebarWidth,
		gui.visualization.sidebarWidth,
		ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.sidebarBackground)

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)
	_ = cell.Fill(gui.visualization.tetrominoColors[gui.nextTetromino])
	cellOptions := ebiten.DrawImageOptions{}

	tetromino := gui.visualization.tetrominoMatrices[gui.nextTetromino]
	for row := range tetromino {
		for col := range tetromino[row] {
			if tetromino[row][col] {
				cellOptions.GeoM.Reset()
				cellOptions.GeoM.Translate(float64(cellSize), float64(cellSize))
				cellOptions.GeoM.Translate(float64(col*cellSize), float64(row*cellSize))
				_ = image.DrawImage(cell, &cellOptions)
			}
		}
	}

	cell.Dispose()
	return image
}

// infoImage creates the image for the information box.
func (gui *GUI) infoImage() *ebiten.Image {
	image, _ := ebiten.NewImage(
		gui.visualization.sidebarWidth,
		gui.visualization.sidebarWidth,
		ebiten.FilterDefault)
	_ = image.Fill(gui.visualization.sidebarBackground)

	automaticMode := "off"
	if gui.automaticMode {
		automaticMode = "on"
	}

	infos := []string{
		fmt.Sprintf("Tetrominoes: %d", gui.ai.Board().DroppedTetrominoes()),
		fmt.Sprintf("Lines: %d", gui.ai.Board().ClearedLines()),
		"",

		fmt.Sprintf("Automatic mode: %s", automaticMode),
		"(press A)",
		"",

		"Drop next",
		"(press space)",
	}

	for i := range infos {
		text.Draw(image, infos[i], gui.visualization.font, 10, 50+i*20, gui.visualization.textColor)
	}

	return image
}

// automaticallyDropTetrominoes generates random tetrominoes and tells the AI to drop them.
// Only does so in automatic mode is set, otherwise blocks on the gui.automaticModeTurnedOn channel.
func (gui *GUI) automaticallyDropTetrominoes() {
	for {
		// time.Sleep(100 * time.Millisecond)

		if !gui.automaticMode {
			<-gui.automaticModeTurnedOn
		}

		if err := gui.ai.DropSetNext(gui.nextTetromino); err != nil {
			fmt.Printf("AI could not drop tetromino: %s", err)
			break
		}
		gui.nextTetromino = tetris.RandomTetromino()
	}
}

// visualizationDetails are details for visualizing the tetris game on the screen.
type visualizationDetails struct {
	cellSize    int
	boardWidth  int
	boardHeight int

	sidebarWidth int

	windowWidth  int
	windowHeight int

	windowTitle string
	scale       float64

	font      font.Face
	textColor color.Color

	tetrominoMatrices map[tetris.Tetromino]tetris.TetrominoMatrix
	tetrominoColors   map[tetris.Tetromino]color.Color
	boardBackground   color.Color
	sidebarBackground color.Color
}

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

func loadTetrominoMatrices() map[tetris.Tetromino]tetris.TetrominoMatrix {
	tm := make(map[tetris.Tetromino]tetris.TetrominoMatrix)
	matrices := tetris.TetrominoMatrices()

	for _, tetromino := range tetris.Tetrominoes() {
		// The zeroth rotation of the next tetromino is displayed.
		tm[tetromino] = matrices[tetromino][0]
	}

	return tm
}

func loadFont() font.Face {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(fmt.Errorf("loadFont: %s", err))
	}

	return truetype.NewFace(tt, &truetype.Options{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
