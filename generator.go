package generator

import (
	"fmt"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

var (
	FPS     float64 = 100.0
	Width   int     = 720
	Height  int     = 480
	Verbose bool    = false
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

	files, err := readDir(dir)
	if err != nil {
		return xerrors.Errorf("readdir() Error: %w", err)
	}

	fm := ""
	digit := 1
	if Verbose {
		base := 1
		for {
			ans := len(files) / base
			if ans == 0 {
				break
			}
			fm = "%0" + fmt.Sprintf("%d", digit) + "d"
			digit++
			base = base * 10
		}
		fmt.Printf("number of files[%d]\n", len(files))
	}
	var now *gocv.Mat
	//var next *gocv.Mat

	for idx, file := range files {

		if Verbose {
			fmt.Printf("files["+fm+"]:%s\n", idx+1, file)
		}
		//if idx == 0 {
		now, err = scale(file, true)
		if err != nil {
			return xerrors.Errorf("scale() Error: %w", err)
		}
		w.Write(now)
		now.Close()
		//} else {

		//now = next
		//}

		//org, err := getOriginal(now)
		//if err != nil {
		//return xerrors.Errorf("getOriginal() error: %w", err)
		//}
		/*
				for row := 0; row < now.Rows()-Height; row++ {

					mat, err := now.FromPtr(Height, Width, gocv.MatTypeCV8UC3, row, 0)
					if err != nil {
						return xerrors.Errorf("Mat FrtomPtr() Error: %w", err)
					}

					p := pasteOriginal(&mat, org)
					w.Write(p)
					p.Close()
				}

			if org != nil {
				org.Close()
			}
				if idx == len(files)-1 {
					next.Close()
					break
				}

				next, err = scale(files[idx+1], true)
				if err != nil {
					return xerrors.Errorf("Scale() Error: %w", err)
				}

				dst := gocv.NewMatWithSize(Height, Width, gocv.MatTypeCV8UC3)

				err = switchPage(w, &dst, now, next)
				if err != nil {
					return xerrors.Errorf("switchPage() error: %w", err)
				}

				dst.Close()
				now.Close()
		*/
	}

	return nil
}

func pasteOriginal(now, org *gocv.Mat) *gocv.Mat {
	if org == nil {
		return now
	}

	e := gocv.NewMatWithSize(Height, Width, gocv.MatTypeCV8UC3)
	alpha := 0.5
	gocv.AddWeighted(*org, float64(alpha), *now, float64(1.0-alpha), 0.0, &e)
	now.Close()
	return &e
}
