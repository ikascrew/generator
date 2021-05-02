package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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

	if Verbose {
		fmt.Printf("number of files[%d]\n", len(files))
	}
	var now *gocv.Mat
	var next *gocv.Mat

	for idx, file := range files {

		if Verbose {
			fmt.Printf("files[%d]:%s\n", idx+1, file)
		}
		if idx == 0 {
			now, err = scale(file)
			if err != nil {
				return xerrors.Errorf("Scale() Error: %w", err)
			}
		} else {
			now = next
		}

		//TODO どっち方向にスライドするかを受け取る

		for row := 0; row < now.Rows()-Hight; row++ {
			mat, err := now.FromPtr(Hight, Width, gocv.MatTypeCV8UC3, row, 0)
			if err != nil {
				return xerrors.Errorf("Mat FrtomPtr() Error: %w", err)
			}
			w.Write(&mat)
			mat.Close()
		}

		if idx == len(files)-1 {
			next.Close()
			break
		}

		next, err = scale(files[idx+1])
		if err != nil {
			return xerrors.Errorf("Scale() Error: %w", err)
		}

		dst := gocv.NewMatWithSize(Hight, Width, gocv.MatTypeCV8UC3)
		for y := 0; y < Hight; y++ {
			err := split(dst, now, next, 0, y)
			if err != nil {
				return xerrors.Errorf("error: %w", err)
			}
			w.Write(&dst)
		}
		dst.Close()
		now.Close()
	}

	return nil
}

func readDir(root string) ([]string, error) {

	entry, err := os.ReadDir(root)
	if err != nil {
		return nil, xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	sort.Slice(entry, func(i, j int) bool {

		var err1 error
		var num1 int
		name1 := entry[i].Name()
		idx1 := strings.LastIndex(name1, ".")
		if idx1 == -1 {
			err1 = fmt.Errorf("[%s/%s] index error\n", root, name1)
		} else {
			num1, err1 = strconv.Atoi(name1[:idx1])
		}

		var err2 error
		var num2 int
		name2 := entry[j].Name()
		idx2 := strings.LastIndex(name2, ".")
		if idx2 == -1 {
			err2 = fmt.Errorf("[%s/%s] index error\n", root, name2)
		} else {
			num2, err2 = strconv.Atoi(name2[:idx2])
		}

		if err1 != nil && err2 != nil {
			return name1 < name2
		} else {
			return num1 < num2
		}

		err := nil
		if err1 != nil {
			err = err1
		} else {
			err = err2
		}

		log.Printf("filename sort warning:%s %s\n%v\n", name1, name2, err)

		return name1 < name2
	})

	dir := false
	file := false

	files := make([]string, 0)
	for _, elm := range entry {
		path := filepath.Join(root, elm.Name())
		if elm.IsDir() {
			dir = true
			subs, err := readDir(path)
			if err != nil {
				return nil, xerrors.Errorf("readDir() error: %w", err)
			}
			files = append(files, subs...)
			continue
		}
		file = true
		files = append(files, path)
	}

	if file && dir {
		return nil, fmt.Errorf("file dir [%s]", root)
	}
	return files, nil
}
