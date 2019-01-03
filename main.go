package main

import (
	"flag"

	"github.com/ozhi/tetris-ai/internal/cli"
	"github.com/ozhi/tetris-ai/internal/gui"
)

var noGUI bool

func init() {
	flag.BoolVar(&noGUI, "nogui", false, "should a graphical user interface be used")
	flag.Parse()
}

func main() {
	if noGUI {
		cli.NewGame().Start()
		return
	}

	gui.NewGame().Start()
}
