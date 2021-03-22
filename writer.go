package generator

import (
	"log"

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
	writer, err := gocv.VideoWriterFile(file, "MJPG", fps, 720, 480, true)
	if err != nil {
		log.Println(err)
	}
	rtn.writer = writer
	return &rtn
}

func (w *FileWriter) Write(mat *gocv.Mat) error {
	w.writer.Write(*mat)
	return nil
}

func (w *FileWriter) Close() error {
	w.writer.Close()
	return nil
}

type DisplayWriter struct {
	Win  *gocv.Window
	Wait int
	Loop bool
}

func NewDisplayWriter(name string, w int) *DisplayWriter {
	var rtn DisplayWriter

	win := gocv.NewWindow("Display:" + name)
	win.ResizeWindow(720, 480)

	rtn.Win = win
	rtn.Wait = w

	return &rtn
}

func (w *DisplayWriter) Write(mat *gocv.Mat) error {
	w.Win.IMShow(*mat)
	w.Win.WaitKey(w.Wait)
	return nil
}

func (w *DisplayWriter) Close() error {
	w.Win.Close()
	return nil
}
