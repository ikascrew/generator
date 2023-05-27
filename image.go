package generator

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
	"golang.org/x/xerrors"
)

func Scale(img image.Image) (image.Image, error) {

	src := img.Bounds()
	odds := float64(Width) / float64(src.Dx())

	dst := image.NewRGBA(image.Rect(0, 0, int(float64(src.Dx())*odds), int(float64(src.Dy())*odds)))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, src, draw.Over, nil)

	return dst, nil
}

//サイズに合わせる
func scale(name string, WH bool) (*gocv.Mat, error) {

	img := gocv.IMRead(name, gocv.IMReadColor)
	if img.Empty() {
		return nil, xerrors.Errorf("gocv.IMRead() error: %s", name)
	}
	defer img.Close()

	dst := gocv.NewMatWithSize(Height, Width, gocv.MatTypeCV8UC3)

	//横幅が足りない時は横をいっぱいにする倍率
	//縦が足りない場合、真ん中にする
	if img.Rows() < Height {

		y := 0
		h := Height - img.Rows()
		y = h / 2

		r := dst.Region(image.Rect(0, y, img.Cols(), img.Rows()+y))
		defer r.Close()

		img.CopyTo(&r)
	} else {
		scale := float64(dst.Cols()) / float64(img.Cols())
		if !WH {
			scale = float64(dst.Rows()) / float64(img.Rows())
		}
		gocv.Resize(img, &dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)
	}

	return &dst, nil
}

const IgnoreHeight = 360

func Whiteouts(img image.Image) ([]image.Image, error) {

	b := img.Bounds()
	var ignores []int

	for y := 0; y < b.Dy(); y++ {
		ignore := true
		for x := 0; x < b.Dx(); x++ {
			c := img.At(x, y)
			if !isWhite(c) {
				ignore = false
				break
			}
		}

		if ignore {
			ignores = append(ignores, y)
		}
	}

	ignores = append(ignores, b.Dy())

	var images []image.Image
	idx := 0
	for _, y := range ignores {

		h := y - idx
		if h <= IgnoreHeight {
			idx = y
			continue
		}

		dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), h))
		ny := 0

		for hy := idx; hy < y; hy++ {
			for x := 0; x < b.Dx(); x++ {
				dst.Set(x, ny, img.At(x, hy))
			}
			ny++
		}
		images = append(images, dst)
		idx = y
	}

	//src の行を取得
	//99% が白だったら行を削除
	//dst := gocv.NewMatWithSize(h, Width, gocv.MatTypeCV8UC3)

	return images, nil
}

const Limit = 58000

func isWhite(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	if r > Limit && g > Limit && b > Limit {
		return true
	}
	return false
}

//上下の画像を位置で変換
func split(dst gocv.Mat, up, down *gocv.Mat, x, y int) error {

	rows := up.Rows()
	y1 := rows - Height + y
	mat1, err := up.FromPtr(Height-y, Width, gocv.MatTypeCV8UC3, y1, 0)
	if err != nil {
		return xerrors.Errorf("up.FromPtr() error: %w", err)
	}
	defer mat1.Close()

	p1 := dst.Region(image.Rect(0, 0, mat1.Cols(), mat1.Rows()))
	defer p1.Close()

	mat2, err := down.FromPtr(y+1, Width, gocv.MatTypeCV8UC3, 0, 0)
	if err != nil {
		return xerrors.Errorf("down.FromPtr() error: %w", err)
	}
	defer mat2.Close()

	y2 := Height - y - 1
	p2 := dst.Region(image.Rect(0, y2, mat2.Cols(), mat2.Rows()+y2))
	defer p2.Close()

	mat1.CopyTo(&p1)
	mat2.CopyTo(&p2)

	return nil
}

func getOriginal(src *gocv.Mat) (*gocv.Mat, error) {

	if src.Rows()/Height < 3 {
		return nil, nil
	}

	scale := float64(Height) / float64(src.Rows())
	dst := gocv.NewMatWithSize(Height, int(float64(Width)*scale), gocv.MatTypeCV8UC3)
	gocv.Resize(*src, &dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)
	pts := [][]image.Point{{
		image.Pt(0, 0),
		image.Pt(Width-1, 0),
		image.Pt(Width-1, Height-1),
		image.Pt(0, Height-1),
	}}
	e := gocv.NewMatWithSize(Height, Width, gocv.MatTypeCV8UC3)
	gocv.FillPoly(&e, pts, color.RGBA{255, 255, 255, 0})

	p := e.Region(image.Rect(0, 0, dst.Cols(), dst.Rows()))

	dst.CopyTo(&p)

	dst.Close()
	p.Close()

	return &e, nil
}
