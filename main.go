package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

func main() {

	nowFile := "test.jpg"

	now, err := Scale(nowFile)
	if err != nil {
		panic(fmt.Sprintf("gocv.IMRead() error[%s] %+v", nowFile, err))
	}

	win := gocv.NewWindow("Test")
	defer win.Close()
	win.ResizeWindow(720, 480)
	//win.SetWindowProperty(gocv.WindowPropertyAspectRatio, gocv.WindowKeepRatio)

	dst := gocv.NewMatWithSize(480, 720, gocv.MatTypeCV8UC3)

	for row := 0; row < now.Rows()-480; row++ {

		mat, err := now.FromPtr(480, 720, gocv.MatTypeCV8UC3, 0, 0)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		win.IMShow(mat)
		win.WaitKey(10)
	}

	win.WaitKey(0)
}

func Scale(name string) (*gocv.Mat, error) {

	img := gocv.IMRead(name, gocv.IMReadColor)
	if img.Empty() {
		panic(fmt.Sprintf("gocv.IMRead() error[%s]", name))
	}
	dst := gocv.NewMatWithSize(480, 720, gocv.MatTypeCV8UC3)

	scale := float64(dst.Cols()) / float64(img.Cols())
	gocv.Resize(img, &dst, image.Point{0, 0}, scale, scale, gocv.InterpolationLinear)

	return &dst, nil
}
