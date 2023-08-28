package find

import (
	"io/fs"
	"log"
	"os"
)

var SkipDot = true
var IncludeHidden = false

// Filter returns
// true, nil -> let through
// false, error -> skip, error can be fs.SkipDir
type Filter func(fs.DirEntry) (bool, error)

func At(dir string, filters ...Filter) (files []string) {
	if !IncludeHidden {
		filters = append(filters, notHidden)
	}

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
