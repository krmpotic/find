package main

import (
	"log"
	"io/fs"
	"os"
	"strings"
)

var SkipDot = true

// Filter returns
// true, nil -> let through
// false, error -> skip, error can be fs.SkipDir
type Filter func (fs.DirEntry) (bool, error)

func NotHidden(d fs.DirEntry) (bool, error) {
	if d.Name() == "." {
		return true, nil
	}

	hidden := strings.HasPrefix(d.Name(), ".")
	if !hidden {
		return true, nil
	}

	var err error
	if d.IsDir() {
		err = fs.SkipDir
	}
	return false, err
}

func Suffix(ext string) Filter {
	return func (d fs.DirEntry) (bool, error) {
		if strings.HasSuffix(d.Name(), ext) {
			return true, nil
		}
		return false, nil
	}
}

func At(dir string, filters ...Filter) (files []string) {
	fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if SkipDot && d.Name() == "." {
			return nil
		}

		for _, f := range filters {
			if ok, err := f(d); !ok {
				return err
			}
		}

		files = append(files, path)
		return nil
	})

	return files
}

func Or(filters ...Filter) Filter {
	return func (d fs.DirEntry) (bool, error) {
		for _, f := range filters {
			if ok, _ := f(d); ok {
				return true, nil
			}
		}
		return false, nil
	}
}

func main() {
	log.Println(At(".", NotHidden, Or(Suffix(".md"), Suffix(".txt"))))
}
