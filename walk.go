package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (w *watcher) feed() error {
	root, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	err = filepath.Walk(root,
		func(path string, info os.FileInfo, _ error) error {
			if info != nil && info.IsDir() {
				if isInIgnoredMount(path) || isIgnoredDir(path) {
					debugf("walk:  skip %s", path)
					return filepath.SkipDir
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
		return fmt.Errorf("walk error: %w", err)
	}
	if len(w.dirs) == 0 {
		return fmt.Errorf("no dirs to watch")
	}
	return nil
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
