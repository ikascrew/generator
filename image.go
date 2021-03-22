package generator

import (
	"image"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

func scale(name string) (*gocv.Mat, error) {

	img := gocv.IMRead(name, gocv.IMReadColor)
	if img.Empty() {
		return nil, xerrors.Errorf("gocv.IMRead() error: %s", name)
	}
	defer img.Close()

	dst := gocv.NewMatWithSize(480, 720, gocv.MatTypeCV8UC3)

	scale := float64(dst.Cols()) / float64(img.Cols())
	gocv.Resize(img, &dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)

	return &dst, nil
}

func split(dst gocv.Mat, up, down *gocv.Mat, x, y int) error {

	rows := up.Rows()
	y1 := rows - 480 + y
	mat1, err := up.FromPtr(480-y, 720, gocv.MatTypeCV8UC3, y1, 0)
	if err != nil {
		return xerrors.Errorf("up.FromPtr() error: %w", err)
	}
	defer mat1.Close()

	p1 := dst.Region(image.Rect(0, 0, mat1.Cols(), mat1.Rows()))
	defer p1.Close()

	mat2, err := down.FromPtr(y+1, 720, gocv.MatTypeCV8UC3, 0, 0)
	if err != nil {
		return xerrors.Errorf("down.FromPtr() error: %w", err)
	}
	defer mat2.Close()

	y2 := 480 - y - 1
	p2 := dst.Region(image.Rect(0, y2, mat2.Cols(), mat2.Rows()+y2))
	defer p2.Close()

	mat1.CopyTo(&p1)
	mat2.CopyTo(&p2)

	return nil
}
