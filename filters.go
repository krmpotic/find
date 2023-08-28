package find

import (
	"io/fs"
	"strings"
)

func Or(filters ...Filter) Filter {
	return func(d fs.DirEntry) (bool, error) {
		for _, f := range filters {
			if ok, _ := f(d); ok {
				return true, nil
			}
		}
		return false, nil
	}
}

func notHidden(d fs.DirEntry) (bool, error) {
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

func Suffix(ext ...string) Filter {
	return func(d fs.DirEntry) (bool, error) {
		for _, e := range ext {
			if strings.HasSuffix(d.Name(), e) {
				return true, nil
			}
		}
		return false, nil
	}
}

func Dir(d fs.DirEntry) (bool, error) {
	return d.IsDir(), nil
}

func NotDir(d fs.DirEntry) (bool, error) {
	return !d.IsDir(), nil
}
