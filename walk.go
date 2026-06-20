package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func (w *watcher) addSingle() error {
	path, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s should be a dir", path)
	}
	if isInIgnoredMount(path) || isIgnoredDir(path) {
		return fmt.Errorf("path %s should be skipped", path)
	}

	err = w.Add(path)
	if err != nil {
		if errors.Is(err, syscall.EMFILE) {
			return err // can't progress with too many open files
		}
		return fmt.Errorf("error adding dir %s: %w", path, err)
	}

	w.dirs = append(w.dirs, path)
	return nil
}

func (w *watcher) addRecursive() error {
	root, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	err = filepath.WalkDir(root,
		func(path string, e fs.DirEntry, _ error) error {
			if e.IsDir() {
				if isInIgnoredMount(path) || isIgnoredDir(path) {
					debugf("walk:  skip %s", path)
					return fs.SkipDir
				}
				err := w.Add(path)
				if err != nil {
					if errors.Is(err, syscall.EMFILE) {
						return err // can't progress with too many open files
					}
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
