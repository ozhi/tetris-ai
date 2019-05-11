// Package gui contains the graphical user interface of Tetris-AI.
package gui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/ozhi/tetris-ai/internal/ai"
	"github.com/ozhi/tetris-ai/internal/tetris"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// GUI is the graphical user interface of Tetris-AI.
// GUI encapsulates an AI that plays tetris and visualization logic.
// The zero value of GUI is not usable, function New should be used to create one.
type GUI struct {
	screen        Screen
	visualization *visualizationOptions

	ai            *ai.AI
	nextTetromino tetris.Tetromino

	automaticMode         bool
	automaticModeTurnedOn chan struct{}

	gameStart time.Time
}

// New creates and initializes a new GUI.
func New() *GUI {
	tetromino := tetris.RandomTetromino()

	ai := ai.New()
	ai.SetNext(tetromino)

	return &GUI{
		screen:        ScreenWelcome,
		visualization: getvisualizationOptions(),

		ai:            ai,
		nextTetromino: tetromino,

		automaticMode:         false,
		automaticModeTurnedOn: make(chan struct{}),
	}
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
		gui.visualization.screenWidth,
		gui.visualization.screenHeight,
		gui.visualization.scale,
		gui.visualization.windowTitle)
	if err != nil {
		return fmt.Errorf("gui.Start: could not start GUI: %s", err)
	}

	return nil
}

// automaticallyDropTetrominoes generates random tetrominoes and tells the AI to drop them.
// Only does so if automatic mode is set, otherwise blocks on the gui.automaticModeTurnedOn channel.
func (gui *GUI) automaticallyDropTetrominoes() {
	for {
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
