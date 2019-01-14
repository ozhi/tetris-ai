package main

import (
	"flag"

	"github.com/ozhi/tetris-ai/internal/cli"
	"github.com/ozhi/tetris-ai/internal/gui"
)

var useCli bool

func init() {
	flag.BoolVar(&useCli, "cli", false, "should the command-line interface be used")
	flag.Parse()
}

func main() {
	if useCli {
		cli.NewGame().Start()
		return
	}

	gui.New().Start()
}
