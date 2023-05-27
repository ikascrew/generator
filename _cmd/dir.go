package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikascrew/generator"
	"gocv.io/x/gocv"
	"golang.org/x/term"
	"golang.org/x/xerrors"
)

var W int
var H int
var FPS float64
var LOOP bool
var SUFFIX string

func init() {
	flag.IntVar(&W, "w", 0, "Width")
	flag.IntVar(&H, "h", 0, "Height")
	flag.Float64Var(&FPS, "fps", 24, "FPS")
	flag.BoolVar(&LOOP, "loop", false, "Loop frames")
	flag.StringVar(&SUFFIX, "suffix", "", "Suffix Filename")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "arguments error: required input directory \n")
		os.Exit(1)
	}

	dir := args[0]

	err := run(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run() error:\n%+v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Success")
}

func run(dir string) error {

	fmt.Println(time.Now())
	s, e, imgs, err := loadDir(dir)
	if err != nil {
		return xerrors.Errorf("cut() error: %w", err)
	}
	fmt.Println(time.Now())

	leng := len(imgs)
	fmt.Printf("Frames:[%d]:FPS[%0.3f] = %0.1fs\n", leng, FPS, float64(leng)/FPS)

	ss := getFrameName(s)
	se := getFrameName(e)
	loop := ""
	if LOOP {
		loop = "_loop"
	}

	suffix := ""
	if SUFFIX != "" {
		suffix = "_" + SUFFIX
	}

	name := filepath.Join(dir, fmt.Sprintf("%s-%s%s%s.mp4", ss, se, loop, suffix))
	err = createMP4(name, imgs)
	if err != nil {
		return xerrors.Errorf("cut() error: %w", err)
	}

	fmt.Println(time.Now())

	return nil
}

func loadDir(dir string) (string, string, []image.Image, error) {

	fmt.Printf("Input Directory:[%s]\n", dir)

	if info, err := os.Stat(dir); err != nil {
		return "", "", nil, xerrors.Errorf("os.Stat() error: %w", err)
	} else if !info.IsDir() {
		return "", "", nil, fmt.Errorf("[%s] is not directory", dir)
	}

	files, err := generator.GetFiles(dir)
	if err != nil {
		return "", "", nil, xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	leng := len(files)

	renames := make([]string, leng)
	imgs := make([]image.Image, leng)
	done := make(chan error, 4)

	for idx, p := range files {
		go func(i int, path string) {
			src, err := load(path)
			if err != nil {
				done <- xerrors.Errorf("load() error: %w", err)
				return
			}
			imgs[i] = src
			renames[i] = rename(path)
			done <- nil
		}(idx, p)
	}

	errs := make([]error, leng)
	p := NewProgress(leng)
	p.Prefix = "  Loading ->"
	for idx := range files {
		err := <-done
		errs[idx] = err
		p.Add()
	}
	p.Done()

	err = errors.Join(errs...)
	if err != nil {
		return "", "", nil, xerrors.Errorf("loading error: %w", err)
	}

	s := files[0]
	e := files[len(files)-1]
	if W == 0 || H == 0 {
		decideSize(imgs)
	}

	if SUFFIX != "" {
		for idx := range files {
			fn := files[idx]
			if strings.Index(fn, SUFFIX) != -1 {
				continue
			}
			rename := renames[idx]
			os.Rename(fn, rename)
		}
	}

	return s, e, imgs, nil
}

func decideSize(imgs []image.Image) error {

	maxW := 0
	maxH := 0

	for _, img := range imgs {
		b := img.Bounds()
		if b.Dx() > maxW {
			maxW = b.Dx()
		}
		if b.Dy() > maxH {
			maxH = b.Dy()
		}
	}

	if W == 0 {
		W = maxW
	}
	if H == 0 {
		H = maxH
	}

	return nil
}

func createMP4(out string, imgs []image.Image) error {

	fmt.Printf("Output:[%s] FPS[%0.3f] Loop[%t]\n", out, FPS, LOOP)

	generator.Width = W
	generator.Height = H

	w := generator.NewFileWriter(out, FPS)
	defer w.Close()

	frames := len(imgs)
	if LOOP {
		//first last frames
		frames = frames*2 - 2
	}

	p := NewProgress(frames)
	defer p.Done()

	p.Prefix = "  Writing ->"
	for _, img := range imgs {
		p.Add()
		mat, err := gocv.ImageToMatRGB(img)
		if err != nil {
			return xerrors.Errorf("gocv.ImageToRGBA() error: %w", err)
		}
		if mat.Empty() {
			return xerrors.Errorf("mat is empty")
		}
		w.Write(&mat)
	}

	if LOOP {
		for idx := len(imgs) - 2; idx > 0; idx-- {
			img := imgs[idx]
			p.Add()
			mat, err := gocv.ImageToMatRGB(img)
			if err != nil {
				return xerrors.Errorf("gocv.ImageToRGBA() error: %w", err)
			}
			if mat.Empty() {
				return xerrors.Errorf("mat is empty")
			}
			w.Write(&mat)
		}
	}

	return nil
}

func load(name string) (image.Image, error) {

	fp, err := os.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("os.Open() error: %w", err)
	}
	defer fp.Close()

	img, _, err := image.Decode(fp)
	if err != nil {
		return nil, xerrors.Errorf("image.Decode() error: %w", err)
	}

	return img, nil
}

func getFrameName(n string) string {

	b := filepath.Base(n)
	idx := strings.LastIndex(b, ".")
	if idx == -1 {
		return b
	}

	//最終指定がある場合
	if SUFFIX != "" {
		//そのファイル名になっているかを検索
		sIdx := strings.LastIndex(b, "_"+SUFFIX)
		if sIdx != -1 {
			idx = sIdx
		}
	}

	return b[0:idx]
}

func rename(n string) string {
	if SUFFIX == "" {
		return n
	}

	dir := filepath.Dir(n)
	base := filepath.Base(n)

	newname := base
	idx := strings.LastIndex(base, ".")
	if idx != -1 {
		newname = base[:idx] + "_" + SUFFIX + base[idx:]
	}
	return filepath.Join(dir, newname)
}

type Progress struct {
	Prefix  string
	width   int
	counter int
	max     int
}

func NewProgress(max int) *Progress {
	var p Progress
	p.counter = 0
	p.max = max

	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		p.width = -1
	} else {
		p.width = w
	}
	return &p
}

func (p *Progress) Add() {

	p.counter++
	odds := float64(p.counter) / float64(p.max) * 100

	line := fmt.Sprintf(p.Prefix+"%5.1f%s[%5d/%5d]", odds, "%", p.counter, p.max)
	leng := len(line)
	if p.width != -1 {
		if leng < p.width {
			remain := p.width - leng
			spacer := strings.Repeat(" ", remain)
			line += spacer
		}
	}

	fmt.Print("\r" + line)
}

func (p *Progress) Done() {
	fmt.Println()
}
