package generator

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

func GetFiles(in string) ([]string, error) {
	return readDir(in)
}

func readDir(root string) ([]string, error) {

    entry,err := filepath.Glob(filepath.Join(root,"*.png"))
	if err != nil {
		return nil, xerrors.Errorf("filepath.Glob() error: %w", err)
	}
    
/*
	entry, err := os.ReadDir(root)
	if err != nil {
		return nil, xerrors.Errorf("os.ReadDir() error: %w", err)
	}
    */

	sort.Slice(entry, func(i, j int) bool {

		var err1 error
		var num1 int
		name1 := filepath.Base(entry[i])
		idx1 := strings.LastIndex(name1, ".")
		if idx1 == -1 {
			err1 = fmt.Errorf("[%s/%s] index error\n", root, name1)
		} else {
			num1, err1 = strconv.Atoi(name1[:idx1])
		}

		var err2 error
		var num2 int
		name2 := filepath.Base(entry[j])
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

		var err error
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
		files = append(files, elm)
	}

	if file && dir {
		return nil, fmt.Errorf("mixed files and directories[%s]", root)
	}
	return files, nil
}
