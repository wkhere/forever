package main

import (
	"os"
	"path/filepath"
)

func (w *watcher) feed() {
	root, err := filepath.Abs(".")
	if err != nil {
		fatal("walk: cannot get absolute path:", err)
	}

	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
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
		fatal("walk error:", err)
	}

	if len(w.dirs) == 0 {
		fatal("no dirs to watch")
	}
}
