package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

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

func feedWatcher(w *fsnotify.Watcher) {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if info != nil && info.IsDir() {
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
				debugf("walk:  add! %s", path)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}
