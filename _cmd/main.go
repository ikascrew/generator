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

	var err error
	if display {
		if len(args) >= 1 {
			dir := args[0]
			err = generator.Display(dir)
		} else {
		}
	} else {
		if len(args) >= 2 {
			dir := args[0]
			name := args[1]
			err = generator.Write(dir, name)
		} else {
		}
	}

	if err != nil {
		return xerrors.Errorf("run() error: %w", err)
	}

	return nil

}
