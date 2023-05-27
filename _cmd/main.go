package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ikascrew/generator"
	"golang.org/x/xerrors"
)

//ディレクトリを動画にする
var (
	width   int
	height  int
	fps     float64
	display bool
	mode    string
	verbose bool
)

func init() {
	flag.IntVar(&width, "w", 720, "")
	flag.IntVar(&height, "h", 480, "")
	flag.StringVar(&mode, "mode", "slide", "")
	flag.Float64Var(&fps, "fps", 100.0, "")
	flag.BoolVar(&display, "d", false, "display flag()")
	flag.BoolVar(&verbose, "v", false, "print verbose")
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
	generator.Verbose = verbose

	if display {
		err = generator.Display(dir)
	} else {
		name := ""
		if len(args) >= 2 {
			name = args[1]
		} else {
			name = filepath.Clean(dir) + ".mp4"
		}
		err = generator.Write(dir, name)
	}

	if err != nil {
		return xerrors.Errorf("run() error: %w", err)
	}

	return nil

}
