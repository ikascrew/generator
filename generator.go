package generator

import (
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

func Write(dir string, out string) error {

	writer := NewFileWriter(out, 30.0)
	defer writer.Close()

	return write(writer, dir)
}

func Display(dir string) error {

	//win.SetWindowProperty(gocv.WindowPropertyAspectRatio, gocv.WindowKeepRatio)
	writer := NewDisplayWriter(dir, 10)
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
	dst := gocv.NewMatWithSize(480, 720, gocv.MatTypeCV8UC3)

	for idx, file := range files {

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

		for y := 0; y < 480; y++ {
			err := split(dst, now, next, 0, y)
			if err != nil {
				return xerrors.Errorf("error: %w", err)
			}
			w.Write(&dst)
		}
		now.Close()

		//fmt.Println(gocv.MatProfile.Count())
		//var b bytes.Buffer
		//gocv.MatProfile.WriteTo(&b, 1)
		//fmt.Print(b.String())
	}
	dst.Close()

	return nil
}
