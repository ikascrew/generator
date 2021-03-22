package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/generator"
	"golang.org/x/xerrors"
)

var (
	width   int
	height  int
	fps     float64
	display bool
)

func init() {
	flag.IntVar(&width, "w", 720, "")
	flag.IntVar(&height, "h", 480, "")
	flag.Float64Var(&fps, "fps", 30.0, "")
	flag.BoolVar(&display, "d", false, "display flag()")
}

func main() {

	flag.Parse()
	args := flag.Args()

	err := run(args)
	if err != nil {
		fmt.Printf("Generator Error: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}

func run(args []string) error {

	if len(args) <= 0 {
		return fmt.Errorf("Arguments error:Directory name Required.")
	}

	dir := args[0]

	_, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("Directory does not exits: %w", err)
	}

	generator.FPS = fps
	generator.Width = width
	generator.Height = height

	if display {
		err = generator.Display(dir)
	} else {
		if len(args) >= 2 {
			name := args[1]
			err = generator.Write(dir, name)
		} else {
			err = fmt.Errorf("Arguments error:Output file Required.")
		}
	}

	if err != nil {
		return xerrors.Errorf("run() error: %w", err)
	}

	return nil

}
