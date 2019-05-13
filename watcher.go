package main

import "github.com/fsnotify/fsnotify"

type watcher struct {
	*fsnotify.Watcher
	dirs []string
}

func newWatcher() (w *watcher, err error) {
	w = new(watcher)
	w.Watcher, err = fsnotify.NewWatcher()
	return
}
