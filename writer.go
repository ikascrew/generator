package generator

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"gocv.io/x/gocv"
)

type Writer interface {
	Write(*gocv.Mat) error
	Close() error
}

type FileWriter struct {
	writer *gocv.VideoWriter
}

func NewFileWriter(file string, fps float64) *FileWriter {
	var rtn FileWriter
	writer, err := gocv.VideoWriterFile(file, "mp4v", fps, 720, 480, true)
	if err != nil {
		log.Println(err)
	}
	rtn.writer = writer
	return &rtn
}

func (w *FileWriter) Write(mat *gocv.Mat) error {
	dst := mat.Clone()
	defer dst.Close()
	w.writer.Write(dst)
	return nil
}

func (w *FileWriter) Close() error {
	w.writer.Close()
	return nil
}

type DisplayWriter struct {
	Win    *gocv.Window
	Wait   int
	Before time.Time
	Loop   bool
}

func NewDisplayWriter(name string, w int) *DisplayWriter {
	var rtn DisplayWriter

	win := gocv.NewWindow("Display:" + name)
	win.ResizeWindow(720, 480)
	win.SetWindowProperty(gocv.WindowPropertyAspectRatio, gocv.WindowKeepRatio)

	rtn.Win = win
	rtn.Wait = w

	return &rtn
}

func (w *DisplayWriter) Write(mat *gocv.Mat) error {

	if Verbose {
		now := time.Now()
		sub := now.Sub(w.Before)
		fps := 1000.0 / float64(sub.Milliseconds())

		buf := fmt.Sprintf("FPS:%0.1f", fps)
		gocv.Rectangle(mat, image.Rect(0, 0, 90, 20), color.RGBA{255, 255, 255, 0}, -1)
		pt := image.Pt(5, 20)
		gocv.PutText(mat, buf, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{50, 255, 10, 0}, 2)
		w.Before = now
	}

	w.Win.IMShow(*mat)

	w.Win.WaitKey(w.Wait)
	return nil
}

func (w *DisplayWriter) Close() error {
	w.Win.Close()
	return nil
}
