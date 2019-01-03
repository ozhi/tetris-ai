package gui

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/ozhi/tetris-ai/internal/ai"
	"github.com/ozhi/tetris-ai/internal/tetris"
)

const (
	windowWidth  = 200 + 90
	windowHeight = 400
	scale        = 2
	windowTitle  = "Tetris AI"

	cellSize = 20

	sidebarWidth = 100
)

var (
	tetrominoColors   map[tetris.Tetromino]color.Color
	tetrominoMatrices map[tetris.Tetromino][][]bool

	boardBackground         = color.RGBA{14, 17, 17, 255}
	boardBackgroundGameOver = color.RGBA{255, 255, 255, 255}
	sidebarBackground       = color.RGBA{14, 17, 17, 255}
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	tetrominoColors = loadTetrominoColors()
	tetrominoMatrices = loadTetrominoMatrices()
}

// document if zero value of structs is usable
type Game struct {
	ai            *ai.AI
	nextTetromino tetris.Tetromino
}

func NewGame() *Game {
	return &Game{
		ai: ai.New(),
	}
}

func (g *Game) Start() {
	update := func(screen *ebiten.Image) error {
		g.update() // add this is user controls next tetromino

		if !ebiten.IsDrawingSkipped() {
			g.draw(screen)
		}

		return nil
	}

	go g.dropRandomTetrominoes()

	if err := ebiten.Run(update, windowWidth, windowHeight, scale, windowTitle); err != nil {
		panic(fmt.Errorf("Could not start gui game: %s", err))
	}
}

func (g *Game) update() {
	// tetromino, ok := nextTetromino()
	// if !ok {
	// 	// No user input.
	// 	return
	// }

	// if !g.initialized {
	// 	g.ai.SetCurrent(tetromino)
	// 	g.initialized = true
	// 	return
	// }

	// g.ai.DropCurrentSetNext(tetromino)
}

/* func nextTetromino() (tetromino tetris.Tetromino, ok bool) {
	switch {
	case inpututil.IsKeyJustReleased(ebiten.KeyI):
		tetromino = tetris.TetrominoI
	case inpututil.IsKeyJustReleased(ebiten.KeyJ):
		tetromino = tetris.TetrominoJ
	case inpututil.IsKeyJustReleased(ebiten.KeyL):
		tetromino = tetris.TetrominoL
	case inpututil.IsKeyJustReleased(ebiten.KeyO):
		tetromino = tetris.TetrominoO
	case inpututil.IsKeyJustReleased(ebiten.KeyS):
		tetromino = tetris.TetrominoS
	case inpututil.IsKeyJustReleased(ebiten.KeyT):
		tetromino = tetris.TetrominoT
	case inpututil.IsKeyJustReleased(ebiten.KeyZ):
		tetromino = tetris.TetrominoZ
	}

	return tetromino, tetromino != tetris.TetrominoEmpty
} */

