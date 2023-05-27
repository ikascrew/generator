package generator

import (
	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

const Stop = 150

func switchPage(w Writer, dst, now, next *gocv.Mat) error {

	normal := false
	if normal {
		for y := 0; y < Height; y++ {
			err := split(*dst, now, next, 0, y)
			if err != nil {
				return xerrors.Errorf("split() error: %w", err)
			}
			w.Write(dst)
		}
	} else {

		row := 0
		if now.Rows() > Height {
			row = now.Rows() - Height
		}

		//最後のカットを作成
		mat1, err := now.FromPtr(Height, Width, gocv.MatTypeCV8UC3, row, 0)
		if err != nil {
			return xerrors.Errorf("now.FromPtr() error: %w", err)
		}
		org1, err := getOriginal(now)
		if err != nil {
			return xerrors.Errorf("getOriginal() error: %w", err)
		}

		mat2, err := next.FromPtr(Height, Width, gocv.MatTypeCV8UC3, 0, 0)
		if err != nil {
			return xerrors.Errorf("next.FromPtr() error: %w", err)
		}
		org2, err := getOriginal(next)
		if err != nil {
			return xerrors.Errorf("getOriginal() error: %w", err)
		}

		p1 := pasteOriginal(&mat1, org1)
		p2 := pasteOriginal(&mat2, org2)

		defer p1.Close()
		defer p2.Close()

		//nowをしばらく書き込む
		for idx := 0; idx < Stop; idx++ {
			w.Write(p1)
		}

		s := Height - Stop
		for idx := 1; idx < s; idx++ {
			alpha := float64(idx) / float64(s)
			gocv.AddWeighted(*p2, float64(alpha), *p1, float64(1.0-alpha), 0.0, dst)
			w.Write(dst)
		}

		//nextを書き込んで終わり
		w.Write(p2)

	}
	return nil
}
