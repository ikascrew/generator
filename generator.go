package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

var (
	FPS    float64 = 30.0
	Width  int     = 720
	Height int     = 480
)

func Write(dir string, out string) error {

	writer := NewFileWriter(out, FPS)
	defer writer.Close()

	return write(writer, dir)
}

func Display(dir string) error {

	writer := NewDisplayWriter(dir, int(1000.0/FPS))
	defer writer.Close()

	return write(writer, dir)
}

func write(w Writer, dir string) error {

	var err error

	files, err := os.ReadDir(dir)
	if err != nil {
		return xerrors.Errorf("os.Readdir() Error: %w", err)
	}

	var now *gocv.Mat
	var next *gocv.Mat

	for idx, file := range files {

		fmt.Printf("files[%d]:%s\n", idx, file.Name())

		if idx == 0 {
			now, err = scale(filepath.Join(dir, file.Name()))
			if err != nil {
				return xerrors.Errorf("Scale() Error: %w", err)
			}
		} else if idx == len(files)-1 {
			next.Close()
			break
		} else {
			now = next
		}

		for row := 0; row < now.Rows()-480; row++ {
			mat, err := now.FromPtr(480, 720, gocv.MatTypeCV8UC3, row, 0)
			if err != nil {
				return xerrors.Errorf("Mat FrtomPtr() Error: %w", err)
			}
			w.Write(&mat)
			mat.Close()
		}

		next, err = scale(filepath.Join(dir, files[idx+1].Name()))
		if err != nil {
			return xerrors.Errorf("Scale() Error: %w", err)
		}

		dst := gocv.NewMatWithSize(480, 720, gocv.MatTypeCV8UC3)
		for y := 0; y < 480; y++ {
			err := split(dst, now, next, 0, y)
			if err != nil {
				return xerrors.Errorf("error: %w", err)
			}
			w.Write(&dst)
		}
		dst.Close()
		now.Close()

		//fmt.Println(gocv.MatProfile.Count())
		//var b bytes.Buffer
		//gocv.MatProfile.WriteTo(&b, 1)
		//fmt.Print(b.String())
	}

	return nil
}