func (g *Game) getBoardImage() *ebiten.Image {
	image, _ := ebiten.NewImage(10*cellSize, 20*cellSize, ebiten.FilterDefault)
	if g.ai.Board().GameOver() {
		image.Fill(boardBackgroundGameOver)
	} else {
		image.Fill(boardBackground)
	}

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)

	board := g.ai.Board()
	for row := 0; row < board.Height(); row++ {
		for col := 0; col < board.Width(); col++ {
			cell.Fill(tetrominoColors[board.At(row, col)])

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

func (g *Game) getNextTetrominoImage() *ebiten.Image {
	image, _ := ebiten.NewImage(4*cellSize+10, 4*cellSize+10, ebiten.FilterDefault)
	image.Fill(sidebarBackground)

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)
	cell.Fill(tetrominoColors[g.nextTetromino])
	cellOptions := ebiten.DrawImageOptions{}

	tetromino := tetrominoMatrices[g.nextTetromino]
	for row := range tetromino {
		for col := range tetromino[row] {
			if tetromino[row][col] {
				cellOptions.GeoM.Reset()
				cellOptions.GeoM.Translate(float64(col*cellSize), float64(row*cellSize))
				_ = image.DrawImage(cell, &cellOptions)
			}
		}
	}

	cell.Dispose()
	return image
}

func (g *Game) getInfoImage() *ebiten.Image {
	image, _ := ebiten.NewImage(4*cellSize+10, 4*cellSize+10, ebiten.FilterDefault)
	image.Fill(sidebarBackground)

	cell, _ := ebiten.NewImage(cellSize-1, cellSize-1, ebiten.FilterDefault)
	cell.Fill(tetrominoColors[g.nextTetromino])
	cellOptions := ebiten.DrawImageOptions{}

	tetromino := tetrominoMatrices[g.nextTetromino]
	for row := range tetromino {
		for col := range tetromino[row] {
			if tetromino[row][col] {
				cellOptions.GeoM.Reset()
				cellOptions.GeoM.Translate(float64(col*cellSize), float64(row*cellSize))
				_ = image.DrawImage(cell, &cellOptions)
			}
		}
	}

	cell.Dispose()
	return image
}

func (g *Game) draw(screen *ebiten.Image) {
	_ = screen.DrawImage(g.getBoardImage(), &ebiten.DrawImageOptions{})

	nextTetromino := g.getNextTetrominoImage()
	nextTetrominoOptions := ebiten.DrawImageOptions{}
	nextTetrominoOptions.GeoM.Reset()
	nextTetrominoOptions.GeoM.Translate(10*cellSize, 0)
	_ = screen.DrawImage(nextTetromino, &nextTetrominoOptions)

	nextTetromino.Dispose()
}

func (g *Game) dropRandomTetrominoes() {
	// g.fillBoard()
	g.nextTetromino = tetris.Tetromino(1 + rand.Intn(7))

	for !g.ai.Board().GameOver() {
		time.Sleep(1000 * time.Millisecond)
		g.ai.Drop(g.nextTetromino)
		g.nextTetromino = tetris.Tetromino(1 + rand.Intn(7))
		fmt.Printf("Dropped tetrominoes: %d, Cleared lines: %d\n", g.ai.Board().DroppedTetrominoes(), g.ai.Board().ClearedLines())
	}
}

func (g *Game) fillBoard() {
	for _, offset := range []int{0, 1, 2, 1, 0, 1, 2, 1, 0} {
		for _, col := range []int{0, 2, 4, 6} {
			time.Sleep(50 * time.Millisecond)
			g.ai.Board().Drop(tetris.TetrominoO, 0, offset+col)
		}
	}
	tetrominoes := []tetris.Tetromino{tetris.TetrominoL, tetris.TetrominoO}
	for _, tetromino := range tetrominoes {
		time.Sleep(1000 * time.Millisecond)
		g.ai.Drop(tetromino)
	}
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
	// return map[tetris.Tetromino]color.Color{
	// 	tetris.TetrominoEmpty: color.RGBA{0, 0, 0, 0},
	// 	tetris.TetrominoI:     color.RGBA{0, 255, 255, 255},
	// 	tetris.TetrominoJ:     color.RGBA{0, 0, 255, 255},
	// 	tetris.TetrominoL:     color.RGBA{255, 165, 0, 255},
	// 	tetris.TetrominoO:     color.RGBA{255, 255, 0, 255},
	// 	tetris.TetrominoS:     color.RGBA{0, 255, 0, 255},
	// 	tetris.TetrominoT:     color.RGBA{128, 0, 128, 255},
	// 	tetris.TetrominoZ:     color.RGBA{255, 0, 0, 255},
	// }
}

func loadTetrominoMatrices() map[tetris.Tetromino][][]bool {
	return map[tetris.Tetromino][][]bool{
		tetris.TetrominoEmpty: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, false, false, false},
			[]bool{false, false, false, false},
			[]bool{false, false, false, false},
		},
		tetris.TetrominoI: [][]bool{
			[]bool{false, false, true, false},
			[]bool{false, false, true, false},
			[]bool{false, false, true, false},
			[]bool{false, false, true, false},
		},
		tetris.TetrominoJ: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, false, true, false},
			[]bool{false, false, true, false},
			[]bool{false, true, true, false},
		},
		tetris.TetrominoL: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, true, false, false},
			[]bool{false, true, false, false},
			[]bool{false, true, true, false},
		},
		tetris.TetrominoO: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, true, true, false},
			[]bool{false, true, true, false},
			[]bool{false, false, false, false},
		},
		tetris.TetrominoS: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, false, true, true},
			[]bool{false, true, true, false},
			[]bool{false, false, false, false},
		},
		tetris.TetrominoT: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, true, true, true},
			[]bool{false, false, true, false},
			[]bool{false, false, false, false},
		},
		tetris.TetrominoZ: [][]bool{
			[]bool{false, false, false, false},
			[]bool{false, true, true, false},
			[]bool{false, false, true, true},
			[]bool{false, false, false, false},
		},
	}
}
