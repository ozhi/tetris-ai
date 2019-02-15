package main

import (
	"flag"
	"fmt"

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
		cli.New().Start()
		return
	}

	err := gui.New().Start()
	if err != nil {
		fmt.Println(err)
	}
}
