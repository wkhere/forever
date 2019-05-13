package main

import (
	"os"
	"path/filepath"
)

var ignoredPaths = []string{
	"/dev",
	"/proc",
	"/sys",
}

var ignoredDirs = []string{
	".git",
	"vendor",
	"__pycache__",
	".mypy_cache",
	"ebin",
	"deps",
	"_build",
	"classes",
}

func (w *watcher) feed() {
	root, err := filepath.Abs(".")
	if err != nil {
		fatal("walk: cannot get absolute path:", err)
	}

	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info != nil && info.IsDir() {
				for _, p := range ignoredPaths {
					if dirContains(p, path) {
						debugf("walk:  skip %s", path)
						return filepath.SkipDir
					}
				}
				_, last := filepath.Split(path)
				for _, d := range ignoredDirs {
					if d == last {
						debugf("walk:  skip %s", path)
						return filepath.SkipDir
					}
				}
				err := w.Add(path)
				if err != nil {
					logf("error adding dir %s: %s", path, err)
					return nil
				}
				w.dirs = append(w.dirs, path)
				debugf("walk:  add! %s", path)
			}
			return nil
		})
	if err != nil {
		fatal("walk error:", err)
	}

	if len(w.dirs) == 0 {
		fatal("no dirs to watch")
	}
}

func dirContains(base, path string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	if rel == ".." || (len(rel) >= 3 && rel[:3] == "../") {
		return false
	}
	return true
}
